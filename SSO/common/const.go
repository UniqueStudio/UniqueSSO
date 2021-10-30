package common

import (
	"fmt"
	"regexp"
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

func LARK_DEPARTMENT_INFO(id string) string {
	return "https://open.feishu.cn/open-apis/contact/v3/departments/" + id
}

var (
	aiRegexp      *regexp.Regexp
	androidRegexp *regexp.Regexp
	designRegexp  *regexp.Regexp
	gameRegexp    *regexp.Regexp
	labRegexp     *regexp.Regexp
	pmRegexp      *regexp.Regexp
	webRegexp     *regexp.Regexp
	iOSRegexp     *regexp.Regexp
)

func LARK_GROUP_NAME_2_GROUP(name string) sso.Group {
	switch {
	case aiRegexp.MatchString(name):
		return sso.Group_AI
	case androidRegexp.MatchString(name):
		return sso.Group_ANDROID
	case designRegexp.MatchString(name):
		return sso.Group_DESIGN
	case gameRegexp.MatchString(name):
		return sso.Group_GAME
	case labRegexp.MatchString(name):
		return sso.Group_LAB
	case pmRegexp.MatchString(name):
		return sso.Group_PM
	case webRegexp.MatchString(name):
		return sso.Group_WEB
	case iOSRegexp.MatchString(name):
		return sso.Group_IOS
	default:
		return sso.Group_invalidGroup
	}
}

func init() {
	aiRegexp, _ = regexp.Compile(`AI.*`)
	androidRegexp, _ = regexp.Compile(`Android.*`)
	designRegexp, _ = regexp.Compile(`Design.*`)
	gameRegexp, _ = regexp.Compile(`Game.*`)
	labRegexp, _ = regexp.Compile(`Lab.*`)
	pmRegexp, _ = regexp.Compile(`PM.*`)
	webRegexp, _ = regexp.Compile(`Web.*`)
	iOSRegexp, _ = regexp.Compile(`iOS.*`)
}
