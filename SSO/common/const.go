package common

import (
	"fmt"
	"time"

	"github.com/UniqueStudio/UniqueSSO/conf"
	"github.com/UniqueStudio/UniqueSSO/pb/sso"
)

const (
	SignTypePhonePassword = "phone"
	SignTypePhoneSms      = "sms"
	SignTypeEmailPassword = "email"
	SignTypeLark          = "lark"
)

const (
	CASErrInvalidRequest    = "INVALID_REQUEST"
	CASErrInvalidTicketSpec = "INVALID_TICKET_SPEC"
	CASErrInvalidTicket     = "INVALID_TICKET"
	CASErrInvalidService    = "INVALID_SERVICE"
	CASErrInternalError     = "INTERNAL_ERROR"
	CASErrUnauthorized      = "UNAUTHENTICATED"
)

const (
	DebugMode = "debug"
)

const (
	SESSION_NAME_UID = "UID"
	SESSION_MAX_AGE  = 4 * 60 * 60 * 1000
	// CAS_COOKIE_NAME    = "CASTGC"
	// CAS_TGT_EXPIRES    = time.Hour
	// CAS_TICKET_EXPIRES = time.Minute * 5
	// DEFAULT_TIMEOUT    = 10000000

	SMS_CODE_EXPIRES = time.Minute * 3
)

const (
	SMS_TEMPO_CODE = "verificationCode"
)

const (
	EXTERNAL_NAME_LARK = "lark"
)

const (
	REDIS_LARK_TENANT_TOKEN_KEY = "LARK_SSO_TENANT_ACCESS_TOKEN"
	REDIS_LARK_APP_TOKEN_KEY    = "LARK_SSO_APP_TOKEN"
)

func REDIS_LARK_USER_TOKEN_KEY(unionId string) string {
	return "LARK_SSO_USER_TOKEN:" + unionId
}

const (
	LARK_OAUTH = "https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=%s&app_id=%s&state=%s"

	LARK_AUTH_CODE2TOKEN    = "https://open.feishu.cn/open-apis/authen/v1/access_token"
	LARK_USER_TOKEN_REFRESH = "https://open.feishu.cn/open-apis/authen/v1/refresh_access_token"
	LARK_TENANT_TOKEN       = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	LARK_APP_TOKEN          = "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal"

	LARK_FETCH_USER_INFO = "https://open.feishu.cn/open-apis/authen/v1/user_info"
)

func LARK_OAUTH_URL(state string) string {
	return fmt.Sprintf(LARK_OAUTH,
		conf.SSOConf.Lark.RedirectUri,
		conf.SSOConf.Lark.AppId,
		state,
	)
}

func LARK_FETCH_USER_CONTACT_INFO(id string) string {
	return "https://open.feishu.cn/open-apis/contact/v3/users/" + id
}

func LARK_DELETE_USER(id string) string {
	return "https://open.feishu.cn/open-apis/contact/v3/users/" + id
}

var (
	DEFAULT_PERMISSION = []*sso.Permission{
		{Action: sso.Action_READ, Resource: sso.Resource_BBS_REPORT, Scope: sso.Scope_ALL},
	}
)

var (
	USER_TYPE_PERMISSION map[sso.UserType][]*sso.Permission
)

func init() {
	USER_TYPE_PERMISSION = make(map[sso.UserType][]*sso.Permission)
	USER_TYPE_PERMISSION[sso.UserType_Intern] = []*sso.Permission{
		{Action: sso.Action_READ, Resource: sso.Resource_BBS_REPORT, Scope: sso.Scope_ALL},
		{Action: sso.Action_WRITE, Resource: sso.Resource_BBS_REPORT, Scope: sso.Scope_SELF},
	}
	USER_TYPE_PERMISSION[sso.UserType_Regular] = []*sso.Permission{
		{Action: sso.Action_READ, Resource: sso.Resource_BBS_REPORT, Scope: sso.Scope_ALL},
		{Action: sso.Action_WRITE, Resource: sso.Resource_BBS_REPORT, Scope: sso.Scope_SELF},
	}
	USER_TYPE_PERMISSION[sso.UserType_Graduated] = []*sso.Permission{
		{Action: sso.Action_READ, Resource: sso.Resource_BBS_REPORT, Scope: sso.Scope_ALL},
		{Action: sso.Action_WRITE, Resource: sso.Resource_BBS_REPORT, Scope: sso.Scope_SELF},
	}
}
