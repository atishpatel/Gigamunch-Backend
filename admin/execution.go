package admin

import (
	"context"
	"net/http"

	"github.com/atishpatel/Gigamunch-Backend/subserver"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/common"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
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
		Executions: subserver.PBExecutions(executions),
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
	execution, err := exeC.Get(req.IdOrDate)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution")
	}

	resp := &pb.GetExecutionResp{
		Execution: subserver.PBExecution(execution),
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

	exeC, err := execution.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution client")
	}
	execution, err := exeC.Update(executionFromPb(req.Execution))
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to update execution")
	}

	resp := &pb.UpdateExecutionResp{
		Execution: subserver.PBExecution(execution),
	}
	return resp
}

// helper functions

func executionFromPb(exe *pbcommon.Execution) *execution.Execution {
	return &execution.Execution{
		ID:              exe.Id,
		Date:            exe.Date,
		Location:        common.Location(exe.Location),
		Publish:         exe.Publish,
		CreatedDatetime: getDatetime(exe.CreatedDatetime),
		Culture:         cultureFromPb(exe.Culture),
		Content:         contentFromPb(exe.Content),
		CultureCook:     cultureCookFromPb(exe.CultureCook),
		CultureGuide:    cultureGuideFromPb(exe.CultureGuide),
		Dishes:          dishesFromPb(exe.Dishes),
		Notifications:   notificationsFromPb(exe.Notifications),
		HasPork:         exe.HasPork,
		HasBeef:         exe.HasBeef,
		HasChicken:      exe.HasChicken,
	}
}

func notificationsFromPb(notifications *pbcommon.Notifications) *execution.Notifications {
	return &execution.Notifications{
		DeliverySMS: notifications.DeliverySms,
		RatingSMS:   notifications.RatingSms,
	}
}

func cultureFromPb(culture *pbcommon.Culture) *execution.Culture {
	return &execution.Culture{
		Country:     culture.Country,
		City:        culture.City,
		Description: culture.Description,
		Nationality: culture.Nationality,
		Greeting:    culture.Greeting,
		FlagEmoji:   culture.FlagEmoji,
	}
}

func infoBoxesFromPb(pbd []*pbcommon.InfoBox) []*execution.InfoBox {
	infoBoxes := make([]*execution.InfoBox, len(pbd))
	for i := range pbd {
		infoBoxes[i] = infoBoxFromPb(pbd[i])
	}
	return infoBoxes
}

func infoBoxFromPb(infoBox *pbcommon.InfoBox) *execution.InfoBox {
	return &execution.InfoBox{
		Title:   infoBox.Title,
		Text:    infoBox.Text,
		Caption: infoBox.Caption,
		Image:   infoBox.Image,
	}
}

func cultureGuideFromPb(cultureGuide *pbcommon.CultureGuide) *execution.CultureGuide {
	return &execution.CultureGuide{
		InfoBoxes:          infoBoxesFromPb(cultureGuide.InfoBoxes),
		DinnerInstructions: cultureGuide.DinnerInstructions,
		MainColor:          cultureGuide.MainColor,
		FontName:           cultureGuide.FontName,
		FontStyle:          cultureGuide.FontStyle,
		FontCaps:           cultureGuide.FontCaps,
	}
}

func contentFromPb(content *pbcommon.Content) *execution.Content {
	return &execution.Content{
		HeroImageURL:             content.HeroImageUrl,
		CookImageURL:             content.CookImageUrl,
		HandsPlateNonVegImageURL: content.HandsPlateNonVegImageUrl,
		HandsPlateVegImageURL:    content.HandsPlateVegImageUrl,
		DinnerImageURL:           content.DinnerImageUrl,
		SpotifyURL:               content.SpotifyUrl,
		YoutubeURL:               content.YoutubeUrl,
		FontURL:                  content.FontUrl,
	}
}

func qandasFromPb(pbd []*pbcommon.QandA) []*execution.QandA {
	qandas := make([]*execution.QandA, len(pbd))
	for i := range pbd {
		qandas[i] = qandaFromPb(pbd[i])
	}
	return qandas
}

func qandaFromPb(qanda *pbcommon.QandA) *execution.QandA {
	return &execution.QandA{
		Question: qanda.Question,
		Answer:   qanda.Answer,
	}
}

func cultureCookFromPb(cultureCook *pbcommon.CultureCook) *execution.CultureCook {
	return &execution.CultureCook{
		FirstName:    cultureCook.FirstName,
		LastName:     cultureCook.LastName,
		Story:        cultureCook.Story,
		StoryPreview: cultureCook.StoryPreview,
		QandA:        qandasFromPb(cultureCook.QAndA),
	}
}

func dishesFromPb(pbd []*pbcommon.Dish) []*execution.Dish {
	dishes := make([]*execution.Dish, len(pbd))
	for i := range pbd {
		dishes[i] = dishFromPb(pbd[i])
	}
	return dishes
}

func dishFromPb(dish *pbcommon.Dish) *execution.Dish {
	return &execution.Dish{
		Number:             dish.Number,
		Color:              dish.Color,
		Name:               dish.Name,
		Description:        dish.Description,
		Ingredients:        dish.Ingredients,
		IsForVegetarian:    dish.IsForVegetarian,
		IsForNonVegetarian: dish.IsForNonVegetarian,
		Stickers:           stickersFromPb(dish.Stickers),
	}
}

func stickersFromPb(pbd []*pbcommon.Sticker) []*execution.Sticker {
	stickers := make([]*execution.Sticker, len(pbd))
	for i := range pbd {
		stickers[i] = stickerFromPb(pbd[i])
	}
	return stickers
}

func stickerFromPb(sticker *pbcommon.Sticker) *execution.Sticker {
	return &execution.Sticker{
		Name:                sticker.Name,
		Ingredients:         sticker.Ingredients,
		ExtraInstructions:   sticker.ExtraInstructions,
		ReheatOption1:       sticker.ReheatOption_1,
		ReheatOption2:       sticker.ReheatOption_2,
		ReheatTime1:         sticker.ReheatTime_1,
		ReheatTime2:         sticker.ReheatTime_2,
		ReheatInstructions1: sticker.ReheatInstructions_1,
		ReheatInstructions2: sticker.ReheatInstructions_2,
		EatingTemperature:   sticker.EatingTemperature,
	}
}
