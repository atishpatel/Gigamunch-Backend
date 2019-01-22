package main

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/atishpatel/Gigamunch-Backend/core/serverhelper"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbadmin"
	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbcommon"

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
	exes, err := serverhelper.PBExecutions(executions)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to PBExecutions")
	}
	resp := &pb.GetExecutionsResp{
		Executions: exes,
		Progress:   getProgress(ctx, executions),
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

	exe, err := serverhelper.PBExecution(execution)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to PBExecution")
	}
	resp := &pb.GetExecutionResp{
		Execution: exe,
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
	exeNew, err := serverhelper.ExecutionFromPb(req.Execution)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to ExecutionFromPb")
	}
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
		exe = getExecutionByMode(ctx, req.Mode, exeOld, exeNew)
	}
	execution, err := exeC.Update(exe)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to update execution")
	}
	exeResp, err := serverhelper.PBExecution(execution)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to PBExecution")
	}
	resp := &pb.UpdateExecutionResp{
		Execution: exeResp,
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

func getProgress(ctx context.Context, exes []*execution.Execution) []*pbcommon.ExecutionProgress {
	exeProgresses := make([]*pbcommon.ExecutionProgress, len(exes))

	for i, exe := range exes {
		exeProgress := &pbcommon.ExecutionProgress{}
		// Head Chef
		hc := progressCounter{}
		cw := progressCounter{}
		cg := progressCounter{}
		if len(exe.Dishes) <= 3 {
			hc.TotalExepectedCount += int8(5-len(exe.Dishes)) * 3
		}

		for _, sticker := range exe.Stickers {
			if sticker.EatingTemperature == "hot" {
				hc.checkEmpty(sticker)
				if sticker.ExtraInstructions == "" {
					hc.TotalExepectedCount--
				}
				if sticker.ReheatInstructions2 == "" {
					hc.TotalExepectedCount--
				}
				if sticker.ReheatOption2 == "" {
					hc.TotalExepectedCount--
				}
				if sticker.ReheatTime2 == "" {
					hc.TotalExepectedCount--
				}
				if sticker.Color == "" {
					hc.TotalExepectedCount--
				}
			}
		}
		if len(exe.Stickers) <= 3 {
			hc.TotalExepectedCount += int8(4-len(exe.Stickers)) * 4
		}

		// Content Writer
		cw.checkEmpty(exe.Culture)
		cw.checkEmpty(exe.CultureCook)

		if len(exe.CultureGuide.InfoBoxes) < 2 {
			cw.TotalExepectedCount += int8(2 - len(exe.CultureGuide.InfoBoxes))
		}

		// Culture Guide
		cg.checkEmpty(exe.Content)
		if exe.Content.DinnerNonVegImageURL == "" {
			cg.TotalExepectedCount--
		}
		if exe.Content.DinnerVegImageURL == "" {
			cg.TotalExepectedCount--
		}
		cg.checkEmpty(exe.Email)
		cg.checkEmpty(exe.Notifications)

		cg.checkEmpty(exe.CultureGuide.MainColor)
		cg.checkEmpty(exe.CultureGuide.FontName)

		// Dishes
		for _, dish := range exe.Dishes {
			// Head Chef
			hc.checkEmpty(dish.Name)
			hc.checkEmpty(dish.Ingredients)
			hc.checkEmpty(dish.ContainerSize)
			hc.addCheck(dish.IsForNonVegetarian || dish.IsForVegetarian)
			// Content Writer
			cw.checkEmpty(dish.Description)
			// cw.checkEmpty(dish.DescriptionPreview)
			// Culture Guide
			cg.checkEmpty(dish.Color)
			// if !dish.IsOnMainPlate {
			// 	cg.checkEmpty(dish.ImageURL)
			// }
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

func getExecutionByMode(ctx context.Context, mode string, exeOld, exeNew *execution.Execution) *execution.Execution {
	exe := exeOld
	switch mode {
	case "captain":
		exe = exeNew
	case "head_chef":
		exe.CultureGuide.DinnerInstructions = exeNew.CultureGuide.DinnerInstructions
		exe.CultureGuide.VegetarianDinnerInstructions = exeNew.CultureGuide.VegetarianDinnerInstructions
		exe.HasChicken = exeNew.HasChicken
		exe.HasBeef = exeNew.HasBeef
		exe.HasPork = exeNew.HasPork
	case "content_writer":
		exe.Culture = exeNew.Culture
		exe.CultureCook = exeNew.CultureCook
	case "culture_guide":
		exe.Content = exeNew.Content
		exe.Email = exeNew.Email
		exe.Notifications = exeNew.Notifications
		exe.CultureGuide.DinnerInstructions = exeNew.CultureGuide.DinnerInstructions
		exe.CultureGuide.VegetarianDinnerInstructions = exeNew.CultureGuide.VegetarianDinnerInstructions
		exe.CultureGuide.MainColor = exeNew.CultureGuide.MainColor
		exe.CultureGuide.FontNamePostScript = exeNew.CultureGuide.FontNamePostScript
		exe.CultureGuide.FontName = exeNew.CultureGuide.FontName
		exe.CultureGuide.FontStyle = exeNew.CultureGuide.FontStyle
		exe.CultureGuide.FontCaps = exeNew.CultureGuide.FontCaps
	}
	// Dishes
	logging.Infof(ctx, "oldExe dishes len(%d) newExe dishes len(%d", len(exe.Dishes), len(exeNew.Dishes))
	if len(exe.Dishes) > len(exeNew.Dishes) {
		exe.Dishes = exe.Dishes[:len(exeNew.Dishes)-1]
	} else if len(exe.Dishes) < len(exeNew.Dishes) {
		d := make([]execution.Dish, len(exeNew.Dishes)-len(exe.Dishes))
		exe.Dishes = append(exe.Dishes, d...)
	}
	logging.Infof(ctx, "oldExe dishes len(%d) newExe dishes len(%d", len(exe.Dishes), len(exeNew.Dishes))
	for i := range exeNew.Dishes {
		switch mode {
		case "head_chef":
			exe.Dishes[i].Number = exeNew.Dishes[i].Number
			exe.Dishes[i].Name = exeNew.Dishes[i].Name
			exe.Dishes[i].Ingredients = exeNew.Dishes[i].Ingredients
			exe.Dishes[i].IsForNonVegetarian = exeNew.Dishes[i].IsForNonVegetarian
			exe.Dishes[i].IsForVegetarian = exeNew.Dishes[i].IsForVegetarian
			exe.Dishes[i].ContainerSize = exeNew.Dishes[i].ContainerSize
		case "content_writer":
			exe.Dishes[i].Description = exeNew.Dishes[i].Description
			exe.Dishes[i].DescriptionPreview = exeNew.Dishes[i].DescriptionPreview
		case "culture_guide":
			exe.Dishes[i].Color = exeNew.Dishes[i].Color
			exe.Dishes[i].IsOnMainPlate = exeNew.Dishes[i].IsOnMainPlate
			exe.Dishes[i].ImageURL = exeNew.Dishes[i].ImageURL
		}
	}
	// Stickers
	if len(exe.Stickers) > len(exeNew.Stickers) {
		exe.Stickers = exe.Stickers[:len(exeNew.Stickers)-1]
	} else if len(exe.Stickers) < len(exeNew.Stickers) {
		d := make([]execution.Sticker, len(exeNew.Stickers)-len(exe.Stickers))
		exe.Stickers = append(exe.Stickers, d...)
	}
	for i := range exeNew.Stickers {
		switch mode {
		case "head_chef":
			exe.Stickers[i].Name = exeNew.Stickers[i].Name
			exe.Stickers[i].Ingredients = exeNew.Stickers[i].Ingredients
			exe.Stickers[i].ExtraInstructions = exeNew.Stickers[i].ExtraInstructions
			exe.Stickers[i].ReheatOption1 = exeNew.Stickers[i].ReheatOption1
			exe.Stickers[i].ReheatOption2 = exeNew.Stickers[i].ReheatOption2
			exe.Stickers[i].ReheatOption1Preferred = exeNew.Stickers[i].ReheatOption1Preferred
			exe.Stickers[i].ReheatTime1 = exeNew.Stickers[i].ReheatTime1
			exe.Stickers[i].ReheatTime2 = exeNew.Stickers[i].ReheatTime2
			exe.Stickers[i].ReheatInstructions1 = exeNew.Stickers[i].ReheatInstructions1
			exe.Stickers[i].ReheatInstructions2 = exeNew.Stickers[i].ReheatInstructions2
			exe.Stickers[i].EatingTemperature = exeNew.Stickers[i].EatingTemperature
			exe.Stickers[i].Number = exeNew.Stickers[i].Number
			exe.Stickers[i].IsForNonVegetarian = exeNew.Stickers[i].IsForNonVegetarian
			exe.Stickers[i].IsForVegetarian = exeNew.Stickers[i].IsForVegetarian
		case "culture_guide":
			exe.Stickers[i].Color = exeNew.Stickers[i].Color
		}
	}

	// CultureGuide.InfoBoxes
	if len(exe.CultureGuide.InfoBoxes) > len(exeNew.CultureGuide.InfoBoxes) {
		exe.CultureGuide.InfoBoxes = exe.CultureGuide.InfoBoxes[:len(exeNew.CultureGuide.InfoBoxes)-1]
	} else if len(exe.CultureGuide.InfoBoxes) < len(exeNew.CultureGuide.InfoBoxes) {
		d := make([]execution.InfoBox, len(exeNew.CultureGuide.InfoBoxes)-len(exe.CultureGuide.InfoBoxes))
		exe.CultureGuide.InfoBoxes = append(exe.CultureGuide.InfoBoxes, d...)
	}
	for i := range exeNew.CultureGuide.InfoBoxes {
		switch mode {
		case "content_writer":
			exe.CultureGuide.InfoBoxes[i].Title = exeNew.CultureGuide.InfoBoxes[i].Title
			exe.CultureGuide.InfoBoxes[i].Text = exeNew.CultureGuide.InfoBoxes[i].Text
			exe.CultureGuide.InfoBoxes[i].Caption = exeNew.CultureGuide.InfoBoxes[i].Caption
		case "culture_guide":
			exe.CultureGuide.InfoBoxes[i].Image = exeNew.CultureGuide.InfoBoxes[i].Image
			exe.CultureGuide.InfoBoxes[i].Caption = exeNew.CultureGuide.InfoBoxes[i].Caption

		}
	}

	return exe
}
