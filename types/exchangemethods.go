package types

const (
	// PickupOnlyExchangeMethod is an ExchangeMethods that is pickup only
	PickupOnlyExchangeMethod ExchangeMethods = ExchangeMethods(1)
	// ChefDeliveryOnlyExchangeMethod is an ExchangeMethods that is chef delivery only
	ChefDeliveryOnlyExchangeMethod ExchangeMethods = ExchangeMethods(2)
)

// ExchangeMethods shows ExchangeMethod options
type ExchangeMethods int32

// IsZero returns true if there are no ExchangeMethods selected
func (em *ExchangeMethods) IsZero() bool {
	return int32(*em) == 0
}

// Pickup returns Pickup
func (em *ExchangeMethods) Pickup() bool {
	return getKthBit(int32(*em), 0)
}

// SetPickup sets Pickup
func (em *ExchangeMethods) SetPickup(x bool) {
	*em = ExchangeMethods(setKthBit(int32(*em), 0, x))
}

// Delivery returns if any type of Delivery is true
func (em *ExchangeMethods) Delivery() bool {
	return (int32(*em) % 2) == 0
}

// ChefDelivery returns ChefDelivery
func (em *ExchangeMethods) ChefDelivery() bool {
	return getKthBit(int32(*em), 1)
}

// SetChefDelivery sets ChefDelivery
func (em *ExchangeMethods) SetChefDelivery(x bool) {
	*em = ExchangeMethods(setKthBit(int32(*em), 1, x))
}

// Equal compares two exchangemethods
func (em ExchangeMethods) Equal(em2 ExchangeMethods) bool {
	return int32(em) == int32(em2)
}

func (em ExchangeMethods) String() string {
	v := int32(em)
	switch v {
	case 1:
		return "Pickup"
	case 2:
		return "Cook Delivery"
	case 4:
		return "Gigamunch Delivery"
	}
	return ""
}
