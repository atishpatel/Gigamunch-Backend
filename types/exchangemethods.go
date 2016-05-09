package types

// ExchangeMethods shows ExchangeMethod options
type ExchangeMethods int32

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
