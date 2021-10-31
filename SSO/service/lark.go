package service

import (
	"context"
	"fmt"

	"github.com/UniqueStudio/UniqueSSO/conf"
	"github.com/UniqueStudio/UniqueSSO/pb/sso"
	"github.com/UniqueStudio/UniqueSSO/pkg"
	"github.com/UniqueStudio/UniqueSSO/util"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

// TODO
func UpdateUserInfo(extInfo *sso.User) {
	// update ext user info

}

func PushTOTPToken2User(ctx context.Context, uid, token string) error {
	apmCtx, span := util.Tracer.Start(ctx, "PushTOTPToken")
	defer span.End()

	err := util.PushMessageByBot(apmCtx, &pkg.LarkBotPushMessage{
		ReceiveID:   uid,
		Content:     fmt.Sprintf(TOTPTokenMessageTemplate, token, conf.SSOConf.Application.LoginRedirectURI),
		MessageType: "interactive",
	})

	if err != nil {
		zapx.WithContext(apmCtx).Error("push totp token to user failed", zap.Error(err))
		return err
	}

	zapx.WithContext(apmCtx).Info("push totp token to user successfully")

	return nil
}

const TOTPTokenMessageTemplate = `{"config": {"wide_screen_mode": true},"i18n_elements": {"zh_cn": [{"tag": "markdown","content": "初次登陆，为您初始化TOTP Token：\n**%s**\n请注意保存，不要泄漏😃\n请前往[UniqueSSO](%s)初始化您的密码\n"},{"tag": "hr"},{"tag": "note","elements": [{"tag": "plain_text","content": "sso-dev"}]}]}}`
