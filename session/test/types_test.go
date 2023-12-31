package test

import (
	"net/http"
	_ "os"
	"testing"
	"time"
	web "webdemo/code"
	"webdemo/session"
	"webdemo/session/cookie"
	"webdemo/session/memory"
)

func TestSession(t *testing.T) {
	//登陆校验

	var m *session.Manager = &session.Manager{
		Propagator: cookie.NewPropagator(),
		Store:      memory.NewStore(time.Minute * 15),
		CtxSessKey: "sessKey",
	}
	server := web.NewHTTPServer(web.ServerWithMiddleware(func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			if ctx.Req.URL.Path == "/login" {
				//放过去用户准备登陆
				next(ctx)
				return
			}
			//从网页拿session id
			_, err := m.GetSession(ctx)
			if err != nil {
				ctx.RespStatusCode = http.StatusUnauthorized
				ctx.RespData = []byte("请重新登陆")
				return
			}
			//刷新session的过期使劲按
			m.RefreshSession(ctx)
			//next
			next(ctx)
		}
	}))

	server.Post("/login", func(ctx *web.Context) {
		//要在这之前校验登陆密码
		sess, err := m.InitSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登陆失败")
		}
		err = sess.Set(ctx.Req.Context(), "nikename", "xiaoming")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登陆失败")
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("登陆成功")

		return
	})

	//退出
	server.Post("/login", func(ctx *web.Context) {
		//要在这之前校验登陆密码
		err := m.RemoveSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("退出失败")
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("退出成功")
	})
	//把拿到的session数据展示出来
	server.Get("/user", func(ctx *web.Context) {
		//根据session id找session
		sess, _ := m.GetSession(ctx)

		val, _ := sess.Get(ctx.Req.Context(), "nickname")
		ctx.RespData = []byte(val.(string))
	})

	server.Start(":8081")

}
