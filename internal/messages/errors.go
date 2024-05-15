package messages

const (
	// UnauthorizedError is the error message for unauthorized requests
	UnauthorizedError = "unauthorized"
	// InvalidGraphqlOperationMessage is the error message for invalid GraphQL requests
	InvalidGraphqlOperationMessage = "invalid GraphQL operation"
	// GinContextError is the error message for gin context errors
	GinContextError = "error getting gin context"
	// UserNotFoundError is the error message for user not found
	UserNotFoundError = "user not found"
	// InvalidEmailAddressError is the error message for invalid email address
	InvalidEmailAddressError = "invalid email address"
	// UserAlreadyExistsError is the error message for user already exists
	UserAlreadyExistsError = "user already exists"
	// InvalidMfaSessionError is the error message for invalid MFA session
	InvalidMfaSessionError = "invalid MFA session"
	// InvalidOtpError is the error message for invalid OTP
	InvalidOtpError = "invalid OTP"
	// ErrorGeneratingSession is the error message for session generation error
	ErrorGeneratingSession = "error generating session"
	// ErrorVerifyingSession is the error message for session verification error
	ErrorVerifyingSession = "error verifying session"
	// InternalServerError is the error message for internal server error
	InternalServerError = "internal server error"
	// InvalidSessionError is the error message for invalid session
	InvalidSessionError = "invalid session"
	// InvalidDeploymentRequestError is the error message for invalid deployment request
	InvalidDeploymentRequestError = "invalid deployment request"
	// DeploymentNotFoundError is the error message for deployment not found
	DeploymentNotFoundError = "deployment not found"
)
