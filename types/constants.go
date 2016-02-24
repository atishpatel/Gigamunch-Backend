package types

const (
	/* ---------------------------------------------------------------------------
	 * Datastore Kinds
	 * ---------------------------------------------------------------------------
	 */

	// KindGigachef is used for the basic Gigachef account info
	KindGigachef = "Gigachef"
	// KindGigamuncher is used for the basic Gigamuncher account info
	KindGigamuncher = "Gigamuncher"
	// KindMealTemplate is used for a Gigachef's meal template
	KindMealTemplate = "MealTemplate"
	// KindMealTemplateTag is used for orginizing meal templates by a Gigachef
	KindMealTemplateTag = "MealTemplateTag"
	// KindMeal is a meal that is either live or finished
	KindMeal = "Meal"
	// KindOrder is an order made by a Gigamuncher on a meal
	KindOrder = "Order"
	// KindGigachefApplication is an application for a chef to become verfified
	KindGigachefApplication = "GigachefApplication"
	/* ---------------------------------------------------------------------------
	 * URLs
	 * ---------------------------------------------------------------------------
	 */

	// LoginURL is the url for chefs to login
	LoginURL = "/login"
	// SendEmailURL sends emails for forgot password for gitkit
	SendEmailURL = "/sendemail"
	// SignOutURL is the url for signing out users
	SignOutURL = "/signout"

	/* ---------------------------------------------------------------------------
	 * Cookies
	 * ---------------------------------------------------------------------------
	 */

	// GitkitCookieName is the name of the cookie for gitkit
	GitkitCookieName = "gtoken"
	// SessionTokenCookieName is the name of the cookie for the session token
	SessionTokenCookieName = "GIGATKN"

	/* ---------------------------------------------------------------------------
	 * Other
	 * ---------------------------------------------------------------------------
	 */

	// DefaultMaxDistanceInMiles is the default distance for live meals
	DefaultMaxDistanceInMiles = 300
)
