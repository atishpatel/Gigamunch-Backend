package payment

const (
	dateOfBirthFormat = "02-01-2006"
)

// var (
// 	fakeDOB = time.Date(1990, 1, 1, 1, 1, 1, 1, time.UTC)
// )

// // GetSubMerchant returns a SubMerchant
// func (c *Client) GetSubMerchant(subMerchantID string) (*SubMerchantInfo, error) {
// 	if len(subMerchantID) != 32 {
// 		return nil, errInvalidParameter.WithMessage("ID must be length 32.")
// 	}
// 	ma, err := c.bt.MerchantAccount().Find(c.ctx, subMerchantID)
// 	if err != nil {
// 		return nil, errBT.WithError(err).Wrapf("cannot bt.MerchantAccount().Find sub-merchant(%s)", subMerchantID)
// 	}
// 	if ma == nil {
// 		return new(SubMerchantInfo), nil
// 	}
// 	dob, err := time.Parse(dateOfBirthFormat, ma.Individual.DateOfBirth)
// 	if err != nil {
// 		return nil, errInternal.WithError(err).Wrapf("failed to parse time from string(%s)", ma.Individual.DateOfBirth)
// 	}
// 	sm := &SubMerchantInfo{
// 		ID:            ma.Id,
// 		FirstName:     ma.Individual.FirstName,
// 		LastName:      ma.Individual.LastName,
// 		Email:         ma.Individual.Email,
// 		DateOfBirth:   dob,
// 		AccountNumber: ma.FundingOptions.AccountNumber,
// 		RoutingNumber: ma.FundingOptions.RoutingNumber,
// 		Address:       *getAddress(ma.Individual.Address),
// 	}
// 	return sm, nil
// }

// // CreateFakeSubMerchant creates a submerchant account with fake info.
// func (c *Client) CreateFakeSubMerchant(user *types.User, id string) error {
// 	req := &SubMerchantInfo{
// 		ID:            id,
// 		FirstName:     "Fake",
// 		LastName:      "Info",
// 		Email:         user.Email,
// 		DateOfBirth:   fakeDOB,
// 		AccountNumber: "1234567",
// 		RoutingNumber: "064102740",
// 		Address: types.Address{
// 			Street:  "140 W 7th St",
// 			City:    "Cookeville",
// 			State:   "TN",
// 			Zip:     "38501",
// 			Country: "USA",
// 		},
// 	}
// 	nameArray := strings.Split(user.Name, " ")
// 	switch len(nameArray) {
// 	case 3:
// 		if len(nameArray[0]) > 2 {
// 			req.FirstName = nameArray[0]
// 		}
// 		if len(nameArray[2]) > 2 {
// 			req.LastName = nameArray[2]
// 		} else if len(nameArray[1]) > 2 {
// 			req.LastName = nameArray[1]
// 		}
// 	case 2:
// 		if len(nameArray[1]) > 2 {
// 			req.LastName = nameArray[1]
// 		}
// 		fallthrough
// 	case 1:
// 		if len(nameArray[0]) > 2 {
// 			req.FirstName = nameArray[0]
// 		}
// 	}
// 	cookC := cook.New(c.ctx)
// 	err := updateSubMerchant(c.ctx, c.bt, cookC, user, req)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = cookC.UpdateSubMerchantStatus(req.ID, "fake")
// 	if err != nil {
// 		return errors.Wrap("failed to cook.UpdateSubMerchantStatus", err)
// 	}
// 	return nil
// }

// // SubMerchantInfo is the request to UpdateSubMerchant
// // ID must be len 32
// type SubMerchantInfo struct {
// 	ID            string        `json:"-"`
// 	FirstName     string        `json:"first_name"`
// 	LastName      string        `json:"last_name"`
// 	Email         string        `json:"email"`
// 	DateOfBirth   time.Time     `json:"date_of_birth"`
// 	AccountNumber string        `json:"account_number"`
// 	RoutingNumber string        `json:"routing_number"`
// 	Address       types.Address `json:"address"`
// }

// func (req *SubMerchantInfo) valid() error {
// 	if len(req.ID) != 32 {
// 		return errInvalidParameter.WithMessage("ID must be length 32.")
// 	}
// 	return nil
// }

// // UpdateSubMerchant creates or updates sub-merchant info
// func (c *Client) UpdateSubMerchant(user *types.User, req *SubMerchantInfo) error {
// 	err := req.valid()
// 	if err != nil {
// 		return err
// 	}
// 	cookC := cook.New(c.ctx)
// 	return updateSubMerchant(c.ctx, c.bt, cookC, user, req)
// }

// func updateSubMerchant(ctx context.Context, bt *braintree.Braintree, cookC cookInterface, user *types.User, req *SubMerchantInfo) error {
// 	ma := &braintree.MerchantAccount{
// 		MasterMerchantAccountId: "Gigamunch_marketplace",
// 		Id:          req.ID,
// 		TOSAccepted: true,
// 		Individual: &braintree.MerchantAccountPerson{
// 			FirstName:   req.FirstName,
// 			LastName:    req.LastName,
// 			Email:       req.Email,
// 			DateOfBirth: req.DateOfBirth.Format(dateOfBirthFormat),
// 			Address:     getBTAddress(req.Address),
// 		},
// 		FundingOptions: &braintree.MerchantAccountFundingOptions{
// 			Destination:   braintree.FUNDING_DEST_BANK,
// 			AccountNumber: req.AccountNumber,
// 			RoutingNumber: req.RoutingNumber,
// 		},
// 	}
// 	_, err := bt.MerchantAccount().Find(ctx, req.ID)
// 	if err == nil {
// 		ma, err = bt.MerchantAccount().Update(ctx, ma)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "number is invalid") {
// 				return errInvalidParameter.WithMessage(err.Error())
// 			}
// 			return errBT.WithError(err).Wrapf("cannot bt.MerchantAccount().Update sub-merchant(%s)", req.ID)
// 		}
// 	} else {
// 		ma, err = bt.MerchantAccount().Create(ctx, ma)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "number is invalid") {
// 				return errInvalidParameter.WithMessage(err.Error())
// 			}
// 			return errBT.WithError(err).Wrapf("cannot create sub-merchant(%s)", req.ID)
// 		}
// 	}
// 	if ma.Status != "pending" && ma.Status != "active" {
// 		return errBT.WithMessage("Error creating sub-merchant account with status " + ma.Status)
// 	}
// 	if !user.HasSubMerchantID() {
// 		user.SetSubMerchantID(true)
// 		err = auth.SaveUser(ctx, user)
// 		if err != nil {
// 			return errors.Wrap("failed update user to has sub-merchant account", err)
// 		}
// 	}
// 	_, err = cookC.UpdateSubMerchantStatus(ma.Id, ma.Status)
// 	if err != nil {
// 		return errors.Wrap("failed to cookC.UpdateSubMerchantStatus", err)
// 	}
// 	return nil
// }

// type cookInterface interface {
// 	FindBySubMerchantID(string) (*cook.Cook, error)
// 	UpdateSubMerchantStatus(string, string) (*cook.Cook, error)
// 	Notify(string, string, string) error
// }

// // DisbursementException handles a disbursement exception notification
// func (c Client) DisbursementException(signature, payload string) ([]string, error) {
// 	cookC := cook.New(c.ctx)
// 	return disbursementException(c.ctx, signature, payload, c.bt, cookC)
// }

// func disbursementException(ctx context.Context, signature, payload string, bt *braintree.Braintree, cookC cookInterface) ([]string, error) {
// 	notification, err := parseNotification(ctx, signature, payload, braintree.DisbursementExceptionWebhook, bt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if notification == nil || notification.Disbursement().MerchantAccount == nil {
// 		return nil, errInternal.Wrapf("there was an error with notification (%+v) subject(%+v) MerchantAccount(%+v) Disbursement(%+v)", notification, notification.Subject, notification.MerchantAccount(), notification.Disbursement())
// 	}
// 	cook, err := cookC.UpdateSubMerchantStatus(notification.Disbursement().MerchantAccount.Id, "disbursement")
// 	if err != nil {
// 		return nil, errors.Wrap("failed to update cook by submerchantID", err)
// 	}
// 	disbursement := notification.Disbursement()
// 	message := fmt.Sprintf("A transaction to your account failed because '%s' please take the following action: '%s'", disbursement.ExceptionMessage, disbursement.FollowUpAction)
// 	err = cookC.Notify(cook.ID, "There was a problem sending money to you! Please update your banking info so you can get paid. - Gigamunch", message) // TODO add update your info button
// 	if err != nil {
// 		return nil, errors.Wrap("failed to notify cook", err)
// 	}
// 	return disbursement.TransactionIds, nil
// }

// // SubMerchantApproved handles a submerchant approved notification
// func (c *Client) SubMerchantApproved(signature, payload string) error {
// 	cookC := cook.New(c.ctx)
// 	return subMerchantApproved(c.ctx, signature, payload, c.bt, cookC)
// }

// func subMerchantApproved(ctx context.Context, signature, payload string, bt *braintree.Braintree, cookC cookInterface) error {
// 	notification, err := parseNotification(ctx, signature, payload, braintree.SubMerchantAccountApprovedWebhook, bt)
// 	if err != nil {
// 		return err
// 	}
// 	merch := notification.MerchantAccount()
// 	_, err = cookC.UpdateSubMerchantStatus(merch.Id, merch.Status)
// 	if err != nil {
// 		return errors.Wrap(fmt.Sprintf("failed to update submerchant(%s) status(%s)", merch.Id, merch.Status), err)
// 	}
// 	return nil
// }

// // SubMerchantDeclined handles a submerchant declined notification
// func (c *Client) SubMerchantDeclined(signature, payload string) error {
// 	cookC := cook.New(c.ctx)
// 	return subMerchantDeclined(c.ctx, signature, payload, c.bt, cookC)
// }

// func subMerchantDeclined(ctx context.Context, signature, payload string, bt *braintree.Braintree, cookC cookInterface) error {
// 	notification, err := parseNotification(ctx, signature, payload, braintree.SubMerchantAccountDeclinedWebhook, bt)
// 	if err != nil {
// 		return err
// 	}
// 	merch := notification.MerchantAccount()
// 	cook, err := cookC.UpdateSubMerchantStatus(merch.Id, merch.Status)
// 	if err != nil {
// 		return errors.Wrap(fmt.Sprintf("failed to update submerchant(%s) status(%s)", merch.Id, merch.Status), err)
// 	}
// 	err = cookC.Notify(cook.ID, "There was a problem with the approving your bank info - Gigamunch", notification.Subject.APIErrorResponse.ErrorMessage)
// 	if err != nil {
// 		return errors.Wrap("failed to notify cook", err)
// 	}
// 	return nil
// }

// func parseNotification(ctx context.Context, signature, payload, wantKind string, bt *braintree.Braintree) (*braintree.WebhookNotification, error) {
// 	notification, err := bt.WebhookNotification().Parse(signature, payload)
// 	if err != nil {
// 		return nil, errBT.WithError(err).Wrap("cannot parse payment notification")
// 	}

// 	if notification.Kind != wantKind {
// 		return nil, errInvalidParameter.Wrapf("%#v notification is not a %s notification", notification, wantKind)
// 	}
// 	return notification, nil
// }
