package constants

const (
	NO_WEBSITES_ADDED         = "No websites added, use /POST to add websites to status check"
	UNEXPECTED_ENDPOINT       = "POST/GET Method expected on this endpoint"
	BAD_REQUEST_UNKNOWN_FIELD = "Bad Request. Incorrect format for field : "
	BAD_REQUEST               = "Bad Request "
	POLLING_RATE              = 10 /* Time after which status should be checked again */
)
