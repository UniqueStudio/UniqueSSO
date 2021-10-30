package router

import (
	"github.com/UniqueStudio/UniqueSSO/controller"
	"github.com/UniqueStudio/UniqueSSO/middleware"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// global middleware
	v1r := r.Group("/v1")
	v1r.Use(middleware.TracingMiddleware())
	v1r.Use(middleware.Cors())
	v1r.Use(gin.Recovery())
	v1r.Use(gin.Logger())
	v1r.Use(sessions.Sessions("SSO_SESSION", middleware.RedisSessionStore))

	// traefik validate
	v1r.GET("/gateway/validate/traefik", controller.TraefikAuthValidate)

	// sms
	smsrouter := v1r.Group("/sms")
	smsrouter.POST("code", controller.SendSmsCode)

	// normal login
	router := v1r.Group("")
	router.Use(middleware.SessionRedirect())
	router.POST("/login", controller.Login)
	router.POST("/logout", controller.Logout)

	rolerouter := router.Group("/role")
	rolerouter.PUT("/user", controller.UpdateUserRoleHandler)
	rolerouter.GET("/permission", controller.GetRoleAllPermissions)
	rolerouter.POST("/permission", controller.AddRolePermissionHandler)
	rolerouter.DELETE("/permission", controller.DeleteRolePermissionHandler)

	// oauth login
	oauth := v1r.Group("/login/oauth")
	oauth.GET("/lark", controller.LarkOauthCallbackHandler)
}
