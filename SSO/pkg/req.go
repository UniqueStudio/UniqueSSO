package pkg

import (
	"github.com/UniqueStudio/UniqueSSO/pb/sso"
)

type QrCodeStatus struct {
	Status   string `json:"status"`
	AuthCode string `json:"auth_code"`
}

type LoginUser struct {
	Phone        string `json:"phone,omitempty"`
	Email        string `json:"email,omitempty"`
	Password     string `json:"password,omitempty"`
	TOTPPasscode string `json:"totp_token,omitempty"`
	Code         string `json:"code,omitempty"`
}

type RolePermissionReq struct {
	Role       sso.Role       `json:"role"`
	Permission sso.Permission `json:"permission"`
}
