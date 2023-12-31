package web

import "strings"

// 用来支持对路由书的操作
// 代表路由树森林
type router struct {
	//到路由树更节电
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}

}

//加一些限制
//path 必须以 / 开头 ，不能以 / 结尾， 中间也不能有连续的//

func (r *router) AddRoute(method string, path string, handleFunc HandleFunc) {
	//首先找到树来
	if path == "" {
		panic("web不能为空字符串")
	}
	root, ok := r.trees[method]
	if !ok {
		//说明没有根节点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	if path[0] != '/' {
		panic("web必须以/开头")
	}
	if path == "/" {
		if root.handler != nil {
			panic("web:重复注册")
		}
		root.handler = handleFunc
		root.route = "/"
		return
	}
	path = path[1:]
	//切割这个 path
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		//递归下去，找准位置
		//如果中途有不在，就创建
		children := root.childOrCreate(seg)
		root = children
	}
	i := segs
	print(i)
	root.handler = handleFunc
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	//基本上是沿着树深度查找下去
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}
	path = strings.Trim(path, "/")
	seqs := strings.Split(path, "/")
	var pathParams map[string]string
	for _, seq := range seqs {
		child, paramChild, found := root.childOf(seq)
		if !found {
			return nil, false
		}
		//命中了路径参数
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string)

			}
			pathParams[child.path[1:]] = seq
		}
		root = child

	}
	//daibiaoyoujiedian
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true
}

func (n *node) childOrCreate(seg string) *node {

	if seg[0] == ':' {
		if n.starChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配,已有通配符匹配")
		}

		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}

	if seg == "*" {
		if n.paramChild != nil {
		}
		panic("web: 不允许同时注册路径参数和通配符匹配，已有路径参数")
		n.starChild = &node{

			path: seg,
		}
		return n.starChild
	}

	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		//新建
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

//type tree struct {
//	root *node
//}

type node struct {
	route string
	path  string
	//子 path 到子节点的映射
	children map[string]*node
	//缺一个代表用户注册的业务逻辑
	handler    HandleFunc
	starChild  *node
	paramChild *node
}

// childOf 优先考虑静态匹配， 匹配不上，再考虑通配符匹配
// 第一个返回值是子节点
// 第二个是标记是否是路径参数
// 第三个标记命中了没有
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}

		return n.starChild, false, n.starChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return child, false, ok
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}
