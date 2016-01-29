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
	// SessionCookieName is the name of the cookie for session id
	SessionCookieName = "GIGASID"
)
