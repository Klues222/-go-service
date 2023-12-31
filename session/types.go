package session

import (
	"context"
	"net/http"
)

//var (ErrKeyNotFound = errors.New(""))

// 管理session本身
type Store interface {
	//session ID 谁来指定
	Generate(ctx context.Context, id string) (Session, error)
	//更新session
	Refresh(ctx context.Context, id string) error
	//销毁
	Remove(ctx context.Context, id string) error
	//
	Get(ctx context.Context, id string) (Session, error)
}

// session本体
type Session interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any) error
	ID() string
}

// 联系http
type Propagator interface {
	//session id 注入到HTTP响应里面
	Inject(id string, writer http.ResponseWriter) error
	//将 session id 从HTTP 请求中提取出来
	Extract(req *http.Request) (string, error)
	//销毁session id
	Remove(writer http.ResponseWriter) error
}
