package service

import (
	"net/http"

	"github.com/UniqueStudio/UniqueSSO/pb/sso"
)

func HTTPMethod2Action(method string) sso.Action {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		return sso.Action_FULL
	default:
		return sso.Action_READ
	}
}
