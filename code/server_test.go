package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHTTPServer_ServeHTTP(t *testing.T) {
	server := NewHTTPServer()
	server.mdls = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println(" 第一个before")
				next(ctx)
				fmt.Println("第1个after")

			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个before")
				next(ctx)
				fmt.Println("第二个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第3个中断")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("你看不到")
			}
		},
	}
	server.ServeHTTP(nil, &http.Request{})
}
