package web

import (
	"github.com/mijia/sweb/server"
	"golang.org/x/net/context"
	"net/http"
)

func RoleCookieInjection(handle func(http.ResponseWriter, *http.Request, string, int) (int, interface{})) JsonHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
		userId, _, userType, err := getUserCookie(r)
		if err != nil {
			return http.StatusOK, wrapJsonError("请先登录")
		}
		return handle(w, r, userId, userType)
	}
}

func LegacyAdminPageInjection(handle func(context.Context, http.ResponseWriter, *http.Request, string, int) context.Context) server.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		redirectUrl := "/reservation/admin/login"
		userId, _, userType, err := getUserCookie(r)
		if err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		}
		return handle(ctx, w, r, userId, userType)
	}
}
