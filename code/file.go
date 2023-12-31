package web

import (
	lru "github.com/hashicorp/golang-lru"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	_ "unicode"
)

type FileUploader struct {
	FileField string
	//为什么要用户传？
	//要考虑文件名冲突的问题
	//所以很多时候，目标文件名是随机的
	DstPathFunc func(*multipart.FileHeader) string
}

func (f FileUploader) Handle() HandleFunc {
	return func(ctx *Context) {
		//上传逻辑在这

		//第一步：读到文件内容
		//第二步计算出目标路径
		//第三步保存文件
		//第四步返回响应
		file, fileHeader, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		//我怎么知道目标路径
		//这种做法将计算目标路径的逻辑交给路径
		dst := f.DstPathFunc(fileHeader)
		//将dst上不存在的目录全部建立起来
		//os.MkdirAll()
		//O_WRONLY 写入数据
		//O_TRUNC 如果文件本身存在，清空数据
		//O_CREATE创建一个新的
		dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
		//buf 会影响你的性能
		//你要考虑复用
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		defer dstFile.Close()
		_, err = io.CopyBuffer(dstFile, file, nil)
		if err != err {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("上传成功")

	}
}

type FileDownloader struct {
	Dir string
}

func (d FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		//用的是xxx?file=xxx
		req, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("找不到目标文件")
			return
		}
		dst := filepath.Join(d.Dir, req)
		fn := filepath.Base(dst)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "must-revalidate")
		http.ServeFile(ctx.Resp, ctx.Req, dst)
	}
}

type StaticReSourceHandlerOption func(handler *StaticResourceHandler)

type StaticResourceHandler struct {
	dir                     string
	cache                   *lru.Cache
	extensionContextTypeMap map[string]string
	maxSize                 int
}

func NewStaticResourceHandler(dir string, opts ...StaticReSourceHandlerOption) (*StaticResourceHandler, error) {
	c, err := lru.New(1000)
	if err != nil {
		return nil, err
	}

	res := &StaticResourceHandler{
		dir:     dir,
		cache:   c,
		maxSize: 1024 * 1024 * 10,
		extensionContextTypeMap: map[string]string{
			"jpeg": "image/jpge",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"dpf":  "image/pdf",
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func StaticWithMaxFileSize(maxSize int) StaticReSourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.maxSize = maxSize
	}
}

func StaticWithCache(c *lru.Cache) StaticReSourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.cache = c
	}
}

func staticWithMoreExtension(extMap map[string]string) StaticReSourceHandlerOption {
	return func(h *StaticResourceHandler) {
		for ext, contentType := range extMap {
			h.extensionContextTypeMap[ext] = contentType
		}
	}
}

func (s StaticResourceHandler) Handle(ctx *Context) {
	//w无缓存
	//拿到目标文件
	//定位到
	//返回给前端

	file, err := ctx.PathValue("file")
	if err != nil {
		ctx.RespStatusCode = http.StatusBadRequest
		ctx.RespData = []byte("请求路径不对")
		return
	}
	dst := filepath.Join(s.dir, file)
	ext := filepath.Ext(dst)[1:]
	header := ctx.Resp.Header()
	if data, ok := s.cache.Get(file); ok {
		ext := filepath.Ext(dst)[1:]

		contentType := s.extensionContextTypeMap[ext]
		header.Set("Content-Type", contentType)
		header.Set("Content-length", strconv.Itoa(len(data.([]byte))))
		ctx.RespData = data.([]byte)
		ctx.RespStatusCode = http.StatusOK

	}

	data, err := os.ReadFile(dst)

	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("服务器错误")
		return
	}
	//大文件不加
	if len(data) <= s.maxSize {
		s.cache.Add(file, data)
	}

	s.cache.Add(file, data)
	contentType := s.extensionContextTypeMap[ext]
	header.Set("Content-Type", contentType)
	header.Set("Content-length", strconv.Itoa(len(data)))
	ctx.RespData = data
	ctx.RespStatusCode = http.StatusOK
}
