package web

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/ifs"
	"bytes"
	"errors"
	"fmt"
	"github.com/mijia/sweb/log"
	"github.com/mijia/sweb/render"
	"github.com/mijia/sweb/server"
	"golang.org/x/net/context"
	"html/template"
	"io"
	"net/http"
	"runtime"
	"time"
)

type Server struct {
	*server.Server
	render         *render.Render
	confPath       string
	isDebug        bool
	muxControllers []ifs.MuxController
}

func (s *Server) AddMuxController(mcs ...ifs.MuxController) {
	s.muxControllers = append(s.muxControllers, mcs...)
}

func (s *Server) ListenAndServe(addr string) error {
	s.AddMuxController(&UserController{})
	s.AddMuxController(&ReservationController{})

	s.render = s.initRender()

	ignoredUrls := []string{"/javascripts/", "/images/", "/stylesheets/", "/fonts/", "/debug/vars", "/favicon", "/robots"}
	s.Middleware(NewRecoveryWare(s.isDebug))
	s.Middleware(server.NewStatWare(ignoredUrls...))
	s.Middleware(server.NewRuntimeWare(ignoredUrls, true, 15*time.Minute))
	s.EnableExtraAssetsJson("assets_map.json")

	// Change the asset prefix to CDN host for external prod server
	//if config.Instance().Env != "staging" && !config.Instance().IsSmokeServer() {
	//	s.EnableAssetsPrefix("https://assets.91zhiwang.com")
	//} else if config.Instance().IsSmokeServer() {
	//	s.EnableAssetsPrefix(config.Instance().H5Host)
	//}

	s.Get("/debug/vars", "RuntimeStat", s.getRuntimeStat)
	s.Files("/assets/*filepath", http.Dir("public"))

	for _, mc := range s.muxControllers {
		mc.SetResponseRenderer(s)
		mc.SetUrlReverser(s)
		mc.MuxHandlers(s)
	}

	return s.Run(addr)
}

func (s *Server) initRender() *render.Render {
	tSets := []*render.TemplateSet{}
	for _, mc := range s.muxControllers {
		mcTSets := mc.GetTemplates()
		tSets = append(tSets, mcTSets...)
	}
	r := render.New(render.Options{
		Directory:     "templates",
		Funcs:         s.renderFuncMaps(),
		Delims:        render.Delims{"{{", "}}"},
		IndentJson:    true,
		UseBufPool:    true,
		IsDevelopment: s.isDebug,
	}, tSets)
	log.Info("Templates loaded ...")
	return r
}

func (s *Server) getRuntimeStat(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	http.DefaultServeMux.ServeHTTP(w, r)
	return ctx
}

func formatTime(tm time.Time, layout string) string {
	return tm.Format(layout)
}

func (s *Server) renderFuncMaps() []template.FuncMap {
	funcs := make([]template.FuncMap, 0)
	funcs = append(funcs, s.DefaultRouteFuncs())
	funcs = append(funcs, template.FuncMap{
		"add": func(input interface{}, toAdd int) float64 {
			switch t := input.(type) {
			case int:
				return float64(t) + float64(toAdd)
			case int64:
				return float64(t) + float64(toAdd)
			case int32:
				return float64(t) + float64(toAdd)
			case float32:
				return float64(t) + float64(toAdd)
			case float64:
				return t + float64(toAdd)
			default:
				return float64(toAdd)
			}
		},
		"formatTime": formatTime,
	})
	return funcs
}

func (s *Server) RenderJsonOr500(w http.ResponseWriter, status int, v interface{}) {
	s.renderJsonOr500(w, status, v)
}

func (s *Server) renderJsonOr500(w http.ResponseWriter, status int, v interface{}) {
	if err := s.render.Json(w, status, v); err != nil {
		log.Errorf("Server got a json rendering error, %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) RenderHtmlOr500(w http.ResponseWriter, status int, name string, binding interface{}) {
	s.renderHtmlOr500(w, status, name, binding)
}

func (s *Server) RenderRawHtml(w http.ResponseWriter, status int, htmlString string) {
	s.renderString(w, status, htmlString)
}

func (s *Server) renderHtmlOr500(w http.ResponseWriter, status int, name string, binding interface{}) {
	w.Header().Set("Cache-Control", "no-store, no-cache")
	if err := s.render.Html(w, status, name, binding); err != nil {
		log.Errorf("Server got a rendering error, %q", err)
		if s.isDebug {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			//渲染一个500错误页面
			s.RenderError500(w, err)
		}
	}
}

func (s *Server) renderString(w http.ResponseWriter, status int, data string) {
	out := new(bytes.Buffer)
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(status)
	out.Write([]byte(data))
	io.Copy(w, out)
}

func (s *Server) Get(path, name string, handle server.Handler) {
	newHandle := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (newCtx context.Context) {
		newCtx = ctx
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, s.isDebug)]
				msg := fmt.Sprintf("Request: %s \r\n PANIC: %s\n%s", r.URL.String(), err, stack)
				log.Error(msg)
				s.RenderError500(w, errors.New(msg))
			}
		}()
		newCtx = handle(ctx, w, r)
		return
	}
	s.Server.Get(path, name, newHandle)
}

func (s *Server) GetJson(path string, name string, handle ifs.JsonHandler) {
	s.Get(path, name, s.makeJsonHandler(handle))
}

func (s *Server) PostJson(path string, name string, handle ifs.JsonHandler) {
	s.Post(path, name, s.makeJsonHandler(handle))
}
func (s *Server) makeJsonHandler(handle ifs.JsonHandler) server.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		// request log
		switch r.Method {
		case http.MethodGet:
			if r.URL.RawQuery != "" {
				log.Infof("%s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
			} else {
				log.Infof("%s %s", r.Method, r.URL.Path)
			}
		case http.MethodPost:
			r.ParseForm()
			log.Infof("%s %s %+v", r.Method, r.URL.Path, r.PostForm)
		}

		resp := handle(ctx, w, r)

		// response log
		if r.Method != http.MethodGet {
			log.Infof("response: %s %+v", r.URL.Path, resp)
		}

		s.renderJsonOr500(w, 200, resp)
		return ctx
	}
}

func NewServer(confPath string, isDebug bool) *Server {
	if isDebug {
		log.EnableDebug()
	}
	srv := &Server{
		confPath:       confPath,
		isDebug:        isDebug,
		muxControllers: []ifs.MuxController{},
	}

	ctx := context.Background()
	srv.Server = server.New(ctx, isDebug)

	// utils.GetAllWxUserAndUpdate(buslogic.WX_ACCOUNT_LICAI)
	// allWxUserTicker := time.NewTicker(time.Minute * 10)
	// go func(){
	// 	utils.GetAllWxUserAndUpdate(buslogic.WX_ACCOUNT_LICAI)
	// 	for  t := range allWxUserTicker.C {
	// 		logger.Info("get all wx user %+v", t)
	// 		utils.GetAllWxUserAndUpdate(buslogic.WX_ACCOUNT_LICAI)
	// 	}
	// }()
	return srv
}

type BaseHandler struct {
	rr ifs.ResponseRenderer
	s  ifs.UrlReverser
}

func (d *BaseHandler) SetResponseRenderer(rr ifs.ResponseRenderer) {
	d.rr = rr
}

func (d *BaseHandler) SetUrlReverser(s ifs.UrlReverser) {
	d.s = s
}
