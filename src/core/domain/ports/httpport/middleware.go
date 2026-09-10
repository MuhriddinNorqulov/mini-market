package httpport

type Middleware func(next HandlerFunc) HandlerFunc
