package async

type AsyncTaskMiddleware func(handler AsyncTaskHandler) AsyncTaskHandler
