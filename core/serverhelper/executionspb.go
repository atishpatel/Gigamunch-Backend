package serverhelper

import (
	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/common"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/execution"
)

// PBExecutions turns an array of executions into a protobuff array of executions.
func PBExecutions(in []*execution.Execution) []*pbcommon.Execution {
	pbe := make([]*pbcommon.Execution, len(in))
	for i := range in {
		pbe[i] = PBExecution(in[i])
	}
	return pbe
}

// PBExecution turns an execution into a protobuff executions.
func PBExecution(in *execution.Execution) *pbcommon.Execution {
	return &pbcommon.Execution{
		ID:              in.ID,
		Date:            in.Date,
		Location:        int32(in.Location),
		Publish:         in.Publish,
		CreatedDatetime: in.CreatedDatetime.String(),
		Culture:         pbCulture(&in.Culture),
		Content:         pbContent(&in.Content),
		CultureCook:     pbCultureCook(&in.CultureCook),
		CultureGuide:    pbCultureGuide(&in.CultureGuide),
		Notifications:   pbNotifications(&in.Notifications),
		Email:           pbEmail(&in.Email),
		Dishes:          pbDishes(in.Dishes),
		Stickers:        pbStickers(in.Stickers),
		HasPork:         in.HasPork,
		HasBeef:         in.HasBeef,
		HasChicken:      in.HasChicken,
	}
}

func pbNotifications(in *execution.Notifications) *pbcommon.Notifications {
	return &pbcommon.Notifications{
		DeliverySMS: in.DeliverySMS,
		RatingSMS:   in.RatingSMS,
	}
}

func pbEmail(in *execution.Email) *pbcommon.Email {
	return &pbcommon.Email{
		DinnerNonVegImageURL: in.DinnerNonVegImageURL,
		DinnerVegImageURL:    in.DinnerVegImageURL,
		CookImageURL:         in.CookImageURL,
		LandscapeImageURL:    in.LandscapeImageURL,
	}
}

func pbInfoBoxes(in []execution.InfoBox) []*pbcommon.InfoBox {
	pbd := make([]*pbcommon.InfoBox, len(in))
	for i := range in {
		pbd[i] = pbInfoBox(&in[i])
	}
	return pbd
}

func pbInfoBox(in *execution.InfoBox) *pbcommon.InfoBox {
	return &pbcommon.InfoBox{
		Title:   in.Title,
		Text:    in.Text,
		Caption: in.Caption,
		Image:   in.Image,
	}
}

func pbCultureGuide(in *execution.CultureGuide) *pbcommon.CultureGuide {
	return &pbcommon.CultureGuide{
		InfoBoxes:                    pbInfoBoxes(in.InfoBoxes),
		DinnerInstructions:           in.DinnerInstructions,
		MainColor:                    in.MainColor,
		FontName:                     in.FontName,
		FontStyle:                    in.FontStyle,
		FontCaps:                     in.FontCaps,
		FontNamePostScript:           in.FontNamePostScript,
		VegetarianDinnerInstructions: in.VegetarianDinnerInstructions,
	}
}

func pbCulture(in *execution.Culture) *pbcommon.Culture {
	return &pbcommon.Culture{
		Country:            in.Country,
		City:               in.City,
		Description:        in.Description,
		DescriptionPreview: in.DescriptionPreview,
		Nationality:        in.Nationality,
		Greeting:           in.Greeting,
		FlagEmoji:          in.FlagEmoji,
	}
}

func pbContent(in *execution.Content) *pbcommon.Content {
	return &pbcommon.Content{
		LandscapeImageURL:        in.LandscapeImageURL,
		CookImageURL:             in.CookImageURL,
		HandsPlateNonVegImageURL: in.HandsPlateNonVegImageURL,
		HandsPlateVegImageURL:    in.HandsPlateVegImageURL,
		DinnerNonVegImageURL:     in.DinnerNonVegImageURL,
		DinnerVegImageURL:        in.DinnerVegImageURL,
		CoverImageURL:            in.CoverImageURL,
		MapImageURL:              in.MapImageURL,
		SpotifyURL:               in.SpotifyURL,
		YoutubeURL:               in.YoutubeURL,
		FontURL:                  in.FontURL,
	}
}

func pbQandAs(in []execution.QandA) []*pbcommon.QandA {
	pbd := make([]*pbcommon.QandA, len(in))
	for i := range in {
		pbd[i] = pbQandA(&in[i])
	}
	return pbd
}

func pbQandA(in *execution.QandA) *pbcommon.QandA {
	return &pbcommon.QandA{
		Question: in.Question,
		Answer:   in.Answer,
	}
}

func pbCultureCook(in *execution.CultureCook) *pbcommon.CultureCook {
	return &pbcommon.CultureCook{
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Story:        in.Story,
		StoryPreview: in.StoryPreview,
		QAndA:        pbQandAs(in.QandA),
	}
}

func pbDishes(in []execution.Dish) []*pbcommon.Dish {
	pbd := make([]*pbcommon.Dish, len(in))
	for i := range in {
		pbd[i] = pbDish(&in[i])
	}
	return pbd
}

func pbDish(in *execution.Dish) *pbcommon.Dish {
	return &pbcommon.Dish{
		Number:             in.Number,
		Color:              in.Color,
		Name:               in.Name,
		Description:        in.Description,
		DescriptionPreview: in.DescriptionPreview,
		Ingredients:        in.Ingredients,
		IsForVegetarian:    in.IsForVegetarian,
		IsForNonVegetarian: in.IsForNonVegetarian,
		IsOnMainPlate:      in.IsOnMainPlate,
		ImageURL:           in.ImageURL,
	}
}

func pbStickers(in []execution.Sticker) []*pbcommon.Sticker {
	pbd := make([]*pbcommon.Sticker, len(in))
	for i := range in {
		pbd[i] = pbSticker(&in[i])
	}
	return pbd
}

func pbSticker(in *execution.Sticker) *pbcommon.Sticker {
	return &pbcommon.Sticker{
		Name:                   in.Name,
		Ingredients:            in.Ingredients,
		ExtraInstructions:      in.ExtraInstructions,
		ReheatOption1:          in.ReheatOption1,
		ReheatOption2:          in.ReheatOption2,
		ReheatTime1:            in.ReheatTime1,
		ReheatTime2:            in.ReheatTime2,
		ReheatInstructions1:    in.ReheatInstructions1,
		ReheatInstructions2:    in.ReheatInstructions2,
		EatingTemperature:      in.EatingTemperature,
		ReheatOption1Preferred: in.ReheatOption1Preferred,
	}
}

// ExecutionFromPb turns pbcommon.Execution to execution.Execution
func ExecutionFromPb(in *pbcommon.Execution) *execution.Execution {
	if in.CultureGuide == nil {
		in.CultureGuide = &pbcommon.CultureGuide{}
	}
	if in.CultureGuide.InfoBoxes == nil {
		in.CultureGuide.InfoBoxes = []*pbcommon.InfoBox{}
	}
	if in.Culture == nil {
		in.Culture = &pbcommon.Culture{}
	}
	if in.CultureCook == nil {
		in.CultureCook = &pbcommon.CultureCook{}
	}
	if in.Content == nil {
		in.Content = &pbcommon.Content{}
	}
	if in.Notifications == nil {
		in.Notifications = &pbcommon.Notifications{}
	}
	return &execution.Execution{
		ID:              in.ID,
		Date:            in.Date,
		Location:        common.Location(in.Location),
		Publish:         in.Publish,
		CreatedDatetime: GetDatetime(in.CreatedDatetime),
		Culture:         *cultureFromPb(in.Culture),
		Content:         *contentFromPb(in.Content),
		CultureCook:     *cultureCookFromPb(in.CultureCook),
		CultureGuide:    *cultureGuideFromPb(in.CultureGuide),
		Dishes:          dishesFromPb(in.Dishes),
		Email:           *emailFromPb(in.Email),
		Stickers:        stickersFromPb(in.Stickers),
		Notifications:   *notificationsFromPb(in.Notifications),
		HasPork:         in.HasPork,
		HasBeef:         in.HasBeef,
		HasChicken:      in.HasChicken,
	}
}

func notificationsFromPb(notifications *pbcommon.Notifications) *execution.Notifications {
	return &execution.Notifications{
		DeliverySMS: notifications.DeliverySMS,
		RatingSMS:   notifications.RatingSMS,
	}
}

func emailFromPb(in *pbcommon.Email) *execution.Email {
	return &execution.Email{
		DinnerNonVegImageURL: in.DinnerNonVegImageURL,
		DinnerVegImageURL:    in.DinnerVegImageURL,
		CookImageURL:         in.CookImageURL,
		LandscapeImageURL:    in.LandscapeImageURL,
	}
}

func cultureFromPb(in *pbcommon.Culture) *execution.Culture {
	return &execution.Culture{
		Country:            in.Country,
		City:               in.City,
		Description:        in.Description,
		DescriptionPreview: in.DescriptionPreview,
		Nationality:        in.Nationality,
		Greeting:           in.Greeting,
		FlagEmoji:          in.FlagEmoji,
	}
}

func infoBoxesFromPb(in []*pbcommon.InfoBox) []execution.InfoBox {
	infoBoxes := make([]execution.InfoBox, len(in))
	for i := range in {
		infoBoxes[i] = *infoBoxFromPb(in[i])
	}
	return infoBoxes
}

func infoBoxFromPb(in *pbcommon.InfoBox) *execution.InfoBox {
	return &execution.InfoBox{
		Title:   in.Title,
		Text:    in.Text,
		Caption: in.Caption,
		Image:   in.Image,
	}
}

func cultureGuideFromPb(in *pbcommon.CultureGuide) *execution.CultureGuide {
	return &execution.CultureGuide{
		InfoBoxes:                    infoBoxesFromPb(in.InfoBoxes),
		DinnerInstructions:           in.DinnerInstructions,
		MainColor:                    in.MainColor,
		FontName:                     in.FontName,
		FontStyle:                    in.FontStyle,
		FontCaps:                     in.FontCaps,
		FontNamePostScript:           in.FontNamePostScript,
		VegetarianDinnerInstructions: in.VegetarianDinnerInstructions,
	}
}

func contentFromPb(in *pbcommon.Content) *execution.Content {
	return &execution.Content{
		LandscapeImageURL:        in.LandscapeImageURL,
		CookImageURL:             in.CookImageURL,
		HandsPlateNonVegImageURL: in.HandsPlateNonVegImageURL,
		HandsPlateVegImageURL:    in.HandsPlateVegImageURL,
		DinnerNonVegImageURL:     in.DinnerNonVegImageURL,
		DinnerVegImageURL:        in.DinnerVegImageURL,
		CoverImageURL:            in.CoverImageURL,
		MapImageURL:              in.MapImageURL,
		SpotifyURL:               in.SpotifyURL,
		YoutubeURL:               in.YoutubeURL,
		FontURL:                  in.FontURL,
	}
}

func qandasFromPb(in []*pbcommon.QandA) []execution.QandA {
	qandas := make([]execution.QandA, len(in))
	for i := range in {
		qandas[i] = *qandaFromPb(in[i])
	}
	return qandas
}

func qandaFromPb(in *pbcommon.QandA) *execution.QandA {
	return &execution.QandA{
		Question: in.Question,
		Answer:   in.Answer,
	}
}

func cultureCookFromPb(in *pbcommon.CultureCook) *execution.CultureCook {
	return &execution.CultureCook{
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Story:        in.Story,
		StoryPreview: in.StoryPreview,
		QandA:        qandasFromPb(in.QAndA),
	}
}

func dishesFromPb(in []*pbcommon.Dish) []execution.Dish {
	dishes := make([]execution.Dish, len(in))
	for i := range in {
		dishes[i] = *dishFromPb(in[i])
	}
	return dishes
}

func dishFromPb(in *pbcommon.Dish) *execution.Dish {
	return &execution.Dish{
		Number:             in.Number,
		Color:              in.Color,
		Name:               in.Name,
		Description:        in.Description,
		DescriptionPreview: in.DescriptionPreview,
		Ingredients:        in.Ingredients,
		IsForVegetarian:    in.IsForVegetarian,
		IsForNonVegetarian: in.IsForNonVegetarian,
		IsOnMainPlate:      in.IsOnMainPlate,
		ImageURL:           in.ImageURL,
	}
}

func stickersFromPb(in []*pbcommon.Sticker) []execution.Sticker {
	if in == nil {
		return []execution.Sticker{}
	}
	stickers := make([]execution.Sticker, len(in))
	for i := range in {
		stickers[i] = *stickerFromPb(in[i])
	}
	return stickers
}

func stickerFromPb(in *pbcommon.Sticker) *execution.Sticker {
	return &execution.Sticker{
		Name:                   in.Name,
		Ingredients:            in.Ingredients,
		ExtraInstructions:      in.ExtraInstructions,
		ReheatOption1:          in.ReheatOption1,
		ReheatOption2:          in.ReheatOption2,
		ReheatTime1:            in.ReheatTime1,
		ReheatTime2:            in.ReheatTime2,
		ReheatInstructions1:    in.ReheatInstructions1,
		ReheatInstructions2:    in.ReheatInstructions2,
		EatingTemperature:      in.EatingTemperature,
		ReheatOption1Preferred: in.ReheatOption1Preferred,
	}
}
