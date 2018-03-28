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

type Payroll struct {
	Name    string
	Hours   string
	Postion string // Headcook, Linecook, Prepcook, etc
}
