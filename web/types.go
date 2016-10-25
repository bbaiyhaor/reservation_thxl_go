package web

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/buslogic"
	"github.com/mijia/sweb/render"
	"github.com/mijia/sweb/server"
	"golang.org/x/net/context"
	"net/http"
)

type JsonHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{})

type JsonMuxer interface {
	server.Muxer
	GetJson(path string, name string, handle JsonHandler)
	PostJson(path string, name string, handle JsonHandler)
	PutJson(path string, name string, handle JsonHandler)
	PatchJson(path string, name string, handle JsonHandler)
	DeleteJson(path string, name string, handle JsonHandler)
}

type ResponseRender interface {
	RenderJsonOr500(w http.ResponseWriter, status int, v interface{})
	RenderHtmlOr500(w http.ResponseWriter, status int, name string, binding interface{})
	RenderRawHtml(w http.ResponseWriter, status int, htmlString string)
	RenderError500(w http.ResponseWriter, err error)
	RenderError404(w http.ResponseWriter)
}

type UrlReverser interface {
	Reverse(name string, params ...interface{}) string
	Assets(path string) string
}

type MuxController interface {
	MuxHandlers(m JsonMuxer)
	SetResponseRenderer(r ResponseRender)
	SetUrlReverser(r UrlReverser)
	GetTemplates() []*render.TemplateSet
}

type BaseMuxController struct {
	ResponseRender
	UrlReverser
}

func (b *BaseMuxController) GetTemplates() []*render.TemplateSet {
	return nil
}

func (b *BaseMuxController) SetResponseRenderer(rr ResponseRender) {
	b.ResponseRender = rr
}

func (b *BaseMuxController) SetUrlReverser(ur UrlReverser) {
	b.UrlReverser = ur
}

type JsonData struct {
	Status  string      `json:"status"`
	ErrCode int         `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Payload interface{} `json:"payload"`
}

func wrapJsonOk(payload interface{}) JsonData {
	return JsonData{
		Status:  "OK",
		Payload: payload,
	}
}

func wrapJsonError(msg string, payloads ...interface{}) JsonData {
	if msg == "" {
		msg = "服务器开小差了，请稍候重试！"
	}
	data := JsonData{
		Status:  "FAIL",
		ErrCode: -1,
		ErrMsg:  msg,
	}
	if msg == buslogic.CHECK_FORCE_ERROR {
		data = JsonData{
			Status: buslogic.CHECK_FORCE_ERROR,
		}
	}
	if len(payloads) > 0 {
		data.Payload = payloads[0]
	}
	return data
}
