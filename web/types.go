package web

import (
	"github.com/mijia/sweb/render"
	"github.com/mijia/sweb/server"
	"github.com/shudiwsh2009/reservation_thxl_go/buslogic"
	re "github.com/shudiwsh2009/reservation_thxl_go/rerror"
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
		ErrCode: re.OK,
		ErrMsg:  re.ReturnMessage(re.OK),
		Payload: payload,
	}
}

func wrapJsonError(err error, payloads ...interface{}) JsonData {
	if err == nil {
		err = re.NewRErrorCode("", nil, re.ERROR_UNKNOWN)
	}
	var data JsonData
	rerr, ok := err.(*re.RError)
	if ok {
		data = JsonData{
			Status:  "FAIL",
			ErrCode: rerr.Code(),
			ErrMsg:  rerr.DisplayMessage(),
		}
		if rerr.Code() == re.CHECK {
			data.Status = buslogic.CHECK_FORCE_ERROR
		}
	} else {
		data = JsonData{
			Status:  "FAIL",
			ErrCode: re.ERROR_UNKNOWN,
			ErrMsg:  re.ReturnMessage(re.ERROR_UNKNOWN),
		}
	}
	if len(payloads) > 0 {
		data.Payload = payloads[0]
	}
	return data
}
