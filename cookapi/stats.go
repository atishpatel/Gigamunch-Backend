package main

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"golang.org/x/net/context"
)

// CohortCell is a cell.
type CohortCell struct {
	AmountLeft       int     `json:"amount_left,omitempty"`
	AmountLost       int     `json:"amount_lost,omitempty"`
	RetentionPercent float32 `json:"retention_percent,omitempty"`
}

// CohortRow is a row.
type CohortRow struct {
	Label      string        `json:"label,omitempty"`
	CohortCell []*CohortCell `json:"cohort_cell,omitempty"`
}

// CohortAnalysis is a full cohortAnalysis.
type CohortAnalysis struct {
	CohortRows []*CohortRow `json:"cohort_rows,omitempty"`
	Interval   int16        `json:"interval,omitempty"`
}

// GetGeneralStatsResp is a response for GetGeneralStats.
type GetGeneralStatsResp struct {
	Activities            []*sub.SublogSummary `json:"activities"`
	WeeklyCohortAnalysis  *CohortAnalysis      `json:"weekly_cohort_analysis"`
	MonthlyCohortAnalysis *CohortAnalysis      `json:"monthly_cohort_analysis"`
	PercentCohortAnalysis *CohortAnalysis      `json:"percent_cohort_analysis"`
	ErrorOnlyResp
}

// GetGeneralStats returns general stats.
func (service *Service) GetGeneralStats(ctx context.Context, req *GigatokenReq) (*GetGeneralStatsResp, error) {
	resp := new(GetGeneralStatsResp)
	defer handleResp(ctx, "GetGeneralStats", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.IsAdmin() {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	subC := sub.New(ctx)
	// subs, err = subC.GetHasSubscribed(time.Now())
	// if err != nil {
	// 	resp.Err = errors.GetErrorWithCode(err).Wrap("failed to sub.GetHasSubscribed")
	// 	return resp, nil
	// }

	activities, err := subC.GetSublogSummaries()
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to sub.GetSublogSummaries")
		return resp, nil
	}
	// resp.Activities = activities
	resp.WeeklyCohortAnalysis = getWeeklyCohort(activities)
	resp.MonthlyCohortAnalysis = getMonthlyCohort(activities)
	return resp, nil
}

func getWeeklyCohort(activities []*sub.SublogSummary) *CohortAnalysis {
	analysis := &CohortAnalysis{
		Interval: 7,
	}
	if len(activities) == 0 {
		return analysis
	}
	// setup for first row
	lastMinDate := activities[0].MinDate
	row := &CohortRow{
		Label: lastMinDate.String(),
	}
	// for every activity
	for i := 0; i < len(activities); i++ {
		// next cohort row
		diff := activities[i].MinDate.Sub(lastMinDate)
		if diff < time.Duration(0) {
			diff *= -1
		}
		if diff > 12*time.Hour {
			analysis.CohortRows = append(analysis.CohortRows, row)
			lastMinDate = activities[i].MinDate
			row = &CohortRow{
				Label: lastMinDate.String(),
			}
		}
		for j := 0; j < activities[i].NumTotal; j++ {
			if len(row.CohortCell)-1 < j {
				row.CohortCell = append(row.CohortCell, new(CohortCell))
			}
			row.CohortCell[j].AmountLeft++
		}
	}
	// fill in amount lost and retention percent
	for _, r := range analysis.CohortRows {
		startAmount := r.CohortCell[0].AmountLeft
		for _, c := range r.CohortCell {
			c.AmountLost = startAmount - c.AmountLeft
			c.RetentionPercent = float32(c.AmountLeft) / float32(startAmount)
		}
	}
	return analysis
}

// Assumes the activities are sorted by min date.
func getMonthlyCohort(activities []*sub.SublogSummary) *CohortAnalysis {
	analysis := &CohortAnalysis{
		Interval: 30,
	}
	if len(activities) == 0 {
		return analysis
	}
	const labelFormat = "%s-%d"
	// setup for first row
	lastMonthYear := fmt.Sprintf(labelFormat, activities[0].MinDate.Month().String(), activities[0].MinDate.Year())
	row := &CohortRow{
		Label: lastMonthYear,
	}
	// for every activity
	for i := 0; i < len(activities); i++ {
		// next cohort row
		monthYear := fmt.Sprintf(labelFormat, activities[i].MinDate.Month().String(), activities[i].MinDate.Year())
		if monthYear != lastMonthYear {
			analysis.CohortRows = append(analysis.CohortRows, row)
			lastMonthYear = monthYear
			row = &CohortRow{
				Label: lastMonthYear,
			}
		}
		for j := 0; j < activities[i].NumTotal; j++ {
			if len(row.CohortCell)-1 < j {
				row.CohortCell = append(row.CohortCell, new(CohortCell))
			}
			row.CohortCell[j].AmountLeft++
		}
	}
	// fill in amount lost and retention percent
	for _, r := range analysis.CohortRows {
		startAmount := r.CohortCell[0].AmountLeft
		for _, c := range r.CohortCell {
			c.AmountLost = startAmount - c.AmountLeft
			c.RetentionPercent = float32(c.AmountLeft) / float32(startAmount)
		}
	}
	return analysis
}
