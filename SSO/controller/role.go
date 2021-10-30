package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/UniqueStudio/UniqueSSO/common"
	"github.com/UniqueStudio/UniqueSSO/pb/sso"
	"github.com/UniqueStudio/UniqueSSO/pkg"
	"github.com/UniqueStudio/UniqueSSO/repo"
	"github.com/UniqueStudio/UniqueSSO/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

func UpdateUserRoleHandler(ctx *gin.Context) {
	apmCtx, span := util.Tracer.Start(ctx.Request.Context(), "UpdateUserRoleHandler")
	defer span.End()

	sess := sessions.Default(ctx)
	uid, ok := sess.Get(common.SESSION_NAME_UID).(string)
	if !ok {
		zapx.WithContext(apmCtx).Info("get uid from session failed")
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResp(errors.New("can't get uid from session")))
		return
	}

	data := pkg.RolePermissionReq{}
	if err := ctx.ShouldBindJSON(data); err != nil {
		zapx.WithContext(apmCtx).Error("bind request body failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResp(err))
		return
	}

	if err := repo.UpdateUserRole(apmCtx, uid, data.Role); err != nil {
		zapx.WithContext(apmCtx).Error("update user role failed", zap.Error(err))
	}

	ctx.JSON(http.StatusOK, pkg.SuccessResp(""))
}

func AddRolePermissionHandler(ctx *gin.Context) {
	apmCtx, span := util.Tracer.Start(ctx.Request.Context(), "AddRolePermissions")
	defer span.End()

	data := pkg.RolePermissionReq{}
	if err := ctx.ShouldBindJSON(&data); err != nil {
		zapx.WithContext(apmCtx).Error("bind request body failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResp(err))
		return
	}

	if err := repo.AddRolePermission(apmCtx, data.Role, &data.Permission); err != nil {
		zapx.WithContext(apmCtx).Error("add role permission failed", zap.Error(err), zap.Any("data", data))
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, pkg.SuccessResp(""))
}

func DeleteRolePermissionHandler(ctx *gin.Context) {
	apmCtx, span := util.Tracer.Start(ctx.Request.Context(), "AddRolePermissions")
	defer span.End()

	data := pkg.RolePermissionReq{}
	if err := ctx.ShouldBindJSON(&data); err != nil {
		zapx.WithContext(apmCtx).Error("bind request body failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResp(err))
		return
	}

	if err := repo.DeleteRolePermission(apmCtx, data.Role, &data.Permission); err != nil {
		zapx.WithContext(apmCtx).Error("delete role permission failed", zap.Error(err), zap.Any("data", data))
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, pkg.SuccessResp(""))
}

func GetRoleAllPermissions(ctx *gin.Context) {
	apmCtx, span := util.Tracer.Start(ctx.Request.Context(), "AddRolePermissions")
	defer span.End()

	rolestr := ctx.Query("role")
	role, err := strconv.ParseInt(rolestr, 10, 32)
	if err != nil {
		zapx.WithContext(apmCtx).Error("can't parse role query to int")
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResp(err))
		return
	}

	rps, err := repo.GetRoleAllPermissions(apmCtx, sso.Role(role))
	if err != nil || rps == nil {
		zapx.WithContext(apmCtx).Error("get role all permission failed", zap.Error(err), zap.Any("data", role))
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResp(err))
		return
	}

	perms := make([]*sso.Permission, len(*rps))
	for i := 0; i < len(perms); i++ {
		perms[i] = (*rps)[i].Permission
	}

	ctx.JSON(http.StatusOK, pkg.SuccessResp(perms))
}
