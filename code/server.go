package web

import (
	"fmt"
	"net"
	"net/http"
)

// 确保一定实现了httpserver
var _ Server = &HTTPServer{}

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	Start(addr string) error

	//增加路由注册功能
	//method 是HTTP 方法
	//path 是路由
	//handlefunc 是你的业务逻辑

	AddRoute(method string, path string, handleFunc HandleFunc)
	//AddRoute1(method string, path string, handeFunc HandleFunc, handFunc1 HandleFunc)
	//这种提供多次
}

// ServeHTTP 处理请求入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//你的框架代码就在这里
	ctx := &Context{
		Req:       request,
		Resp:      writer,
		tplEngine: h.tplEngine,
	}
	//最后一个是这个
	root := h.serve
	//然后这里就是利用最后一个不断回溯组装链条
	//从后往前
	//把后一个作为前一个的next构造好链条
	for i := len(h.mdls) - 1; i >= 0; i-- {
		root = h.mdls[i](root)

	}
	//从这里执行就是从后往前
	//这里最后一个步骤就是将RespData和Respstatuscode刷新到响应里

	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			h.flashResp(ctx)
		}
	}
	root = m(root)

	root(ctx)
}

func (h *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)

	}
	n, err := ctx.Resp.Write(ctx.RespData)
	if err != nil || n != len(ctx.RespData) {
		h.log("写入响应失败 %v", err)
	}
}

func (h *HTTPServer) serve(ctx *Context) {
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.n.handler == nil {
		//路由没有命中， 404
		ctx.RespStatusCode = 404
		ctx.RespData = []byte("NotFound")
		return
	}
	ctx.PathParams = info.pathParams
	ctx.MatchedRoute = info.n.route
	info.n.handler(ctx)
}

//func (h *HTTPServer) AddRoute(method string, path string, handleFunc HandleFunc) {
//	//这里注册到路由树里面
//	//panic("i")
//}

//func (h *HTTPServer) AddRoute1(method string, path string, handeFunc HandleFunc, handFunc1 HandleFunc) {
//
//}

func (h *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	//在这里可以让用户注册所谓的after start
	return http.Serve(l, h)
}

type HTTPServerOption func(server *HTTPServer)

type HTTPServer struct {
	//router
	router
	mdls []Middleware

	log func(msg string, args ...any)

	tplEngine TemplateEngine
}

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{
		router: newRouter(),
		log: func(msg string, args ...any) {
			fmt.Printf(msg, args...)
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func ServerWithTemplateEngine(tplEngine TemplateEngine) HTTPServerOption {
	return func(server *HTTPServer) {
		server.tplEngine = tplEngine
	}
}

func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.mdls = mdls
	}

}
func (N *HTTPServer) Get(path string, handlefunc HandleFunc) {
	N.AddRoute(http.MethodGet, path, handlefunc)
}
func (N *HTTPServer) Post(path string, handlefunc HandleFunc) {
	N.AddRoute(http.MethodPost, path, handlefunc)
}
