package admin

import (
	"context"
	"net/http"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/common"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/execution"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// GetExecutions gets list of executions.
func (s *server) GetExecutions(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.GetExecutionsReq)
	var err error
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
		Executions: pbExecutions(executions),
	}
	return resp
}

// GetExecution gets an execution.
func (s *server) GetExecution(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.GetExecutionReq)
	var err error
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
		Execution: pbExecution(execution),
	}
	return resp
}

// UpdateExecution updates or creates an execution.
func (s *server) UpdateExecution(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.UpdateExecutionReq)
	var err error
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
		Execution: pbExecution(execution),
	}
	return resp
}

// helper functions
func pbExecutions(exes []*execution.Execution) []*pbcommon.Execution {
	pbe := make([]*pbcommon.Execution, len(exes))
	for i := range exes {
		pbe[i] = pbExecution(exes[i])
	}
	return pbe
}

func pbExecution(exe *execution.Execution) *pbcommon.Execution {
	return &pbcommon.Execution{
		Id:              exe.ID,
		Date:            exe.Date,
		Location:        int32(exe.Location),
		Publish:         exe.Publish,
		CreatedDatetime: exe.CreatedDatetime.String(),
		Culture:         pbCulture(&exe.Culture),
		Content:         pbContent(&exe.Content),
		CultureCook:     pbCultureCook(&exe.CultureCook),
		Dishes:          pbDishes(exe.Dishes),
		HasPork:         exe.HasPork,
		HasBeef:         exe.HasBeef,
		HasChicken:      exe.HasChicken,
		// HasWeirdMeat:    exe.HasWeirdMeat,
		// HasFish:         exe.HasFish,
		// HasOtherSeafood: exe.HasOtherSeafood,
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
	}
}

func pbCultureCook(cultureCook *execution.CultureCook) *pbcommon.CultureCook {
	return &pbcommon.CultureCook{
		FirstName: cultureCook.FirstName,
		LastName:  cultureCook.LastName,
		Story:     cultureCook.Story,
	}
}

func pbDishes(dishes []execution.Dish) []*pbcommon.Dish {
	pbd := make([]*pbcommon.Dish, len(dishes))
	for i := range dishes {
		pbd[i] = pbDish(dishes[i])
	}
	return pbd
}

func pbDish(dish execution.Dish) *pbcommon.Dish {
	return &pbcommon.Dish{
		Number:             dish.Number,
		Color:              dish.Color,
		Name:               dish.Name,
		Description:        dish.Description,
		Ingredients:        dish.Ingredients,
		IsForVegetarian:    dish.IsForVegetarian,
		IsForNonVegetarian: dish.IsForNonVegetarian,
	}
}

func executionFromPb(exe *pbcommon.Execution) *execution.Execution {
	return &execution.Execution{
		ID:              exe.Id,
		Date:            exe.Date,
		Location:        common.Location(exe.Location),
		Publish:         exe.Publish,
		CreatedDatetime: getDatetime(exe.CreatedDatetime),
		Culture:         *cultureFromPb(exe.Culture),
		Content:         *contentFromPb(exe.Content),
		CultureCook:     *cultureCookFromPb(exe.CultureCook),
		Dishes:          dishesFromPb(exe.Dishes),
		HasPork:         exe.HasPork,
		HasBeef:         exe.HasBeef,
		HasChicken:      exe.HasChicken,
		// HasWeirdMeat:    exe.HasWeirdMeat,
		// HasFish:         exe.HasFish,
		// HasOtherSeafood: exe.HasOtherSeafood,
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

func contentFromPb(content *pbcommon.Content) *execution.Content {
	return &execution.Content{
		HeroImageURL:             content.HeroImageUrl,
		CookImageURL:             content.CookImageUrl,
		HandsPlateNonVegImageURL: content.HandsPlateNonVegImageUrl,
		HandsPlateVegImageURL:    content.HandsPlateVegImageUrl,
		DinnerImageURL:           content.DinnerImageUrl,
		SpotifyURL:               content.SpotifyUrl,
		YoutubeURL:               content.YoutubeUrl,
	}
}

func cultureCookFromPb(cultureCook *pbcommon.CultureCook) *execution.CultureCook {
	return &execution.CultureCook{
		FirstName: cultureCook.FirstName,
		LastName:  cultureCook.LastName,
		Story:     cultureCook.Story,
	}
}

func dishesFromPb(pbd []*pbcommon.Dish) []execution.Dish {
	dishes := make([]execution.Dish, len(pbd))
	for i := range pbd {
		dishes[i] = *dishFromPb(pbd[i])
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
	}
}
