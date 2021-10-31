package controller

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/UniqueStudio/UniqueSSO/common"
	"github.com/UniqueStudio/UniqueSSO/pb/sso"
	"github.com/UniqueStudio/UniqueSSO/pkg"
	"github.com/UniqueStudio/UniqueSSO/repo"
	"github.com/UniqueStudio/UniqueSSO/service"
	"github.com/UniqueStudio/UniqueSSO/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

// for first login, append user info. otherwise just login
func LarkOauthCallbackHandler(ctx *gin.Context) {
	apmCtx, span := util.Tracer.Start(ctx.Request.Context(), "LarkOauthCallbackHandler")
	defer span.End()

	code, ok := ctx.GetQuery("code")
	if !ok {
		ctx.JSON(http.StatusBadGateway, pkg.ErrorResp(errors.New("no code in redirect url")))
		return
	}

	redirectURI := ctx.Query("state")
	if redirectURI != "" {
		uri, err := url.PathUnescape(redirectURI)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, pkg.ErrorResp(err))
		}
		redirectURI = uri
	}

	eid, _, err := util.LarkAuthCode2IDToken(apmCtx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResp(err))
		return
	}

	larkUserInfo, err := util.GetLarkContactUserInfo(apmCtx, eid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResp(err))
		return
	}

	userGroups := make([]sso.Group, len(larkUserInfo.DepartmentIds))
	for i := range userGroups {
		dpinfo, err := util.GetLarkDepartmentInfo(apmCtx, larkUserInfo.DepartmentIds[i])
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, pkg.ErrorResp(err))
			return
		}
		userGroups[i] = common.LARK_GROUP_NAME_2_GROUP(dpinfo.Name)
	}

	zapx.Info("", zap.Any("larkUserInfo", larkUserInfo))

	larkDetail, err := util.MarshelLarkExternalInfo(larkUserInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.InternalError(err))
		return
	}

	role := util.LarkEmployeeType2Role(larkUserInfo.EmployeeType)

	uid, err := repo.GetUserIDByEid(apmCtx, larkUserInfo.UnionID)
	// it means have not logged in. append user info into database
	if err != nil || uid == "" {
		user := &sso.User{
			Name:      larkUserInfo.Name,
			Phone:     strings.Replace(larkUserInfo.Mobile, "+86", "", 1),
			Gender:    util.LarkGender2UserGender(larkUserInfo.Gender),
			Email:     larkUserInfo.Email,
			Group:     userGroups,
			Role:      role,
			AvatarURL: larkUserInfo.Avatar.AvatarOrigin,
			ExternalInfos: []*sso.ExternalInfo{
				{EName: sso.ExternalType_LARK, EID: larkUserInfo.UnionID, Detail: larkDetail},
			},
		}

		err = repo.SaveUser(apmCtx, user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, pkg.InternalError(err))
			return
		}

		service.PushTOTPToken2User(apmCtx, larkUserInfo.UnionID, user.TOTPSecret)
	}

	// TODO: update lark user info async
	// go service.UpdateUserInfo(&sso.User{
	// 	Name:      larkUserInfo.Name,
	// 	Phone:     strings.Replace(larkUserInfo.Mobile, "+86", "", 1),
	// 	Email:     larkUserInfo.Email,
	// 	Group:     userGroups,
	// 	Gender:    util.LarkGender2UserGender(larkUserInfo.Gender),
	// 	Role:      role,
	// 	AvatarURL: larkUserInfo.Avatar.AvatarOrigin,
	// 	ExternalInfos: []*sso.ExternalInfo{
	// 		{EName: sso.ExternalType_LARK, EID: larkUserInfo.UnionID, Detail: larkDetail},
	// 	},
	// })

	// issue session
	sess := sessions.Default(ctx)
	sess.Set(common.SESSION_NAME_UID, uid)
	sess.Options(sessions.Options{MaxAge: common.SESSION_MAX_AGE})
	sess.Save()

	if redirectURI != "" {
		ctx.Redirect(http.StatusFound, redirectURI)
		return
	}

	ctx.JSON(http.StatusOK, pkg.AuthSuccess(&sso.User{
		UID: uid,
	}))
}
