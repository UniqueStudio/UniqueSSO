package pkg

import "github.com/UniqueStudio/UniqueSSO/pb/lark"

type LarkAppTokenReq struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type LarkBasicResp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type LarkAppTokenResp struct {
	LarkBasicResp
	AppAccessToken string `json:"app_access_token"`
	Expire         int    `json:"expire"`
}

type LarkTenantTokenReq struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type LarkTenantTokenResp struct {
	Code              int    `json:"code"`
	Message           string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

type LarkCode2TokenReq struct {
	GrantType string `json:"grant_type"`
	Code      string `json:"code"`
}

type LarkCode2TokenResp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		lark.LarkUserInfo
	} `json:"data"`
}

type LarkGetUserInfoResp struct {
	Code    int               `json:"code"`
	Message string            `json:"msg"`
	Data    lark.LarkUserInfo `json:"data"`
}

type LarkGetContactUserInfoResp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		User lark.LarkUserInfo `json:"user"`
	} `json:"data"`
}

type LarkDeleteContactUserInfoResp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type LarkGetDepartmentInfoResp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		Department lark.Department `json:"department"`
	} `json:"data"`
}

type LarkBotPushMessage struct {
	ReceiveID   string `json:"receive_id"`
	Content     string `json:"content"`
	MessageType string `json:"msg_type"`
}
