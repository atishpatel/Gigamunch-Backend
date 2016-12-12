package types

const (
	// EMPickup is an ExchangeMethod that is pickup.
	EMPickup ExchangeMethod = ExchangeMethod(1)
	// EMCookDelivery is an ExchangeMethod that is cook delivery.
	EMCookDelivery ExchangeMethod = ExchangeMethod(2)
	// EMGigamunchDelivery is an ExchangeMethod that is Gigamunch delivery.
	EMGigamunchDelivery ExchangeMethod = ExchangeMethod(4)
)

var (
	gigamunchPoint = GeoPoint{Latitude: 36.1513632, Longitude: -86.7255927}
)

// ExchangeMethod is the ExchangeMethod.
type ExchangeMethod int32

// Pickup returns Pickup
func (em ExchangeMethod) Pickup() bool {
	return em.Equal(EMPickup)
}

func (em ExchangeMethod) Delivery() bool {
	return !em.Equal(EMPickup)
}

// CookDelivery returns CookDelivery
func (em ExchangeMethod) CookDelivery() bool {
	return em.Equal(EMCookDelivery)
}

// GigamunchDelivery returns GigamunchDelivery
func (em ExchangeMethod) GigamunchDelivery() bool {
	return em.Equal(EMGigamunchDelivery)
}

// Valid checks if a ExchangeMethod is valid.
func (em ExchangeMethod) Valid() bool {
	return em.Pickup() || em.CookDelivery() || em.GigamunchDelivery()
}

// Equal compares two ExchangeMethod
func (em ExchangeMethod) Equal(em2 ExchangeMethod) bool {
	return int32(em) == int32(em2)
}

func (em ExchangeMethod) ID() int64 {
	return int64(em)
}

func (em ExchangeMethod) String() string {
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

type ExchangeMethodWithPrice struct {
	ExchangeMethod
	Price float32
}

// GetExchangeMethods gets the available exchange methods.
func GetExchangeMethods(cookPoint GeoPoint, cookDeliveryRange int32, cookDeliveryPrice float32, eaterPoint GeoPoint) []ExchangeMethodWithPrice {
	ems := []ExchangeMethodWithPrice{ExchangeMethodWithPrice{ExchangeMethod: EMPickup, Price: 0.0}}
	if cookPoint.GreatCircleDistance(eaterPoint) < float32(cookDeliveryRange) {
		ems = append(ems, ExchangeMethodWithPrice{ExchangeMethod: EMCookDelivery, Price: cookDeliveryPrice})
	}
	if InGigadeliveryRange(cookPoint, eaterPoint) {
		ems = append(ems, ExchangeMethodWithPrice{ExchangeMethod: EMGigamunchDelivery, Price: 5.0})
	}
	return ems
}

// InGigadeliveryRange return if the cook and eater are in Gigamunch Delivery Range.
func InGigadeliveryRange(cookPoint GeoPoint, eaterPoint GeoPoint) bool {
	return gigamunchPoint.GreatCircleDistance(cookPoint) < 60 && gigamunchPoint.GreatCircleDistance(eaterPoint) < 30
}
