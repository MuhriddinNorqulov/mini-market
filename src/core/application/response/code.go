package response

type Code string

const (
	CodeSuccess Code = "OK"

	CodNotFound         Code = "NOT_FOUND"
	CodeConflict        Code = "CONFLICT"
	CodeBadRequest      Code = "BAD_REQUEST"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeDatabaseError   Code = "DATABASE_ERROR"
	CodeGatewayError    Code = "GATEWAY_ERROR"
	CodeInvalidToken    Code = "INVALID_TOKEN"
	CodeExpiredToken    Code = "EXPIRED_TOKEN"
	CodeInfluentBalance Code = "INFLUENT_BALANCE"
	CodeFileError       Code = "FILE_ERROR"
	CodeForbidden       Code = "FORBIDDEN"
	CodeQuotaExceeded   Code = "QUOTA_EXCEEDED"
	CodeAsyncError      Code = "ASYNC_ERROR"
	CodeNotVerified     Code = "NOT_VERIFIED"
	CodeTooManyRequests Code = "TOO_MANY_REQUESTS"

	CodeUserAlreadyAssigned Code = "USER_ALREADY_ASSIGNED"
	CodeDepartmentNotFound  Code = "DEPARTMENT_NOT_FOUND"
	CodeNotDepartmentMember Code = "NOT_DEPARTMENT_MEMBER"

	CodeInternalServerError Code = "INTERNAL_SERVER_ERROR"
)
