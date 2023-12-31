package template_demo

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"text/template"
)

func TestHelloWorld(t *testing.T) {
	type User struct {
		Name string
	}
	//定义模板的名字
	tpl := template.New("hello-world")
	//解析模板
	//点代表普通结构体或者指针，map或者其他切片或数组
	tpl, err := tpl.Parse(`Hello, {{ .Name}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, User{Name: "Tom"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, Tom`, buffer.String())
}
func TestMap(t *testing.T) {
	type User struct {
		Name string
	}
	//定义模板的名字
	tpl := template.New("hello-world")
	//解析模板
	//点代表普通结构体或者指针，map或者其他切片或数组
	tpl, err := tpl.Parse(`Hello, {{ .Name}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, map[string]string{"Name": "Tom"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, Tom`, buffer.String())
}

func TestSlice(t *testing.T) {
	type User struct {
		Name string
	}
	//定义模板的名字
	tpl := template.New("hello-world")
	//解析模板
	//点代表普通结构体或者指针，map或者其他切片或数组
	tpl, err := tpl.Parse(`Hello, {{index . 0}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, []string{"Tom"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, Tom`, buffer.String())
}
func TestBasic(t *testing.T) {
	//定义模板的名字
	tpl := template.New("hello-world")
	//解析模板
	//点代表普通结构体或者指针，map或者其他切片或数组
	tpl, err := tpl.Parse(`Hello, {{.}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, 123)
	require.NoError(t, err)
	assert.Equal(t, `Hello, 123`, buffer.String())
}

// 传方法
func TestFuncCall(t *testing.T) {
	tpl := template.New("hello-world")
	//解析模板
	//点代表普通结构体或者指针，map或者其他切片或数组
	tpl, err := tpl.Parse(`
切片长度: {{len .Slice}}
{{printf "%.2f" 1.2345}}
Hello, {{.Hello "Tom" "Jerry"}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, FuncCall{
		Slice: []string{"a", "b"},
	})
	require.NoError(t, err)
	assert.Equal(t, `
切片长度: 2
1.23
Hello, Tom Jerry`, buffer.String())

}

func TestForloop(t *testing.T) {
	tpl := template.New("hello-world")
	//解析模板
	//点代表普通结构体或者指针，map或者其他切片或数组
	tpl, err := tpl.Parse(`
{{- range $idx, $ele := .Slice}}
{{- .}}
{{ end}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, FuncCall{
		Slice: []string{"a", "b"},
	})
	require.NoError(t, err)
	assert.Equal(t, `a
b

`, buffer.String())

}

func TestIfElse(t *testing.T) {
	type user struct {
		Age int
	}
	tpl := template.New("hello-world")
	//解析模板
	//点代表普通结构体或者指针，map或者其他切片或数组
	tpl, err := tpl.Parse(`
{{- if and (gt .Age 0) (le .Age 6)}}
我是儿童: (0, 6]
{{- else if and (gt .Age 6) (le .Age 18)}}
我是少年: (6, 18]
{{- else}}
我是成人: >18
{{end -}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, user{Age: 100})
	require.NoError(t, err)
	assert.Equal(t, ``, buffer.String())

}

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(first string, last string) string {
	return fmt.Sprintf("%s %s", first, last)
}
