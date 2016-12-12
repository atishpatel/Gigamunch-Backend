package promo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	// driver for mysql
	mysql "github.com/go-sql-driver/mysql"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/types"
	"google.golang.org/appengine"
)

const (
	datetimeFormat               = "2006-01-02 15:04:05" // "Jan 2, 2006 at 3:04pm (MST)"
	insertUsedPromoCodeStatement = "INSERT INTO `used_promo_code` (code,eater_id,inquiry_id,state) VALUES ('%s','%s',%d,%d)"
	selectUsedPromoCodeStatement = "SELECT inquiry_id,state FROM `used_promo_code` WHERE eater_id='%s' AND code='%s'"
	updateUsedPromoCodeStatement = "UPDATE `used_promo_code` SET state=%d WHERE eater_id='%s' AND inquiry_id=%d"
	insertPromoCodeStatement     = "INSERT INTO `promo_code` (code,free_delivery,percent_off,amount_off,discount_cap,free_dish,buy_one_get_one_free,start_datetime,end_datetime,num_uses) VALUES ('%s',%t,%d,%f,%f,%t,%t,'%s','%s',%d)"
	selectPromoCodesStatement    = "SELECT created_datetime,free_delivery,percent_off,amount_off,discount_cap,free_dish,buy_one_get_one_free,start_datetime,end_datetime,num_uses FROM `promo_code` WHERE code='%s'"
)

var (
	connectOnce = sync.Once{}
	mysqlDB     *sql.DB
	errSQLDB    = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with cloud sql database."}
	// errBuffer           = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "An unknown error occured."}
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errInvalidPromoCode = errors.ErrorWithCode{Code: errors.CodeInvalidPromoCode, Message: "Invalid promo code."}
)

// Client is the client for likes
type Client struct {
	ctx context.Context
}

// New returns a new Client for user likes
func New(ctx context.Context) *Client {
	connectOnce.Do(func() {
		connectSQL(ctx)
	})
	return &Client{ctx: ctx}
}

// InsertCode inserts a promo code.
func (c *Client) InsertCode(user *types.User, code *Code) error {
	if !user.IsAdmin() {
		return errInvalidParameter.WithMessage("User is not an admin.")
	}
	if code == nil {
		return errInvalidParameter.WithMessage("Code cannot be nil.")
	}
	if strings.Contains(code.Code, "'") {
		return errInvalidParameter.WithMessage("Code cannot contain '.")
	}
	code.Code = strings.ToUpper(code.Code)
	if code.PercentOff > 0 || code.AmountOff > 0 || code.BuyOneGetOneFree || code.FreeDish && code.DiscountCap < .001 {
		return errInvalidParameter.WithMessage("Discount Cap cannot be 0.")
	}
	st := fmt.Sprintf(insertPromoCodeStatement, code.Code, code.FreeDelivery, code.PercentOff, code.AmountOff, code.DiscountCap, code.FreeDish, code.BuyOneGetOneFree, code.StartDatetime.Format(datetimeFormat), code.EndDatetime.Format(datetimeFormat), code.NumUses)
	_, err := mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute statement: " + st)
	}
	return nil
}

// GetCodeInfo checks if a promo code is valid.
func (c *Client) GetCodeInfo(code string, eaterPoint, cookPoint types.GeoPoint) (*Code, error) {
	if code == "" || strings.Contains(code, "'") {
		return nil, errInvalidParameter.WithMessage("Promo code is empty.")
	}

	code = strings.ToUpper(code)
	st := fmt.Sprintf(selectPromoCodesStatement, code)
	codeInfo := &Code{Code: code}
	rows, err := mysqlDB.Query(st)
	if err != nil {
		return nil, errSQLDB.WithError(err).Wrap("failed to query statement:" + st)
	}
	defer handleCloser(c.ctx, rows)
	found := false
	for rows.Next() {
		var createdNulltime mysql.NullTime
		var startNulltime mysql.NullTime
		var endNulltime mysql.NullTime
		err = rows.Scan(&createdNulltime, &codeInfo.FreeDelivery, &codeInfo.PercentOff, &codeInfo.AmountOff, &codeInfo.DiscountCap, &codeInfo.FreeDish, &codeInfo.BuyOneGetOneFree, &startNulltime, &endNulltime, &codeInfo.NumUses)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("failed to rows.Scan")
		}
		if createdNulltime.Valid {
			codeInfo.CreatedDatetime = createdNulltime.Time
		}
		if startNulltime.Valid {
			codeInfo.StartDatetime = startNulltime.Time
		}
		if endNulltime.Valid {
			codeInfo.EndDatetime = endNulltime.Time
		}
		found = true
	}
	if !found {
		return nil, errInvalidPromoCode.WithMessage("Invalid Promo Code.")
	}
	now := time.Now()
	if now.Before(codeInfo.StartDatetime) || now.After(codeInfo.EndDatetime) {
		return nil, errInvalidPromoCode.WithMessage("The promo code has expired.")
	}
	if codeInfo.FreeDelivery {
		inRange := types.InGigadeliveryRange(cookPoint, eaterPoint)
		if !inRange {
			return nil, errInvalidPromoCode.WithMessage("Free delivery is only available in Nashville.")
		}
	}
	// TODO add num uses code
	return codeInfo, nil
}

// GetUsableCodeForUser checks if a code is valid and if a user has already used it.
func (c *Client) GetUsableCodeForUser(code, eaterID string, eaterPoint, cookPoint types.GeoPoint) (*Code, error) {
	code = strings.ToUpper(code)
	promoCode, err := c.GetCodeInfo(code, eaterPoint, cookPoint)
	if err != nil {
		return nil, err
	}
	if eaterID != "" {
		st := fmt.Sprintf(selectUsedPromoCodeStatement, eaterID, code)
		rows, err := mysqlDB.Query(st)
		if err != nil {
			return nil, errSQLDB.WithError(err).Wrap("failed to query statement:" + st)
		}
		defer handleCloser(c.ctx, rows)
		var inquiryID int64
		var state State
		for rows.Next() {
			err = rows.Scan(&inquiryID, &state)
			if err != nil {
				return nil, errSQLDB.WithError(err).Wrap("failed to rows.Scan")
			}
			if state == Used {
				return nil, errInvalidPromoCode.WithMessage("Promo code has already been used.").Wrapf("promo code used by inquiry(%d)", inquiryID)
			}
			if state == Pending {
				return nil, errInvalidPromoCode.WithMessage("There is an open request with the promo code.")
			}
		}
	}
	return promoCode, nil
}

// InsertUsedCode inserts the promo code info.
func (c *Client) InsertUsedCode(code, eaterID string, inquiryID int64, state State) error {
	if code == "" || strings.Contains(code, "'") {
		return errInvalidParameter.Wrap("Code cannot be empty")
	}
	if eaterID == "" {
		return errInvalidParameter.Wrap("EaterID cannot be empty")
	}
	if inquiryID == 0 {
		return errInvalidParameter.Wrap("InquiryID cannot be 0")
	}
	code = strings.ToUpper(code)
	st := fmt.Sprintf(insertUsedPromoCodeStatement, code, eaterID, inquiryID, state)
	_, err := mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute statement: " + st)
	}
	return nil
}

// UpdateUsedCodeState updates the state of promo code.
func (c *Client) UpdateUsedCodeState(eaterID string, inquiryID int64, state State) error {
	if inquiryID == 0 {
		return errInvalidParameter.Wrap("InquiryID cannot be 0")
	}
	if eaterID == "" {
		return errInvalidParameter.Wrap("EaterID cannot be empty")
	}
	st := fmt.Sprintf(updateUsedPromoCodeStatement, state, eaterID, inquiryID)
	_, err := mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrap("failed to execute statement: " + st)
	}
	return nil
}

func connectSQL(ctx context.Context) {
	var err error
	var connectionString string
	if appengine.IsDevAppServer() {
		// "user:password@/dbname"
		connectionString = "root@/gigamunch"
	} else {
		connectionString = os.Getenv("MYSQL_CONNECTION")
	}
	mysqlDB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Couldn't connect to mysql database")
	}
}

type closer interface {
	Close() error
}

func handleCloser(ctx context.Context, c closer) {
	err := c.Close()
	if err != nil {
		utils.Errorf(ctx, "Error closing rows: %v", err)
	}
}
