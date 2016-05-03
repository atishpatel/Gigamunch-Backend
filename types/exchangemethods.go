package types

// ExchangeMethods shows ExchangeMethod options
type ExchangeMethods int32

// Pickup returns Pickup
func (em ExchangeMethods) Pickup() bool {
	return getKthBit(int32(em), 0)
}

// ChefDelivery returns ChefDelivery
func (em ExchangeMethods) ChefDelivery() bool {
	return getKthBit(int32(em), 1)
}

// SetPickup sets Pickup
func (em ExchangeMethods) SetPickup(x bool) {
	setKthBit(int32(em), 0, x)
}

// SetChefDelivery sets ChefDelivery
func (em ExchangeMethods) SetChefDelivery(x bool) {
	setKthBit(int32(em), 1, x)
}
