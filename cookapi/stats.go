package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/utils"

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

type ChurnWithWeeksRetained struct {
	Label         string  `json:"label"`
	Operator      string  `json:"operator"`
	WeeksRetained int     `json:"weeks_retained"`
	Churn         float32 `json:"churn"`
	CancelNum     int     `json:"cancel_num"`
}

type ChurnSummary struct {
	Label       string                    `json:"label"`
	Month       time.Month                `json:"month"`
	Year        int                       `json:"year"`
	MonthStart  time.Time                 `json:"month_start"`
	MonthEnd    time.Time                 `json:"month_end"`
	SubsAtTime  int                       `json:"subs_at_time"`
	ChurnGroups []*ChurnWithWeeksRetained `json:"churn_groups"`
}

type ChurnAnalysisSummary struct {
	ChurnList []*ChurnSummary `json:"churn_list"`
	Interval  int16           `json:"interval,omitempty"`
}

type LTVHistogram struct {
	AverageWeeks            float32   `json:"average_weeks,omitempty"`
	AveragePaidWeeks        float32   `json:"average_paid_weeks,omitempty"`
	AveragePaidRevenue      float32   `json:"average_paid_revenue,omitempty"`
	Percentile50PaidRevenue float32   `json:"percentile_50_paid_revenue,omitempty"`
	Percentile50PaidWeeks   float32   `json:"percentile_50_paid_weeks,omitempty"`
	Weeks                   []int     `json:"weeks,omitempty"`
	PaidWeeks               []int     `json:"paid_weeks,omitempty"`
	PaidRevenue             []float32 `json:"paid_revenue,omitempty"`
}

type LifeTimeValueSummary struct {
	AverageChurn       float32       `json:"average_churn,omitempty"`
	ProjectedHistogram *LTVHistogram `json:"projected_histogram,omitempty"`
	ActualHistogram    *LTVHistogram `json:"actual_histogram,omitempty"`
	CanceledHistogram  *LTVHistogram `json:"canceled_histogram,omitempty"`
}

type BagTypeBreakDown struct {
	NonVeg2 int `json:"non_veg_2"`
	Veg2    int `json:"veg_2"`
	NonVeg4 int `json:"non_veg_4"`
	Veg4    int `json:"veg_4"`
}

type BagPriceBreakDown struct {
	Price       float32 `json:"price"`
	Count       int     `json:"count"`
	ActiveCount int     `json:"active_count"`
}

// GetGeneralStatsResp is a response for GetGeneralStats.
type GetGeneralStatsResp struct {
	Activities []*subold.SublogSummary `json:"activities"`
	// ProjectedActivities           []*subold.SublogSummary `json:"projected_activities"`
	WeeklyCohortAnalysis          *CohortAnalysis        `json:"weekly_cohort_analysis"`
	ProjectedWeeklyCohortAnalysis *CohortAnalysisSummary `json:"projected_weekly_cohort_analysis"`
	WeeklyCohortAnalysis2NonVeg   *CohortAnalysisSummary `json:"weekly_cohort_analysis_2_non_veg"`
	WeeklyCohortAnalysis4NonVeg   *CohortAnalysisSummary `json:"weekly_cohort_analysis_4_non_veg"`
	WeeklyCohortAnalysis2Veg      *CohortAnalysisSummary `json:"weekly_cohort_analysis_2_veg"`
	WeeklyCohortAnalysis4Veg      *CohortAnalysisSummary `json:"weekly_cohort_analysis_4_veg"`
	BagTypeBreakDown              *BagTypeBreakDown      `json:"bag_type_break_down"`
	BagTypeBreakDownActive        *BagTypeBreakDown      `json:"bag_type_break_down_active"`
	BagPriceBreakDown             []*BagPriceBreakDown   `json:"bag_price_break_down"`
	LifeTimeValue                 *LifeTimeValueSummary  `json:"life_time_value"`
	MonthlyChurn                  *ChurnAnalysisSummary  `json:"monthly_churn"`
	ErrorOnlyResp
}

type GetGeneralStatsReq struct {
	GigatokenReq
	StartDateMin time.Time `json:"start_date_min"`
	StartDateMax time.Time `json:"start_date_max"`
}

// GetGeneralStats returns general stats.
func (service *Service) GetGeneralStats(ctx context.Context, req *GetGeneralStatsReq) (*GetGeneralStatsResp, error) {
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

	activities, err := subC.GetSublogSummaries(req.StartDateMin, req.StartDateMax)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to subold.GetSublogSummaries")
		return resp, nil
	}
	resp.Activities = activities
	resp.WeeklyCohortAnalysis = getWeeklyCohort(activities)
	resp.MonthlyChurn = getMonthlyChurn(activities)
	churn := getChurn(resp.WeeklyCohortAnalysis.Summary)
	projectedActivites := projectActivites(ctx, activities, churn)
	var activities2NonVeg []*subold.SublogSummary
	var activities4NonVeg []*subold.SublogSummary
	var activities2Veg []*subold.SublogSummary
	var activities4Veg []*subold.SublogSummary
	for _, activity := range projectedActivites {
		vegServings := float64(activity.TotalVegServings) / float64(activity.NumTotal)
		nonVegServings := float64(activity.TotalNonVegServings) / float64(activity.NumTotal)
		if vegServings >= 3 {
			activities4Veg = append(activities4Veg, activity)
			continue
		}
		if vegServings > .01 {
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
	var cohort *CohortAnalysis
	cohort = getWeeklyCohort(activities2NonVeg)
	resp.WeeklyCohortAnalysis2NonVeg = &CohortAnalysisSummary{Label: "2 non-veg " + cohort.Label, Summary: cohort.Summary}
	cohort = getWeeklyCohort(activities2Veg)
	resp.WeeklyCohortAnalysis2Veg = &CohortAnalysisSummary{Label: "2 veg " + cohort.Label, Summary: cohort.Summary}
	cohort = getWeeklyCohort(activities4NonVeg)
	resp.WeeklyCohortAnalysis4NonVeg = &CohortAnalysisSummary{Label: "4 non-veg " + cohort.Label, Summary: cohort.Summary}
	cohort = getWeeklyCohort(activities4Veg)
	resp.WeeklyCohortAnalysis4Veg = &CohortAnalysisSummary{Label: "4 veg " + cohort.Label, Summary: cohort.Summary}
	resp.LifeTimeValue = getLTVSummary(ctx, activities, projectedActivites, churn)
	// resp.ProjectedActivities = projectedActivites
	cohort = getWeeklyCohort(projectedActivites)
	resp.ProjectedWeeklyCohortAnalysis = &CohortAnalysisSummary{Label: "Projected " + cohort.Label, Summary: cohort.Summary}
	resp.BagTypeBreakDown, resp.BagTypeBreakDownActive = getBagTypeBreakDown(activities)
	resp.BagPriceBreakDown = getBagPriceBreadDown(activities)
	return resp, nil
}

func getWeeklyCohort(activities []*subold.SublogSummary) *CohortAnalysis {
	analysis := &CohortAnalysis{
		Interval: 7,
	}
	if len(activities) < 10 {
		return analysis
	}
	const labelFormat = "2006-01-02"
	maxCohortCellLength := 0
	// setup for first row
	lastMinDate := activities[0].MinDate
	row := &CohortRow{
		Label: lastMinDate.Format(labelFormat),
	}
	analysis.Label = lastMinDate.Format(labelFormat) + " to " + activities[len(activities)-1].MaxDate.Format(labelFormat)
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

func getMonthlyChurn(activities []*subold.SublogSummary) *ChurnAnalysisSummary {
	analysis := &ChurnAnalysisSummary{
		Interval: 30,
	}
	if len(activities) == 0 {
		return analysis
	}
	const labelFormat = "%s-%d"
	const dateFormat = "2006-01-02"
	// fill out array
	timeIndex := activities[0].MinDate
	now := time.Now()

	getMonthString := func(m time.Month) string {
		if m < 10 {
			return fmt.Sprintf("0%d", m)
		}
		return fmt.Sprintf("%d", m)
	}

	for timeIndex.Before(now) {
		monthStart, _ := time.Parse(dateFormat, fmt.Sprintf("%d-%s-01", timeIndex.Year(), getMonthString(timeIndex.Month())))
		nextMonth := timeIndex.Month() + 1
		nextYear := timeIndex.Year()
		if int(nextMonth) > 12 {
			nextYear++
			nextMonth = time.January
		}
		tmp, _ := time.Parse(dateFormat, fmt.Sprintf("%d-%s-01", nextYear, getMonthString(nextMonth)))
		monthEnd := tmp.Add(-1 * 24 * time.Hour)

		row := &ChurnSummary{
			Label:      fmt.Sprintf(labelFormat, timeIndex.Month().String(), timeIndex.Year()),
			Month:      timeIndex.Month(),
			Year:       timeIndex.Year(),
			MonthStart: monthStart,
			MonthEnd:   monthEnd,
		}

		row.ChurnGroups = []*ChurnWithWeeksRetained{
			&ChurnWithWeeksRetained{
				Label:         ">0",
				Operator:      ">",
				WeeksRetained: 0,
			},
			&ChurnWithWeeksRetained{
				Label:         "<=4",
				Operator:      "<",
				WeeksRetained: 5,
			},
			&ChurnWithWeeksRetained{
				Label:         ">4",
				Operator:      ">",
				WeeksRetained: 4,
			},
			&ChurnWithWeeksRetained{
				Label:         "<=12",
				Operator:      "<",
				WeeksRetained: 13,
			},
			&ChurnWithWeeksRetained{
				Label:         ">12",
				Operator:      ">",
				WeeksRetained: 12,
			},
			&ChurnWithWeeksRetained{
				Label:         "<=20",
				Operator:      "<",
				WeeksRetained: 21,
			},
			&ChurnWithWeeksRetained{
				Label:         ">20",
				Operator:      ">",
				WeeksRetained: 20,
			},
		}

		analysis.ChurnList = append(analysis.ChurnList, row)
		// next month
		m := timeIndex.Month()
		for timeIndex.Month() == m {
			timeIndex = timeIndex.Add(time.Hour * 7 * 24)
		}
	}

	for _, activity := range activities {
		for _, monthAnalysis := range analysis.ChurnList {
			if activity.MinDate.Before(monthAnalysis.MonthStart) && activity.MaxDate.After(monthAnalysis.MonthStart) {
				// was sub at start of month
				monthAnalysis.SubsAtTime++
				if monthAnalysis.MonthStart.After(activity.MinDate) && monthAnalysis.MonthStart.Before(activity.MaxDate) && monthAnalysis.MonthEnd.After(activity.MaxDate) {
					// canceled between this month
					for _, a := range monthAnalysis.ChurnGroups {
						if a.Operator == "" {
							a.CancelNum++
						} else if a.Operator == "<" {
							if activity.NumTotal < a.WeeksRetained {
								a.CancelNum++
							}
						} else if a.Operator == ">" {
							if activity.NumTotal > a.WeeksRetained {
								a.CancelNum++
							}
						}
					}
				}
			}
		}
	}

	for _, month := range analysis.ChurnList {
		for _, churnGroup := range month.ChurnGroups {
			if month.SubsAtTime > 0 {
				churnGroup.Churn = (float32(churnGroup.CancelNum) / float32(month.SubsAtTime)) * 100.0
			}
		}
	}

	return analysis
}

func getChurn(summaries []*CohortSummary) float32 {
	var totalChurn float32
	var totalCount float32
	for i := 0; i < len(summaries)-2; i++ {
		weekChurn := summaries[i].AverageRetention - summaries[i+1].AverageRetention
		totalChurn += weekChurn
		totalCount++
		if i >= 10 || weekChurn < .01 {
			break
		}
	}
	churn := totalChurn / totalCount
	if churn < .01 {
		churn = .0375
	}
	return churn
}

func projectActivites(ctx context.Context, activities []*subold.SublogSummary, churn float32) []*subold.SublogSummary {
	projectedActivites := []*subold.SublogSummary{}
	stillSubscribedActivites := []*subold.SublogSummary{}
	for _, activity := range activities {
		if time.Since(activity.MaxDate) < 7*24*time.Hour {
			// still subscribed
			stillSubscribedActivites = append(stillSubscribedActivites, activity)
		} else {
			// already canceled
			projectedActivites = append(projectedActivites, activity)
		}
	}
	utils.Infof(ctx, "churn: %3.3f", churn)
	utils.Infof(ctx, "still subscribed actvities length: ", len(stillSubscribedActivites))
	// get remaining subscribers over time
	subOverTime := []int{len(stillSubscribedActivites)}
	for subOverTime[len(subOverTime)-1] != 0 {
		remainingSubs := int(float32(subOverTime[len(subOverTime)-1]) * (1.0 - churn))
		subOverTime = append(subOverTime, remainingSubs)
	}
	utils.Infof(ctx, "subOverTime: %+v", subOverTime)
	// distribute remaining subscribers over time to
	stillSubscribedIndex := 0
	for i := 0; i < len(subOverTime)-1; i++ {
		lostSubs := subOverTime[i] - subOverTime[i+1]
		for j := 0; j < lostSubs; j++ {
			if stillSubscribedIndex >= len(stillSubscribedActivites) {
				break
			}
			// say they canceled after x weeks
			actualSub := stillSubscribedActivites[stillSubscribedIndex]
			projectedSub := *actualSub // make copy
			// act as if they stayed for this length of time
			projectedSub.MaxDate = projectedSub.MaxDate.Add(time.Duration(i) * 7 * 24 * time.Hour)
			skipRate := float32(projectedSub.NumSkip) / float32(projectedSub.NumTotal)
			costPerPaid := projectedSub.TotalAmountPaid / float32(projectedSub.NumPaid)
			if projectedSub.NumPaid == 0 || projectedSub.TotalAmountPaid < 0.1 {
				costPerPaid = projectedSub.TotalAmount / float32(projectedSub.NumTotal)
			}
			projectedSkips := int(skipRate * float32(i))
			projectedPaid := int((1.0 - skipRate) * float32(i))
			projectedSub.NumSkip += projectedSkips
			projectedSub.NumPaid += projectedPaid
			projectedSub.NumTotal += projectedSkips + projectedPaid
			if projectedSub.TotalVegServings > 2 {
				projectedSub.TotalVegServings += projectedSkips + projectedPaid
			} else {
				projectedSub.TotalNonVegServings += projectedSkips + projectedPaid
			}
			projectedSub.TotalAmountPaid += costPerPaid * float32(projectedPaid)
			projectedSub.TotalAmount += costPerPaid * float32(projectedPaid+projectedSkips)
			projectedActivites = append(projectedActivites, &projectedSub)
			stillSubscribedIndex++
		}
	}

	sort.Slice(projectedActivites, func(i, j int) bool {
		return projectedActivites[i].MinDate.Before(projectedActivites[j].MinDate)
	})

	return projectedActivites
}

func getLTVSummary(ctx context.Context, activities, projectedActivites []*subold.SublogSummary, churn float32) *LifeTimeValueSummary {
	ltvSummary := &LifeTimeValueSummary{}
	ltvSummary.AverageChurn = churn

	canceledActivities := []*subold.SublogSummary{}
	for _, activity := range activities {
		if time.Since(activity.MaxDate) > 7*24*time.Hour {
			// already canceled
			canceledActivities = append(canceledActivities, activity)
		}
	}

	ltvSummary.ActualHistogram = getLTVHistogram(ctx, activities)
	ltvSummary.ProjectedHistogram = getLTVHistogram(ctx, projectedActivites)
	ltvSummary.CanceledHistogram = getLTVHistogram(ctx, canceledActivities)
	return ltvSummary
}

func getLTVHistogram(ctx context.Context, activities []*subold.SublogSummary) *LTVHistogram {
	numSubs := len(activities)
	ltvHistogram := &LTVHistogram{
		Weeks:       make([]int, numSubs),
		PaidWeeks:   make([]int, numSubs),
		PaidRevenue: make([]float32, numSubs),
	}
	if numSubs < 2 {
		return ltvHistogram
	}

	totalWeeks := 0
	totalPaidWeeks := 0
	var totalPaidRevenue float32
	for i, activity := range activities {
		ltvHistogram.Weeks[i] = activity.NumTotal
		ltvHistogram.PaidWeeks[i] = activity.NumPaid
		ltvHistogram.PaidRevenue[i] = activity.TotalAmountPaid
		totalWeeks += activity.NumTotal
		totalPaidWeeks += activity.NumPaid
		totalPaidRevenue += activity.TotalAmountPaid
	}
	ltvHistogram.AverageWeeks = float32(totalWeeks) / float32(numSubs)
	ltvHistogram.AveragePaidWeeks = float32(totalPaidWeeks) / float32(numSubs)
	ltvHistogram.AveragePaidRevenue = totalPaidRevenue / float32(numSubs)

	return ltvHistogram
}

func getBagTypeBreakDown(activities []*subold.SublogSummary) (*BagTypeBreakDown, *BagTypeBreakDown) {
	bagType := &BagTypeBreakDown{}
	bagTypeNewOnly := &BagTypeBreakDown{}
	for _, activity := range activities {
		vegServings := float64(activity.TotalVegServings) / float64(activity.NumTotal)
		nonVegServings := float64(activity.TotalNonVegServings) / float64(activity.NumTotal)
		stillSubed := time.Since(activity.MaxDate) < 7*24*time.Hour
		if vegServings >= 3 {
			bagType.Veg4++
			if stillSubed {
				bagTypeNewOnly.Veg4++
			}
			continue
		}
		if vegServings > .05 {
			bagType.Veg2++
			if stillSubed {
				bagTypeNewOnly.Veg2++
			}
			continue
		}
		if nonVegServings >= 3 {
			bagType.NonVeg4++
			if stillSubed {
				bagTypeNewOnly.NonVeg4++
			}
			continue
		}
		if nonVegServings < 3 {
			bagType.NonVeg2++
			if stillSubed {
				bagTypeNewOnly.NonVeg2++
			}
			continue
		}
	}

	return bagType, bagTypeNewOnly
}

func getBagPriceBreadDown(activities []*subold.SublogSummary) []*BagPriceBreakDown {
	bd := []*BagPriceBreakDown{
		&BagPriceBreakDown{
			Price: 30,
		},
		&BagPriceBreakDown{
			Price: 35.12,
		},
		&BagPriceBreakDown{
			Price: 36.22,
		},
		&BagPriceBreakDown{
			Price: 65.85,
		},
		&BagPriceBreakDown{
			Price: 66.95,
		},
	}
	sort.SliceStable(activities, func(i, j int) bool {
		price1 := activities[i].TotalAmount / float32(activities[i].NumTotal)
		price2 := activities[j].TotalAmount / float32(activities[j].NumTotal)
		return price1 < price2
	})
	for _, activity := range activities {
		price := activity.TotalAmount / float32(activity.NumTotal)
		stillSubed := time.Since(activity.MaxDate) < 7*24*time.Hour
		found := false
		for _, b := range bd {
			if math.Abs(float64(price-b.Price)) < 1 || b.Price == activity.Amount {
				// increase count
				b.Count++
				if stillSubed {
					b.ActiveCount++
				}
				found = true
				break
			}
		}
		// if not found, add type to array
		if !found {
			b := &BagPriceBreakDown{
				Price: activity.Amount,
				Count: 1,
			}
			if stillSubed {
				b.ActiveCount++
			}
			bd = append(bd, b)
		}
	}
	sort.SliceStable(bd, func(i, j int) bool {
		return bd[i].Price < bd[j].Price
	})
	return bd
}
