package executionstats

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
)

// Kind is the kind for datastore.
const Kind = "ExecutionStats"

// ExecutionStats is the status related to a culture execution.
type ExecutionStats struct {
	ID              int64           `json:"id,omitempty"`
	CreatedDatetime time.Time       `json:"created_datetime,omitempty"`
	Date            string          `json:"date,omitempty"`
	Location        common.Location `json:"location,omitempty"`
	Country         string          `json:"country,omitempty"`
	City            string          `json:"city,omitempty"`
	Revenue         float32         `json:"revenue,omitempty"`
	Payroll         []Payroll       `json:"payroll,omitempty"`
	PayrollCosts    float32         `json:"payroll_costs,omitempty"`
	FoodCosts       float32         `json:"food_costs,omitempty"`
	DeliveryCosts   float32         `json:"delivery_costs,omitempty"`
	OnboardingCosts float32         `json:"onboarding_costs,omitempty"`
	ProcessingCosts float32         `json:"processing_costs,omitempty"`
	TaxCosts        float32         `json:"tax_costs,omitempty"`
	PackagingCosts  float32         `json:"packaging_costs,omitempty"`
	OtherCosts      float32         `json:"other_costs,omitempty"`
}

// Payroll is the payroll.
type Payroll struct {
	Name    string  `json:"name,omitempty"`
	Hours   float32 `json:"hours,omitempty"`
	Wage    float32 `json:"wage,omitempty"`
	Postion string  `json:"postion,omitempty"` // Headcook, Linecook, Prepcook, etc
}
