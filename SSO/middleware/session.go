package middleware

import (
	"net/http"

	"github.com/UniqueStudio/UniqueSSO/common"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionRedirect() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sess := sessions.Default(ctx)

		service := ctx.Query("service")
		_, ok := sess.Get(common.SESSION_NAME_UID).(string)

		switch {
		// login with redirect
		case ok && service != "":
			ctx.Redirect(http.StatusFound, service)
			ctx.Abort()
			return

		// not login.
		case !ok && ctx.Request.URL.Path != "/login":
			ctx.Redirect(http.StatusFound, "/login")
			ctx.Abort()
			return

		default:
			ctx.Next()
			return
		}
	}
}
