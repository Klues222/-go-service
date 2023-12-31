package web

// 函数式的责任链模式
type Middleware func(next HandleFunc) HandleFunc

//拦截器
//type MiddlewareV1 interface {
//	Before(ctx *Context)]7
//	After(ctx *Context)
//	Surround(ctx *Context)
//}
