package server

const (
	baseChefURL  = "/chef"
	baseAdminURL = "/admin"
	/*****************************************************************************
	*	URLs
	*****************************************************************************/
	homeURL = "/"

	baseLoginURL = "/login"
	// loginURL is the url for chefs to login
	loginURL = baseLoginURL + "?mode=select"
	// sendEmailURL sends emails for forgot password for gitkit
	sendEmailURL = "/sendemail"
	// signOutURL is the url for signing out users
	signOutURL = "/signout"

	errorURL = "/error"

	// chef stuff
	chefHomeURL = baseChefURL

	chefApplicationURL = "/chef/apply"

	// admin stuff
	adminHomeURL = baseAdminURL
	/*****************************************************************************
	*	Cookies
	*****************************************************************************/
	// gitkitCookieName is the name of the cookie for gitkit
	gitkitCookieName = "gtoken"
	// sessionTokenCookieName is the name of the cookie for the session token
	sessionTokenCookieName = "GIGATKN"
)
