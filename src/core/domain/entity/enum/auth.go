package enum

type LoginMethod string

const (
	LoginMethodPassword LoginMethod = "PASSWORD"
	LoginMethodOtp      LoginMethod = "OTP"
	LoginMethodGoogle   LoginMethod = "GOOGLE"
)

type LoginResult string

const (
	LoginResultSuccess LoginResult = "SUCCESS"
	LoginResultFailure LoginResult = "FAILURE"
)

type LoginFailureReason string

const (
	LoginFailureWrongPassword LoginFailureReason = "WRONG_PASSWORD"
	LoginFailureWrongOtp      LoginFailureReason = "WRONG_OTP"
	LoginFailureExpiredFlow   LoginFailureReason = "EXPIRED_FLOW"
	LoginFailureRateLimited   LoginFailureReason = "RATE_LIMITED"
)

type AuthFlowScope string

const (
	AuthFlowSignin        AuthFlowScope = "SIGNIN"
	AuthFlowSignup        AuthFlowScope = "SIGNUP"
	AuthFlowSetPassword   AuthFlowScope = "SET_PASSWORD"
	AuthFlowResetPassword AuthFlowScope = "RESET_PASSWORD"
)
