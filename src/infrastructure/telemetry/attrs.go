package telemetry

const (
	AttrErrorCode = "app.error_code"
	AttrCaller    = "app.caller"
	AttrUserID    = "app.user_id"
	AttrUserRole  = "app.user_role"
	AttrTaskType  = "app.task_type"
	AttrTaskID    = "app.task_id"
	AttrWsChannel = "app.ws_channel"
	AttrWsConnID  = "app.ws_conn_id"
	AttrWsEvent   = "app.ws_event"

	AttrErrorKind = "app.error_kind"
)

const (
	ErrorKindPanic = "panic"

	ErrorKindError = "error"

	ErrorKindBusiness = "business"
)
