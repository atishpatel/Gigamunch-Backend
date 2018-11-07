package admin

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/atishpatel/Gigamunch-Backend/core/serverhelper"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/common"

	"github.com/atishpatel/Gigamunch-Backend/core/execution"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// GetExecutions gets list of executions.
func (s *server) GetExecutions(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetExecutionsReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	exeC, err := execution.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution client")
	}
	executions, err := exeC.GetAll(int(req.Start), int(req.Limit))
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get all executions")
	}
	log.Infof(ctx, "return %d executions", len(executions))

	resp := &pb.GetExecutionsResp{
		Executions: serverhelper.PBExecutions(executions),
		Progress:   getProgress(executions),
	}
	return resp
}

// GetExecution gets an execution.
func (s *server) GetExecution(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetExecutionReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	exeC, err := execution.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution client")
	}
	execution, err := exeC.Get(req.IDOrDate)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution")
	}

	resp := &pb.GetExecutionResp{
		Execution: serverhelper.PBExecution(execution),
	}
	return resp
}

// UpdateExecution updates or creates an execution.
func (s *server) UpdateExecution(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.UpdateExecutionReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	if req.Mode == "" {
		return errBadRequest.WithMessage("Mode must be selected")
	}

	exeC, err := execution.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution client")
	}
	exeNew := serverhelper.ExecutionFromPb(req.Execution)
	var exe *execution.Execution
	if exeNew.ID == 0 {
		// Create
		exe = exeNew
	} else {
		// Update
		exeOld, err := exeC.Get(strconv.FormatInt(exeNew.ID, 10))
		if err != nil {
			return errors.Annotate(err, "failed to exection.Get")
		}
		exe = getExecutionByMode(req.Mode, exeOld, exeNew)
	}
	execution, err := exeC.Update(exe)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to update execution")
	}

	resp := &pb.UpdateExecutionResp{
		Execution: serverhelper.PBExecution(execution),
	}
	return resp
}

type progressCounter struct {
	ValidCount          int8
	TotalExepectedCount int8
}

func (c *progressCounter) addCheck(success bool) {
	if success {
		c.ValidCount++
	}
	c.TotalExepectedCount++
}

func (c *progressCounter) checkEmpty(object interface{}) {
	if object == nil {
		c.addCheck(false)
		return
	} else if object == "" {
		c.addCheck(false)
		return
	}

	// for arrays
	v := reflect.ValueOf(object)
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			c.checkEmpty(v.Index(i).Interface())
		}
		return
	}

	// for structs
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			c.checkEmpty(v.Field(i).Interface())
		}
		return
	}
	// ignore booleans, ints, floats
	c.addCheck(true)
}

func (c *progressCounter) getPercent() int32 {
	return int32((float32(c.ValidCount) / float32(c.TotalExepectedCount)) * 100)
}

func getProgress(exes []*execution.Execution) []*pbcommon.ExecutionProgress {
	exeProgresses := make([]*pbcommon.ExecutionProgress, len(exes))

	for i, exe := range exes {
		exeProgress := &pbcommon.ExecutionProgress{}
		// Head Chef
		hc := progressCounter{}
		if len(exe.Dishes) < 5 {
			hc.TotalExepectedCount += int8(5-len(exe.Stickers)) * 3
		}

		hc.checkEmpty(exe.CultureGuide.DinnerInstructions)
		hc.checkEmpty(exe.CultureGuide.VegetarianDinnerInstructions)

		// Content Writer
		cw := progressCounter{}
		cw.checkEmpty(exe.Culture)
		cw.checkEmpty(exe.CultureCook)

		if len(exe.CultureGuide.InfoBoxes) < 2 {
			cw.TotalExepectedCount += int8(2 - len(exe.CultureGuide.InfoBoxes))
		}

		// Culture Guide
		cg := progressCounter{}
		cg.checkEmpty(exe.Content)
		cg.checkEmpty(exe.Email)
		cg.checkEmpty(exe.Notifications)
		for _, sticker := range exe.Stickers {
			if sticker.EatingTemperature == "" || sticker.EatingTemperature == "hot" {
				cg.checkEmpty(sticker)
			}
		}
		if len(exe.Stickers) < 4 {
			cg.TotalExepectedCount += int8(4 - len(exe.Stickers))
		}
		hc.checkEmpty(exe.CultureGuide.MainColor)
		hc.checkEmpty(exe.CultureGuide.FontName)

		// Dishes
		for _, dish := range exe.Dishes {
			// Head Chef
			hc.checkEmpty(dish.Name)
			hc.checkEmpty(dish.Ingredients)
			hc.addCheck(!dish.IsForNonVegetarian && !dish.IsForVegetarian)
			// Content Writer
			cw.checkEmpty(dish.Description)
			cw.checkEmpty(dish.DescriptionPreview)
			// Culture Guide
			cg.checkEmpty(dish.Color)
			if !dish.IsOnMainPlate {
				cg.checkEmpty(dish.ImageURL)
			}
		}

		// set progress
		exeProgress.HeadChef = hc.getPercent()
		exeProgress.ContentWriter = cw.getPercent()
		exeProgress.CultureGuide = cg.getPercent()

		// add summary
		dishCountNonVeg := 0
		dishCountVeg := 0
		for _, dish := range exe.Dishes {
			if dish.IsForVegetarian {
				dishCountVeg++
			}
			if dish.IsForNonVegetarian {
				dishCountNonVeg++
			}
		}

		exeProgress.Summary = append(exeProgress.Summary, &pbcommon.ExecutionProgressSummary{
			Message: fmt.Sprintf("Non-veg Dishes: %d", dishCountNonVeg),
			IsError: (dishCountNonVeg < 4),
		})
		exeProgress.Summary = append(exeProgress.Summary, &pbcommon.ExecutionProgressSummary{
			Message: fmt.Sprintf("Veg Dishes: %d", dishCountVeg),
			IsError: (dishCountVeg < 4),
		})
		exeProgress.Summary = append(exeProgress.Summary, &pbcommon.ExecutionProgressSummary{
			Message: fmt.Sprintf("CG Info Boxes: %d", len(exe.CultureGuide.InfoBoxes)),
			IsError: (len(exe.CultureGuide.InfoBoxes) < 2),
		})

		// add progress
		exeProgresses[i] = exeProgress
	}
	return exeProgresses
}

func getExecutionByMode(mode string, exeOld, exeNew *execution.Execution) *execution.Execution {
	exe := exeOld
	switch mode {
	case "captain":
		exe = exeNew
	case "head_chef":
		exe.CultureGuide.DinnerInstructions = exeNew.CultureGuide.DinnerInstructions
		exe.CultureGuide.VegetarianDinnerInstructions = exeNew.CultureGuide.VegetarianDinnerInstructions
	case "content_writer":
		exe.Culture = exeNew.Culture
		exe.CultureCook = exeNew.CultureCook
		exe.CultureGuide.InfoBoxes = exeNew.CultureGuide.InfoBoxes
	case "culture_guide":
		exe.Content = exeNew.Content
		exe.Email = exeNew.Email
		exe.Notifications = exeNew.Notifications
		exe.Stickers = exeNew.Stickers
		exe.CultureGuide.DinnerInstructions = exeNew.CultureGuide.DinnerInstructions
		exe.CultureGuide.VegetarianDinnerInstructions = exeNew.CultureGuide.VegetarianDinnerInstructions
		exe.CultureGuide.MainColor = exeNew.CultureGuide.MainColor
		exe.CultureGuide.FontNamePostScript = exeNew.CultureGuide.FontNamePostScript
		exe.CultureGuide.FontName = exeNew.CultureGuide.FontName
		exe.CultureGuide.FontStyle = exeNew.CultureGuide.FontStyle
		exe.CultureGuide.FontCaps = exeNew.CultureGuide.FontCaps
	}

	if len(exe.Dishes) < len(exeNew.Dishes) {
		exe.Dishes = exe.Dishes[:len(exeNew.Dishes)-1]
	} else if len(exe.Dishes) > len(exeNew.Dishes) {
		d := make([]execution.Dish, len(exe.Dishes)-len(exeNew.Dishes))
		exe.Dishes = append(exe.Dishes, d...)
	}

	for i := range exeNew.Dishes {
		switch mode {
		case "head_chef":
			exe.Dishes[i].Number = exeNew.Dishes[i].Number
			exe.Dishes[i].Name = exeNew.Dishes[i].Name
			exe.Dishes[i].Ingredients = exeNew.Dishes[i].Ingredients
			exe.Dishes[i].IsForNonVegetarian = exeNew.Dishes[i].IsForNonVegetarian
			exe.Dishes[i].IsForVegetarian = exeNew.Dishes[i].IsForVegetarian
		case "content_writer":
			exe.Dishes[i].Description = exeNew.Dishes[i].Description
			exe.Dishes[i].DescriptionPreview = exeNew.Dishes[i].DescriptionPreview
		case "culture_guide":
			exe.Dishes[i].Color = exeNew.Dishes[i].Color
			exe.Dishes[i].IsOnMainPlate = exeNew.Dishes[i].IsOnMainPlate
			exe.Dishes[i].ImageURL = exeNew.Dishes[i].ImageURL
		}
	}
	return exe
}
