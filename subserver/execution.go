package subserver

import (
	"context"
	"net/http"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/common"
	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/sub"

	"github.com/atishpatel/Gigamunch-Backend/core/execution"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// GetExecutions gets list of executions.
func (s *Server) GetExecutions(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pbsub.GetExecutionsReq)
	var err error
	// decode request
	err = DecodeRequest(ctx, r, req)
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

	resp := &pbsub.GetExecutionsResp{
		Executions: PBExecutions(executions),
	}
	return resp
}

// GetExecution gets an execution.
func (s *Server) GetExecution(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pbsub.GetExecutionReq)
	var err error
	// decode request
	err = DecodeRequest(ctx, r, req)
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

	resp := &pbsub.GetExecutionResp{
		Execution: PBExecution(execution),
	}
	return resp
}

// helper functions

// PBExecutions turns an array of executions into a protobuff array of executions.
func PBExecutions(exes []*execution.Execution) []*pbcommon.Execution {
	pbe := make([]*pbcommon.Execution, len(exes))
	for i := range exes {
		pbe[i] = PBExecution(exes[i])
	}
	return pbe
}

// PBExecution turns an execution into a protobuff executions.
func PBExecution(exe *execution.Execution) *pbcommon.Execution {
	return &pbcommon.Execution{
		Id:              exe.ID,
		Date:            exe.Date,
		Location:        int32(exe.Location),
		Publish:         exe.Publish,
		CreatedDatetime: exe.CreatedDatetime.String(),
		Culture:         pbCulture(exe.Culture),
		Content:         pbContent(exe.Content),
		CultureCook:     pbCultureCook(exe.CultureCook),
		CultureGuide:    pbCultureGuide(exe.CultureGuide),
		Notifications:   pbNotifications(exe.Notifications),
		Dishes:          pbDishes(exe.Dishes),
		HasPork:         exe.HasPork,
		HasBeef:         exe.HasBeef,
		HasChicken:      exe.HasChicken,
	}
}

func pbNotifications(notifications *execution.Notifications) *pbcommon.Notifications {
	return &pbcommon.Notifications{
		DeliverySms: notifications.DeliverySMS,
		RatingSms:   notifications.RatingSMS,
	}
}

func pbInfoBoxes(infoBoxes []*execution.InfoBox) []*pbcommon.InfoBox {
	pbd := make([]*pbcommon.InfoBox, len(infoBoxes))
	for i := range infoBoxes {
		pbd[i] = pbInfoBox(infoBoxes[i])
	}
	return pbd
}

func pbInfoBox(infoBox *execution.InfoBox) *pbcommon.InfoBox {
	return &pbcommon.InfoBox{
		Title:   infoBox.Title,
		Text:    infoBox.Text,
		Caption: infoBox.Caption,
		Image:   infoBox.Image,
	}
}

func pbCultureGuide(cultureGuide *execution.CultureGuide) *pbcommon.CultureGuide {
	return &pbcommon.CultureGuide{
		InfoBoxes:          pbInfoBoxes(cultureGuide.InfoBoxes),
		DinnerInstructions: cultureGuide.DinnerInstructions,
		MainColor:          cultureGuide.MainColor,
		FontName:           cultureGuide.FontName,
		FontStyle:          cultureGuide.FontStyle,
		FontCaps:           cultureGuide.FontCaps,
	}
}

func pbCulture(culture *execution.Culture) *pbcommon.Culture {
	return &pbcommon.Culture{
		Country:     culture.Country,
		City:        culture.City,
		Description: culture.Description,
		Nationality: culture.Nationality,
		Greeting:    culture.Greeting,
		FlagEmoji:   culture.FlagEmoji,
	}
}

func pbContent(content *execution.Content) *pbcommon.Content {
	return &pbcommon.Content{
		HeroImageUrl:             content.HeroImageURL,
		CookImageUrl:             content.CookImageURL,
		HandsPlateNonVegImageUrl: content.HandsPlateNonVegImageURL,
		HandsPlateVegImageUrl:    content.HandsPlateVegImageURL,
		DinnerImageUrl:           content.DinnerImageURL,
		SpotifyUrl:               content.SpotifyURL,
		YoutubeUrl:               content.YoutubeURL,
		FontUrl:                  content.FontURL,
	}
}

func pbQandAs(qanda []*execution.QandA) []*pbcommon.QandA {
	pbd := make([]*pbcommon.QandA, len(qanda))
	for i := range qanda {
		pbd[i] = pbQandA(qanda[i])
	}
	return pbd
}

func pbQandA(qanda *execution.QandA) *pbcommon.QandA {
	return &pbcommon.QandA{
		Question: qanda.Question,
		Answer:   qanda.Answer,
	}
}

func pbCultureCook(cultureCook *execution.CultureCook) *pbcommon.CultureCook {
	return &pbcommon.CultureCook{
		FirstName:    cultureCook.FirstName,
		LastName:     cultureCook.LastName,
		Story:        cultureCook.Story,
		StoryPreview: cultureCook.StoryPreview,
		QAndA:        pbQandAs(cultureCook.QandA),
	}
}

func pbDishes(dishes []*execution.Dish) []*pbcommon.Dish {
	pbd := make([]*pbcommon.Dish, len(dishes))
	for i := range dishes {
		pbd[i] = pbDish(dishes[i])
	}
	return pbd
}

func pbDish(dish *execution.Dish) *pbcommon.Dish {
	return &pbcommon.Dish{
		Number:             dish.Number,
		Color:              dish.Color,
		Name:               dish.Name,
		Description:        dish.Description,
		Ingredients:        dish.Ingredients,
		IsForVegetarian:    dish.IsForVegetarian,
		IsForNonVegetarian: dish.IsForNonVegetarian,
		Stickers:           pbStickers(dish.Stickers),
	}
}

func pbStickers(stickers []*execution.Sticker) []*pbcommon.Sticker {
	pbd := make([]*pbcommon.Sticker, len(stickers))
	for i := range stickers {
		pbd[i] = pbSticker(stickers[i])
	}
	return pbd
}

func pbSticker(sticker *execution.Sticker) *pbcommon.Sticker {
	return &pbcommon.Sticker{
		Name:                 sticker.Name,
		Ingredients:          sticker.Ingredients,
		ExtraInstructions:    sticker.ExtraInstructions,
		ReheatOption_1:       sticker.ReheatOption1,
		ReheatOption_2:       sticker.ReheatOption2,
		ReheatTime_1:         sticker.ReheatTime1,
		ReheatTime_2:         sticker.ReheatTime2,
		ReheatInstructions_1: sticker.ReheatInstructions1,
		ReheatInstructions_2: sticker.ReheatInstructions2,
		EatingTemperature:    sticker.EatingTemperature,
	}
}
