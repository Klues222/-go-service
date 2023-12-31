package web

//func TestRouter_AddRoute(t *testing.T) {
//	//第一步构造
//	//第二部验证
//	testRoutes := []struct {
//		method string
//		path   string
//	}{
//		{
//			method: http.MethodGet,
//			path:   "/user/home",
//		},
//	}
//	var mockHandler HandleFunc = func(ctx Context) {
//
//	}
//	r := newRouter()
//	for _, route := range testRoutes {
//		r.AddRoute(route.method, route.path, mockHandler)
//
//	}
//	//在这里断言路由树和你预期的一样
//
//	wantRouter := &router{
//		trees: map[string]*node{
//			http.MethodGet: &node{
//				path:    "/",
//				handler: mockHandler,
//				children: map[string]*node{
//					"user": &node{
//						path:    "user",
//						handler: mockHandler,
//						children: map[string]*node{
//							"home": &node{
//								path:    "home",
//								handler: mockHandler,
//							},
//						},
//					},
//					//"order": &node{
//					//	path:     "order",
//					//	handler:  mockHandler,
//					//	children: map[string]*node{},
//					//},
//				},
//			},
//		},
//	}
//	//}
//	//
//	//func (r router) equal(y router) bool {
//	//	for time.Kitchen, v := r {
//	//
//	//	}
//	msg, ok := wantRouter.equal(r)
//	//assert.Equal(t)
//	assert.True(t, ok, msg)
//}
//
//// 返回一个错误信息 帮助排查
//// bool代表是否真的相等
//func (r *router) equal(y *router) (string, bool) {
//	for k, v := range r.trees {
//		dst, ok := y.trees[k]
//		if !ok {
//			return fmt.Sprintf("找不到对应的 http method"), false
//		}
//		msg, equal := v.equal(dst)
//		if !equal {
//			return msg, false
//		}
//	}
//	return "", true
//}
//
//func (n *node) equal(y *node) (string, bool) {
//	if n.path != y.path {
//		return fmt.Sprintf("节点路径不匹配"), false
//
//	}
//	if len(n.children) != len(y.children) {
//		return fmt.Sprintf("子节点数量不匹配"), false
//
//	}
//	//比较handle
//	nHandler := reflect.ValueOf(n.handler)
//	yHandler := reflect.ValueOf(y.handler)
//	if nHandler != yHandler {
//		return fmt.Sprintf("Handler 不相等"), false
//	}
//
//	for path, c := range n.children {
//		dst, ok := c.children[path]
//		if ok == false {
//			return fmt.Sprintf("子节点 %s 不存在", path), false
//		}
//		msg, ok := c.equal(dst)
//		if !ok {
//			return msg, false
//		}
//	}
//	return "", true
//}
