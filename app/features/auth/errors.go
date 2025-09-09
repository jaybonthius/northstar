package auth

const (
	MsgInvalidCredentials        = "Invalid email or password"
	MsgUserNotFound              = "User not found"
	MsgUserAlreadyExists         = "User already exists"
	MsgUsernameAlreadyExists     = "An account with this username already exists"
	MsgEmailAlreadyExists        = "An account with this email already exists"
	MsgMissingFields             = "All fields are required"
	MsgMissingCredentials        = "Email and password are required"
	MsgPasswordTooShort          = "Password must be at least 6 characters long"
	MsgInvalidFormData           = "Invalid form data"
	MsgInvalidMethod             = "Invalid request method"
	MsgSessionFailed             = "Session management failed"
	MsgLoginFailed               = "Login failed, please try again"
	MsgSignupFailed              = "Account creation failed, please try again"
	MsgLogoutFailed              = "Logout failed, please try again"
	MsgAccountCreatedLoginFailed = "Account created but login failed, please try logging in manually"
)
