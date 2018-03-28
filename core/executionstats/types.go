package executionstats

// ExecutionStats is the status related to a culture execution.
type ExecutionStats struct {
	Payroll         Payroll
	Revenue         float32
	FoodCosts       float32
	DeliveryCosts   float32
	OnboardingCosts float32
	// TODO: add other fields
	OtherCosts float32
}

// Payroll is the payroll.
type Payroll struct {
	Name string
	// TODO: hours as time slots?
	Hours   int
	Postion string // Headcook, Linecook, Prepcook, etc
}
