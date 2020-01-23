package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/execution"
)

var (
	templates  = template.Must(template.New("all").Delims("[[", "]]").ParseGlob("templates/*.html"))
	dateFormat = "2006-01-02"
)

// ExecutionExtended Extends Execution for content generation
type ExecutionExtended struct {
	execution.Execution
	DisplayDate         string
	MondayDisplayDate   string
	ThursdayDisplayDate string
	DishesNonVeg        []execution.Dish
	DishesVeg           []execution.Dish
}

type htmlContentResponse struct {
	previewEmail    string
	cultureEmail    string
	offerEmail      string
	socialMediaText string
}

// Generates preview email html from culture execution json
func getExecutionContent(e *execution.Execution) (*htmlContentResponse, error) {
	exe := &ExecutionExtended{
		Execution: *e,
	}
	date, err := time.Parse(dateFormat, exe.Date)
	if err != nil {
		return nil, err
	}

	exe.DisplayDate = getDisplayDate(date)
	exe.MondayDisplayDate = getDisplayDate(date)
	exe.ThursdayDisplayDate = getDisplayDate(date.AddDate(0, 0, -4))

	// separate nonveg/veg rating links
	links := exe.Notifications.RatingSMS
	linksArray := strings.Split(links, ",")
	if len(linksArray) >= 1 {
		exe.Notifications.RatingLinkNonveg = strings.TrimSpace(linksArray[0])
	}
	if len(linksArray) >= 2 {

		exe.Notifications.RatingLinkVeg = strings.TrimSpace(linksArray[1])
	}

	// sort non veg and veg dishes
	for _, dish := range exe.Dishes {
		if dish.IsForNonVegetarian {
			exe.DishesNonVeg = append(exe.DishesNonVeg, dish)
		}
		if dish.IsForVegetarian {
			exe.DishesVeg = append(exe.DishesVeg, dish)
		}
	}
	// output html files
	resp := &htmlContentResponse{}
	resp.previewEmail, err = generate("preview-email", exe)
	if err != nil {
		return nil, err
	}
	resp.cultureEmail, err = generate("culture-email", exe)
	if err != nil {
		return nil, err
	}
	resp.offerEmail, err = generate("offer-email", exe)
	if err != nil {
		return nil, err
	}
	resp.socialMediaText, err = generate("social-media-text", exe)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getDateSuffix(i int) string {
	j := i % 10
	k := i % 100
	if j == 1 && k != 11 {
		return "st"
	}
	if j == 2 && k != 12 {
		return "nd"
	}
	if j == 3 && k != 13 {
		return "rd"
	}
	return "th"
}

func getDisplayDate(d time.Time) string {
	_, month, day := d.Date()
	weekday := d.Weekday()
	suffix := getDateSuffix(day)
	return fmt.Sprintf("%s, %s %d%s", weekday.String(), month.String(), day, suffix)
}

func generate(inputTemplateName string, exe *ExecutionExtended) (string, error) {
	var err error
	t := templates.Lookup(inputTemplateName)
	buf := new(bytes.Buffer)
	// execute template
	err = t.Execute(buf, exe)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
