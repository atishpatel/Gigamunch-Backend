package main

import (
	"fmt"
	"time"

	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"golang.org/x/net/context"
)

// CohortCell is a cell.
type CohortCell struct {
	AmountLeft       int     `json:"amount_left"`
	AmountLost       int     `json:"amount_lost,omitempty"`
	RetentionPercent float32 `json:"retention_percent"`
}

// CohortRow is a row.
type CohortRow struct {
	Label       string        `json:"label,omitempty"`
	CohortCells []*CohortCell `json:"cohort_cells"`
}

// CohortAnalysis is a full cohortAnalysis.
type CohortAnalysis struct {
	AverageRetention []float32    `json:"average_retention"`
	CohortRows       []*CohortRow `json:"cohort_rows,omitempty"`
	Interval         int16        `json:"interval,omitempty"`
}

// GetGeneralStatsResp is a response for GetGeneralStats.
type GetGeneralStatsResp struct {
	Activities               []*subold.SublogSummary `json:"activities"`
	WeeklyCohortAnalysis     *CohortAnalysis         `json:"weekly_cohort_analysis"`
	WeeklyPaidCohortAnalysis *CohortAnalysis         `json:"weekly_paid_cohort_analysis"`
	MonthlyCohortAnalysis    *CohortAnalysis         `json:"monthly_cohort_analysis"`
	PercentCohortAnalysis    *CohortAnalysis         `json:"percent_cohort_analysis"`
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

	subC := subold.New(ctx)
	// subs, err = subC.GetHasSubscribed(time.Now())
	// if err != nil {
	// 	resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subold.GetHasSubscribed")
	// 	return resp, nil
	// }

	activities, err := subC.GetSublogSummaries()
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subold.GetSublogSummaries")
		return resp, nil
	}
	// resp.Activities = activities
	resp.WeeklyCohortAnalysis = getWeeklyCohort(activities)
	resp.WeeklyPaidCohortAnalysis = getWeeklyPaidCohort(activities)
	resp.MonthlyCohortAnalysis = getMonthlyCohort(activities)
	return resp, nil
}

func getWeeklyCohort(activities []*subold.SublogSummary) *CohortAnalysis {
	analysis := &CohortAnalysis{
		Interval: 7,
	}
	if len(activities) == 0 {
		return analysis
	}
	const labelFormat = "2006-01-02"
	// setup for first row
	lastMinDate := activities[0].MinDate
	row := &CohortRow{
		Label: lastMinDate.Format(labelFormat),
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
				Label: lastMinDate.Format(labelFormat),
			}
			since := time.Since(lastMinDate)
			totalCells := int(since / (time.Duration(analysis.Interval) * time.Hour * 24))
			for z := 0; z < totalCells; z++ {
				row.CohortCells = append(row.CohortCells, new(CohortCell))
			}
		}
		for j := 0; j < activities[i].NumTotal; j++ {
			if len(row.CohortCells)-1 < j {
				row.CohortCells = append(row.CohortCells, new(CohortCell))
			}
			row.CohortCells[j].AmountLeft++
		}
	}
	// fill in amount lost and retention percent
	for _, r := range analysis.CohortRows {
		startAmount := r.CohortCells[0].AmountLeft
		for _, c := range r.CohortCells {
			c.AmountLost = startAmount - c.AmountLeft
			if startAmount > 0 {
				c.RetentionPercent = float32(c.AmountLeft) / float32(startAmount)
			}
		}
	}
	// calculate average retention
	analysis.AverageRetention = make([]float32, len(analysis.CohortRows[0].CohortCells))
	for i := 0; i < len(analysis.AverageRetention); i++ {
		numTotal := 0
		var sumTotal float32
		for _, row := range analysis.CohortRows {
			if i < len(row.CohortCells)-1 {
				sumTotal += row.CohortCells[i].RetentionPercent
				numTotal++
			}
		}
		if numTotal > 0 && sumTotal > 0.01 {
			analysis.AverageRetention[i] = sumTotal / float32(numTotal)
		}
	}
	return analysis
}

func getWeeklyPaidCohort(activities []*subold.SublogSummary) *CohortAnalysis {
	analysis := &CohortAnalysis{
		Interval: 7,
	}
	if len(activities) == 0 {
		return analysis
	}
	const labelFormat = "2006-01-02"
	// setup for first row
	lastMinDate := activities[0].MinDate
	row := &CohortRow{
		Label: lastMinDate.Format(labelFormat),
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
				Label: lastMinDate.Format(labelFormat),
			}
			since := time.Since(lastMinDate)
			totalCells := int(since / (time.Duration(analysis.Interval) * time.Hour * 24))
			for z := 0; z < totalCells; z++ {
				row.CohortCells = append(row.CohortCells, new(CohortCell))
			}
		}
		for j := 0; j < activities[i].NumPaid; j++ {
			if len(row.CohortCells)-1 < j {
				row.CohortCells = append(row.CohortCells, new(CohortCell))
			}
			row.CohortCells[j].AmountLeft++
		}
	}
	// fill in amount lost and retention percent
	for _, r := range analysis.CohortRows {
		startAmount := r.CohortCells[0].AmountLeft
		for _, c := range r.CohortCells {
			c.AmountLost = startAmount - c.AmountLeft
			if startAmount > 0 {
				c.RetentionPercent = float32(c.AmountLeft) / float32(startAmount)
			}
		}
	}
	// calculate average retention
	analysis.AverageRetention = make([]float32, len(analysis.CohortRows[0].CohortCells))
	for i := 0; i < len(analysis.AverageRetention); i++ {
		numTotal := 0
		var sumTotal float32
		for _, row := range analysis.CohortRows {
			if i < len(row.CohortCells)-1 {
				sumTotal += row.CohortCells[i].RetentionPercent
				numTotal++
			}
		}
		if numTotal > 0 && sumTotal > 0.01 {
			analysis.AverageRetention[i] = sumTotal / float32(numTotal)
		}
	}
	return analysis
}

// Assumes the activities are sorted by min date.
func getMonthlyCohort(activities []*subold.SublogSummary) *CohortAnalysis {
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
			since := time.Since(activities[i].MinDate)
			totalCells := int(since / (time.Duration(analysis.Interval) * time.Hour * 24))
			for z := 0; z < totalCells; z++ {
				row.CohortCells = append(row.CohortCells, new(CohortCell))
			}
		}
		for j := 0; j < activities[i].NumTotal; j++ {
			if len(row.CohortCells)-1 < j {
				row.CohortCells = append(row.CohortCells, new(CohortCell))
			}
			row.CohortCells[j].AmountLeft++
		}
	}
	// fill in amount lost and retention percent
	for _, r := range analysis.CohortRows {
		startAmount := r.CohortCells[0].AmountLeft
		for _, c := range r.CohortCells {
			c.AmountLost = startAmount - c.AmountLeft
			if startAmount > 0 {
				c.RetentionPercent = float32(c.AmountLeft) / float32(startAmount)
			}
		}
		r.CohortCells = r.CohortCells[:len(r.CohortCells)-3]
	}
	// calculate average retention
	analysis.AverageRetention = make([]float32, len(analysis.CohortRows[0].CohortCells))
	for i := 0; i < len(analysis.AverageRetention); i++ {
		numTotal := 0
		var sumTotal float32
		for _, row := range analysis.CohortRows {
			if i < len(row.CohortCells)-1 {
				sumTotal += row.CohortCells[i].RetentionPercent
				numTotal++
			}
		}
		if numTotal > 0 && sumTotal > 0.01 {
			analysis.AverageRetention[i] = sumTotal / float32(numTotal)
		}
	}
	return analysis
}
