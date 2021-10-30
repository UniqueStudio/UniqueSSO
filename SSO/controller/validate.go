package controller

import (
	"net/http"

	"github.com/UniqueStudio/UniqueSSO/common"
	"github.com/UniqueStudio/UniqueSSO/conf"
	"github.com/UniqueStudio/UniqueSSO/pb/sso"
	"github.com/UniqueStudio/UniqueSSO/pkg"
	"github.com/UniqueStudio/UniqueSSO/repo"
	"github.com/UniqueStudio/UniqueSSO/service"
	"github.com/UniqueStudio/UniqueSSO/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xylonx/zapx"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func TraefikAuthValidate(ctx *gin.Context) {
	apmCtx, span := util.Tracer.Start(ctx.Request.Context(), "TraefikAuthValidate")
	defer span.End()

	sess := sessions.Default(ctx)

	uid, ok := sess.Get(common.SESSION_NAME_UID).(string)
	if !ok {
		zapx.WithContext(apmCtx).Info("get uid from session failed")
		ctx.Redirect(http.StatusFound, conf.SSOConf.Application.LoginRedirectURI)
		return
	}

	zapx.WithContext(apmCtx).Info("get uid from session successfully", zap.String("uid", uid))
	span.SetAttributes(attribute.String("UID", uid))

	user, err := repo.GetBasicUserById(apmCtx, uid)
	if err != nil {
		zapx.WithContext(apmCtx).Error("get user by uid failed", zap.String("uid", uid))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	switch user.Role {
	case sso.Role_LEADER, sso.Role_DEVOPS:
		ctx.JSON(http.StatusOK, "ok")
		return
	}

	uri := ctx.Request.Header.Get("X-Forwarded-Uri")
	proto := ctx.Request.Header.Get("X-Forwarded-Proto")
	host := ctx.Request.Header.Get("X-Forwarded-Host")
	resource := proto + "://" + host + uri
	method := ctx.Request.Header.Get("X-Forwarded-Method")
	action := service.HTTPMethod2Action(method)

	err = repo.RoleHasPermission(apmCtx, user.Role, &sso.Permission{
		Action:   action,
		Resource: resource,
	})
	if err != nil {
		zapx.WithContext(apmCtx).Error("user don't have the permission to operation resource")
		ctx.JSON(http.StatusForbidden, pkg.ErrorResp(err))
		return
	}

	ctx.Writer.Header().Set("X-UID", uid)
	ctx.JSON(http.StatusOK, "")
	return
}
