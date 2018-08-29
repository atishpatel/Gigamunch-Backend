package admin

import (
	"context"
	"net/http"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/execution"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// GetAllExecutions gets all executions.
func (s *server) GetAllExecutions(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.GetAllExecutionsReq)
	var err error
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	exeC, err := execution.NewClient(ctx, log)
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get execution client")
	}
	executions, err := exeC.GetAll(int(req.Start), int(req.Limit))
	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get all executions")
	}

	resp := &pb.GetAllExecutionsResp{
		Executions: pbExecutions(executions),
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

	exeC, err := execution.NewClient(ctx, log)
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
func pbExecutions(exes []*execution.Execution) []*pb.Execution {
	pbe := make([]*pb.Execution, len(exes))
	for i := range exes {
		pbe[i] = pbExecution(exes[i])
	}
	return pbe
}

func pbExecution(exe *execution.Execution) *pb.Execution {
	return &pb.Execution{
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
		HasWeirdMeat:    exe.HasWeirdMeat,
		HasFish:         exe.HasFish,
		HasOtherSeafood: exe.HasOtherSeafood,
	}
}

func pbCulture(culture *execution.Culture) *pb.Culture {
	return &pb.Culture{
		Country:     culture.Country,
		City:        culture.City,
		Description: culture.Description,
		Nationality: culture.Nationality,
		Greeting:    culture.Greeting,
		FlagEmoji:   culture.FlagEmoji,
	}
}

func pbContent(content *execution.Content) *pb.Content {
	return &pb.Content{
		HeroImageUrl:       content.HeroImageURL,
		CookImageUrl:       content.CookImageURL,
		HandsPlateImageUrl: content.HandsPlateImageURL,
		DinnerImageUrl:     content.DinnerImageURL,
		SpotifyUrl:         content.SpotifyURL,
		YoutubeUrl:         content.YoutubeURL,
	}
}

func pbCultureCook(cultureCook *execution.CultureCook) *pb.CultureCook {
	return &pb.CultureCook{
		FirstName: cultureCook.FirstName,
		LastName:  cultureCook.LastName,
		Story:     cultureCook.Story,
	}
}

func pbDishes(dishes []execution.Dish) []*pb.Dish {
	pbd := make([]*pb.Dish, len(dishes))
	for i := range dishes {
		pbd[i] = pbDish(dishes[i])
	}
	return pbd
}

func pbDish(dish execution.Dish) *pb.Dish {
	return &pb.Dish{
		Number:             dish.Number,
		Color:              dish.Color,
		Name:               dish.Name,
		Description:        dish.Description,
		Ingredients:        dish.Ingredients,
		IsForVegetarian:    dish.IsForVegetarian,
		IsForNonVegetarian: dish.IsForNonVegetarian,
	}
}

func executionFromPb(exe *pb.Execution) *execution.Execution {
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
		HasWeirdMeat:    exe.HasWeirdMeat,
		HasFish:         exe.HasFish,
		HasOtherSeafood: exe.HasOtherSeafood,
	}
}

func cultureFromPb(culture *pb.Culture) *execution.Culture {
	return &execution.Culture{
		Country:     culture.Country,
		City:        culture.City,
		Description: culture.Description,
		Nationality: culture.Nationality,
		Greeting:    culture.Greeting,
		FlagEmoji:   culture.FlagEmoji,
	}
}

func contentFromPb(content *pb.Content) *execution.Content {
	return &execution.Content{
		HeroImageURL:       content.HeroImageUrl,
		CookImageURL:       content.CookImageUrl,
		HandsPlateImageURL: content.HandsPlateImageUrl,
		DinnerImageURL:     content.DinnerImageUrl,
		SpotifyURL:         content.SpotifyUrl,
		YoutubeURL:         content.YoutubeUrl,
	}
}

func cultureCookFromPb(cultureCook *pb.CultureCook) *execution.CultureCook {
	return &execution.CultureCook{
		FirstName: cultureCook.FirstName,
		LastName:  cultureCook.LastName,
		Story:     cultureCook.Story,
	}
}

func dishesFromPb(pbd []*pb.Dish) []execution.Dish {
	dishes := make([]execution.Dish, len(pbd))
	for i := range pbd {
		dishes[i] = *dishFromPb(pbd[i])
	}
	return dishes
}

func dishFromPb(dish *pb.Dish) *execution.Dish {
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
