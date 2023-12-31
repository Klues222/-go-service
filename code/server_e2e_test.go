//go:build e2e

package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	h := &HTTPServer{}

	h.AddRoute(http.MethodGet, "/uer", func(ctx *Context) {
		fmt.Println("first")
		fmt.Println("second")
	})
	handle1 := func(ctx *Context) {
		fmt.Println("first")
	}
	handle2 := func(ctx *Context) {
		fmt.Println("second")
	}
	h.AddRoute(http.MethodGet, "/user", func(ctx *Context) {
		handle1(ctx)
		handle2(ctx)
	})

	//h.AddRoute1(http.MethodGet, "/", func(ctx Context) {
	//	fmt.Println("first")
	//}, func(ctx Context) {
	//	fmt.Println("first")
	//})

	//http.ListenAndServe(":8081", h)
	//http.ListenAndServeTLS(":8088", "", "", h)
	h.Start(":8081")
}
