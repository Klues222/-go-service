package web

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngine interface {
	//Render渲染页面
	//data渲染页面所需的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
	//不需要
	//AddTemplate()
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	bs := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(bs, tplName, data)
	return bs.Bytes(), err
}
