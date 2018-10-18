package main

import (
	"time"

	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"golang.org/x/net/context"
)

// CohortCell is a cell.
type CohortCell struct {
	AmountLeft              int     `json:"amount_left"`
	AverageSkipRatio        float32 `json:"average_skip_ratio"`
	AverageRevenuePerPerson float32 `json:"average_revenue_per_person"`
	// Derived from the other values in cell
	// TotalRevenueForGroup = AverageRevenuePerPerson * AmountLeft
	TotalRevenueForGroup float32 `json:"total_revenue_for_group"` // for this cell
	// TotalRevenueSoFar =  AverageRevenuePerPerson * AmountLeft + previousCell.TotalRevenueSoFar
	TotalRevenueSoFar float32 `json:"total_revenue_so_far"` // for this cell and previous cells
	AmountLost        int     `json:"amount_lost,omitempty"`
	RetentionPercent  float32 `json:"retention_percent"`
}

// CohortRow is a row.
type CohortRow struct {
	Label       string        `json:"label,omitempty"`
	CohortCells []*CohortCell `json:"cohort_cells"`
}

// CohortSummary is a summary of the cohort
type CohortSummary struct {
	GroupSize        int     `json:"group_size"` // from AmountLeft
	AverageRetention float32 `json:"average_retention"`
	// Group
	AverageSkipRateForGroup float32 `json:"average_skip_rate_for_group"`
	TotalRevenueForGroup    float32 `json:"total_revenue_for_group"`
	AverageRevenueForGroup  float32 `json:"average_revenue_for_group"`
	// All, "So far"
	StartGroupSize            int     `json:"start_group_size"`
	TotalRevenueSoFarForAll   float32 `json:"total_revenue_so_far_for_all"`
	AverageRevenueSoFarForAll float32 `json:"average_revenue_so_far_for_all"`
}

// CohortAnalysis is a full cohortAnalysis.
type CohortAnalysis struct {
	Label      string           `json:"label"`
	Summary    []*CohortSummary `json:"summary"`
	CohortRows []*CohortRow     `json:"cohort_rows,omitempty"`
	Interval   int16            `json:"interval,omitempty"`
}

// CohortAnalysisSummary is just the summary of a CohortAnalysis.
type CohortAnalysisSummary struct {
	Label   string           `json:"label"`
	Summary []*CohortSummary `json:"summary"`
}

// GetGeneralStatsResp is a response for GetGeneralStats.
type GetGeneralStatsResp struct {
	Activities                  []*subold.SublogSummary `json:"activities"`
	WeeklyCohortAnalysis        *CohortAnalysis         `json:"weekly_cohort_analysis"`
	WeeklyCohortAnalysis2NonVeg *CohortAnalysisSummary  `json:"weekly_cohort_analysis_2_non_veg"`
	WeeklyCohortAnalysis4NonVeg *CohortAnalysisSummary  `json:"weekly_cohort_analysis_4_non_veg"`
	WeeklyCohortAnalysis2Veg    *CohortAnalysisSummary  `json:"weekly_cohort_analysis_2_veg"`
	WeeklyCohortAnalysis4Veg    *CohortAnalysisSummary  `json:"weekly_cohort_analysis_4_veg"`
	WeeklyPaidCohortAnalysis    *CohortAnalysis         `json:"weekly_paid_cohort_analysis"`
	ErrorOnlyResp
}

// GetGeneralStats returns general stats.
func (service *Service) GetGeneralStats(ctx context.Context, req *DateReq) (*GetGeneralStatsResp, error) {
	resp := new(GetGeneralStatsResp)
	defer handleResp(ctx, "GetGeneralStats", resp.Err)
	user, err := getUserFromRequest(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if !user.Admin {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
		return resp, nil
	}

	subC := subold.New(ctx)

	activities, err := subC.GetSublogSummaries(req.Date)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subold.GetSublogSummaries")
		return resp, nil
	}
	resp.Activities = activities
	var activities2NonVeg []*subold.SublogSummary
	var activities4NonVeg []*subold.SublogSummary
	var activities2Veg []*subold.SublogSummary
	var activities4Veg []*subold.SublogSummary
	for _, activity := range activities {
		vegServings := float64(activity.TotalVegServings) / float64(activity.NumTotal)
		nonVegServings := float64(activity.TotalNonVegServings) / float64(activity.NumTotal)
		if vegServings >= 3 {
			activities4Veg = append(activities4Veg, activity)
			continue
		}
		if vegServings > .05 {
			activities2Veg = append(activities2Veg, activity)
			continue
		}
		if nonVegServings >= 3 {
			activities4NonVeg = append(activities4NonVeg, activity)
			continue
		}
		if nonVegServings < 3 {
			activities2NonVeg = append(activities2NonVeg, activity)
			continue
		}
	}
	resp.WeeklyCohortAnalysis = getWeeklyCohort(activities)
	var cohort *CohortAnalysis
	cohort = getWeeklyCohort(activities2NonVeg)
	resp.WeeklyCohortAnalysis2NonVeg = &CohortAnalysisSummary{Label: "2 non-veg " + cohort.Label, Summary: cohort.Summary}
	cohort = getWeeklyCohort(activities2Veg)
	resp.WeeklyCohortAnalysis2Veg = &CohortAnalysisSummary{Label: "2 veg " + cohort.Label, Summary: cohort.Summary}
	cohort = getWeeklyCohort(activities4NonVeg)
	resp.WeeklyCohortAnalysis4NonVeg = &CohortAnalysisSummary{Label: "4 non-veg " + cohort.Label, Summary: cohort.Summary}
	cohort = getWeeklyCohort(activities4Veg)
	resp.WeeklyCohortAnalysis4Veg = &CohortAnalysisSummary{Label: "4 veg " + cohort.Label, Summary: cohort.Summary}
	resp.WeeklyPaidCohortAnalysis = getWeeklyPaidCohort(activities)
	return resp, nil
}

func getWeeklyCohort(activities []*subold.SublogSummary) *CohortAnalysis {
	analysis := &CohortAnalysis{
		Interval: 7,
	}
	if len(activities) < 3 {
		return analysis
	}
	const labelFormat = "2006-01-02"
	maxCohortCellLength := 0
	// setup for first row
	lastMinDate := activities[0].MinDate
	row := &CohortRow{
		Label: lastMinDate.Format(labelFormat),
	}
	analysis.Label = lastMinDate.Format(labelFormat)
	// for every activity
	for i := 0; i < len(activities); i++ {
		// next cohort row
		diff := activities[i].MinDate.Sub(lastMinDate)
		if diff < time.Duration(0) {
			diff *= -1
		}
		if diff > 12*time.Hour {
			// if different mindate, add next row
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
		subSkipAverage := (float32(activities[i].NumSkip) / float32(activities[i].NumTotal)) * 100
		avgRevenueForSubPerWeek := float32(activities[i].TotalAmountPaid) / float32(activities[i].NumTotal)
		// num total times sub got it = amount left over the weeks
		for j := 0; j < activities[i].NumTotal; j++ {
			if len(row.CohortCells)-1 < j {
				row.CohortCells = append(row.CohortCells, new(CohortCell))
			}
			if j > maxCohortCellLength {
				maxCohortCellLength = j
			}
			row.CohortCells[j].AverageRevenuePerPerson = ((row.CohortCells[j].AverageRevenuePerPerson * float32(row.CohortCells[j].AmountLeft)) + avgRevenueForSubPerWeek) / float32(row.CohortCells[j].AmountLeft+1)
			row.CohortCells[j].AverageSkipRatio = ((row.CohortCells[j].AverageSkipRatio * float32(row.CohortCells[j].AmountLeft)) + subSkipAverage) / float32(row.CohortCells[j].AmountLeft+1)
			row.CohortCells[j].AmountLeft++
		}
	}
	// fill in amount lost and retention percent
	for _, r := range analysis.CohortRows {
		startAmount := r.CohortCells[0].AmountLeft
		for j, c := range r.CohortCells {
			c.AmountLost = startAmount - c.AmountLeft
			if startAmount > 0 {
				c.RetentionPercent = float32(c.AmountLeft) / float32(startAmount)
				c.TotalRevenueForGroup = float32(c.AmountLeft) * float32(c.AverageRevenuePerPerson) * float32(j)
				var previousRevenue float32
				if j > 0 {
					previousRevenue = r.CohortCells[j-1].TotalRevenueSoFar
				}
				c.TotalRevenueSoFar += float32(c.AmountLeft)*float32(c.AverageRevenuePerPerson) + previousRevenue
			}
		}
	}
	// calculate average retention
	analysis.Summary = make([]*CohortSummary, maxCohortCellLength)
	for i := 0; i < len(analysis.Summary); i++ {
		analysis.Summary[i] = new(CohortSummary)
		numTotal := 0
		var sumTotalRetention float32
		var sumTotalRevenueSoFar float32
		var sumTotalRevenueForGroup float32
		var sumTotalAverageSkip float32
		for _, row := range analysis.CohortRows {
			if i < len(row.CohortCells)-1 {
				sumTotalRetention += row.CohortCells[i].RetentionPercent
				sumTotalRevenueSoFar += row.CohortCells[i].TotalRevenueSoFar
				sumTotalRevenueForGroup += row.CohortCells[i].TotalRevenueForGroup
				sumTotalAverageSkip += row.CohortCells[i].AverageSkipRatio
				numTotal++
				analysis.Summary[i].GroupSize += row.CohortCells[i].AmountLeft
			}
		}
		if numTotal > 0 && sumTotalRetention > 0.01 {
			analysis.Summary[i].StartGroupSize = len(activities)
			analysis.Summary[i].AverageRetention = sumTotalRetention / float32(numTotal)
			analysis.Summary[i].TotalRevenueForGroup = sumTotalRevenueForGroup
			analysis.Summary[i].AverageRevenueForGroup = sumTotalRevenueForGroup / float32(numTotal)
			analysis.Summary[i].AverageSkipRateForGroup = sumTotalAverageSkip / float32(numTotal)
			analysis.Summary[i].TotalRevenueSoFarForAll = sumTotalRevenueSoFar
			analysis.Summary[i].AverageRevenueSoFarForAll = sumTotalRevenueSoFar / float32(numTotal)
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
	analysis.Summary = make([]*CohortSummary, len(analysis.CohortRows[0].CohortCells))
	for i := 0; i < len(analysis.Summary); i++ {
		analysis.Summary[i] = new(CohortSummary)
		numTotal := 0
		var sumTotal float32
		for _, row := range analysis.CohortRows {
			if i < len(row.CohortCells)-1 {
				sumTotal += row.CohortCells[i].RetentionPercent
				numTotal++
			}
		}
		if numTotal > 0 && sumTotal > 0.01 {
			analysis.Summary[i].AverageRetention = sumTotal / float32(numTotal)
		}
	}
	return analysis
}
