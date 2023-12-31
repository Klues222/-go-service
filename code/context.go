package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

//var (
//	jsonUseNumber = true
//	jsonDisallowUnknown = true
//)

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter
	//Ctx context.Context
	RespData       []byte
	RespStatusCode int
	PathParams     map[string]string
	QueryValues    url.Values
	MatchedRoute   string
	//cookieSameSite http.SameSite

	tplEngine TemplateEngine
	//存session
	UserValue map[string]any
}

//func (c *Context) ErrPage() {
//
//}

func (c *Context) Render(tplName string, data any) error {
	var err error
	c.RespData, err = c.tplEngine.Render(c.Req.Context(), tplName, data)
	if err != nil {
		c.RespStatusCode = http.StatusInternalServerError
		return err

	}
	c.RespStatusCode = http.StatusOK
	return nil

}

func (c *Context) SetCookie(ck *http.Cookie) {
	//不推荐
	//ck.SameSite = c.cookieSameSite
	http.SetCookie(c.Resp, ck)
}

func (c *Context) RespJSONOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

func (c *Context) RespJSON(status int, val any) error {

	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.WriteHeader(status)
	//c.Resp.Header().Set("Content-Type", "application/json")
	//c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
	c.RespData = data
	c.RespStatusCode = status
	return err
}

// 获取json，body
func (c *Context) BindJSON(val any) error {
	if val == nil {
		return errors.New("web: 输入不能为nil")
	}
	if c.Req.Body == nil {
		return errors.New("web:body为nil")
	}
	//json.unmarshal:将json数据解析后存储到指向目标结构体的指针中
	//marshal将结构体转成json
	decoder := json.NewDecoder(c.Req.Body) //创建json解码器
	//usenumber数字用number表示否则默认float64
	decoder.UseNumber()
	//如果有一个未知字段，就会报错 User只有两个字段 json有三个就报错
	decoder.DisallowUnknownFields()
	return decoder.Decode(val) //解码成json数据流
}

// 处理表单数据
func (c *Context) FormValue(key string) (string, error) {
	err := c.Req.ParseForm()
	//解析表单出错返回表单错误原因
	if err != nil {
		return "", err
	}
	//返回表单内容
	return c.Req.FormValue(key), nil
	//vals, ok := c.Req.Form[key]
	//if !ok {
	//	return "", errors.New("web: key不存在")
	//}
	//return vals[0], nil
	//

}

// 查询参数
func (c *Context) QueryValue(key string) (string, error) {
	if c.QueryValues == nil {

		//Query获取url上的数据
		c.QueryValues = c.Req.URL.Query()

	}
	vals, ok := c.QueryValues[key]
	if !ok {
		return "", errors.New("web : 找不到这个key")

	}

	return vals[0], nil
}
func (c *Context) PathValueV1(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{
			err: errors.New("web: key不存在"),
		}

	}

	return StringValue{
		val: val,
	}

}

// 路由参数
func (c *Context) PathValue(key string) (string, error) {
	val, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("web : 找不到这个key")

	}

	return val, nil

}

type StringValue struct {
	val string
	err error
}

func (s StringValue) AsIn64() (int64, error) {
	if s.err != nil {
		return 0, s.err

	}
	return strconv.ParseInt(s.val, 10, 64)
}
