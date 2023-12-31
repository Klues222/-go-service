package session

import (
	uuid "github.com/satori/go.uuid"
	web "webdemo/code"
)

//用户友好的封装

type Manager struct {
	Propagator
	Store
	CtxSessKey string
}

func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	if ctx.UserValue == nil {
		ctx.UserValue = make(map[string]any, 1)
	}
	val, ok := ctx.UserValue[m.CtxSessKey]
	if ok {
		return val.(Session), nil
	}
	// 尝试缓存住session
	sessId, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	//根据session id找session
	sess, err := m.Get(ctx.Req.Context(), sessId)
	if err != nil {
		return nil, err
	}
	ctx.UserValue[m.CtxSessKey] = sess
	return sess, err
}

func (m *Manager) InitSession(ctx *web.Context) (Session, error) {
	id := uuid.NewV1().String()
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	//注入http响应里面
	err = m.Inject(id, ctx.Resp)
	return sess, err
}

func (m *Manager) RefreshSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	return m.Refresh(ctx.Req.Context(), sess.ID())
}

func (m *Manager) RemoveSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return nil
	}
	m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}
