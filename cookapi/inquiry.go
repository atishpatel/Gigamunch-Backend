package main

import (
	"context"

	"github.com/atishpatel/Gigamunch-Backend/corenew/inquiry"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// GetInquiriesReq is the request for GetInquries.
type GetInquiriesReq struct {
	GigatokenReq
	StartIndex int `json:"start_index"`
	EndIndex   int `json:"end_index"`
}

// InquiryResp is a response with an Inquiry and err.
type InquiryResp struct {
	Inquiry Inquiry              `json:"inquiry"`
	Err     errors.ErrorWithCode `json:"err"`
}

// InquiriesResp is a response with multiple Inquiry and err.
type InquiriesResp struct {
	Inquiries []Inquiry            `json:"inquiries"`
	Err       errors.ErrorWithCode `json:"err"`
}

// GetInquiries gets the cook's Inquiries.
func (service *Service) GetInquiries(ctx context.Context, req *GetInquiriesReq) (*InquiriesResp, error) {
	resp := new(InquiriesResp)
	defer handleResp(ctx, "GetInquiries", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	inquiryC := inquiry.New(ctx)
	inquiries, err := inquiryC.GetByCookID(user.ID, req.StartIndex, req.EndIndex)
	if err != nil {
		resp.Err = errors.Wrap("failed to inquiry.GetByCookID", err)
		return resp, nil
	}

	for i := range inquiries {
		resp.Inquiries = append(resp.Inquiries, Inquiry{
			Inquiry: *inquiries[i],
		})
	}
	return resp, nil
}

// GetInquiry gets the cook's Inquiry.
func (service *Service) GetInquiry(ctx context.Context, req *IDReq) (*InquiryResp, error) {
	resp := new(InquiryResp)
	defer handleResp(ctx, "GetInquiry", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	inquiryC := inquiry.New(ctx)
	inq, err := inquiryC.Get(user, req.ID)
	if err != nil {
		resp.Err = errors.Wrap("failed to inquiry.Get", err)
		return resp, nil
	}
	resp.Inquiry = Inquiry{Inquiry: *inq}
	return resp, nil
}

// AcceptInquiry accepts an inquiry.
func (service *Service) AcceptInquiry(ctx context.Context, req *IDReq) (*InquiryResp, error) {
	resp := new(InquiryResp)
	defer handleResp(ctx, "AcceptInquiry", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	inquiryC := inquiry.New(ctx)
	inquiry, err := inquiryC.CookAccept(user, req.ID)
	if err != nil {
		resp.Err = errors.Wrap("failed to inquiry.CookAccept", err)
		return resp, nil
	}
	resp.Inquiry.Inquiry = *inquiry
	return resp, nil
}

// DeclineInquiry declines an inquiry.
func (service *Service) DeclineInquiry(ctx context.Context, req *IDReq) (*InquiryResp, error) {
	resp := new(InquiryResp)
	defer handleResp(ctx, "DeclineInquiry", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	inquiryC := inquiry.New(ctx)
	inquiry, err := inquiryC.CookDecline(user, req.ID)
	if err != nil {
		resp.Err = errors.Wrap("failed to inquiry.CookDecline", err)
		return resp, nil
	}
	resp.Inquiry.Inquiry = *inquiry
	return resp, nil
}
