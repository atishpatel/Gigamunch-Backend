package common

// PaymentProvider is a payment provider.
type PaymentProvider int8

const (
	// Braintree is the braintree provider.
	Braintree PaymentProvider = 0
	// Stripe is the Stripe provider.
	Stripe PaymentProvider = 1
)
