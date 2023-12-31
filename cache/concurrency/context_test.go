package context

import (
	"context"
	"testing"
	"time"
)

type mykey struct {
}

func TestContext(t *testing.T) {
	//一般是链路起点， 或者调用的起点
	ctx := context.Background()
	// 在你不确定 context 该用啥的时候，用 TODO()
	ctx = context.WithValue(ctx, mykey{}, "my-value")
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	val := ctx.Value(mykey{}).(string)
	t.Log(val)
	newval := ctx.Value("不存在key")
	val, ok := newval.(string)
	if !ok {
		t.Log("类型不对")
		return
	}
	t.Log(val)
}

func TestContext_WithCancel(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	//用完ctx再去调用
	//defer cancel()
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()
	//用ctx
	//<- 将箭头后面的东西发送到通道

	<-ctx.Done()
	t.Log("hello, cancel: ", ctx.Err())

}

func TestContext_WithDeadline(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	deadline, _ := ctx.Deadline()
	t.Log("deadline: ", deadline)
	defer cancel()
	<-ctx.Done()
	t.Log("hello, deadline: ", ctx.Err())
}

func TestContext_WithTimeout(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	deadline, _ := ctx.Deadline()
	t.Log("deadline: ", deadline)
	defer cancel()
	<-ctx.Done()
	t.Log("hello, deadline: ", ctx.Err())
}
