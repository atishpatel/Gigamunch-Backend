package main

import (
	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"golang.org/x/net/context"
)

// CohortCell is a cell.
type CohortCell struct {
	Interval         int16   `json:"interval,omitempty"`
	AmountLeft       int     `json:"amount_left,omitempty"`
	AmountLost       int     `json:"amount_lost,omitempty"`
	RetentionPercent float32 `json:"retention_percent,omitempty"`
}

// CohortRow is a row.
type CohortRow struct {
	StartDate  string       `json:"start_date,omitempty"`
	CohortCell []CohortCell `json:"cohort_cell,omitempty"`
}

// CohortAnalysis is a full cohortAnalysis.
type CohortAnalysis struct {
	CohortRows []CohortRow `json:"cohort_rows,omitempty"`
}

// GetGeneralStatsResp is a response for GetGeneralStats.
type GetGeneralStatsResp struct {
	Activities           []*sub.SublogSummary `json:"activities"`
	WeeklyCohortAnalysis *CohortAnalysis      `json:"weekly_cohort_analysis"`
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
	resp.Activities = activities
	resp.WeeklyCohortAnalysis = getWeeklyCohort(activities)
	return resp, nil
}

func getWeeklyCohort(activities []*sub.SublogSummary) *CohortAnalysis {
	analysis := new(CohortAnalysis)

	return analysis
}
