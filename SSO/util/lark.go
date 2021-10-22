package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/UniqueStudio/UniqueSSO/common"
	"github.com/UniqueStudio/UniqueSSO/conf"
	"github.com/UniqueStudio/UniqueSSO/database"
	"github.com/UniqueStudio/UniqueSSO/model"
	"github.com/UniqueStudio/UniqueSSO/pb/lark"
	"github.com/UniqueStudio/UniqueSSO/pb/sso"
	"github.com/UniqueStudio/UniqueSSO/pkg"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func LarkAuthCode2IDToken(ctx context.Context, authCode string) (string, string, error) {
	apmCtx, span := Tracer.Start(ctx, "GetLarkUserToken")
	defer span.End()

	bs, err := json.Marshal(pkg.LarkCode2TokenReq{
		GrantType: "authorization_code",
		Code:      authCode,
	})
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequestWithContext(
		apmCtx,
		http.MethodPost,
		common.LARK_AUTH_CODE2TOKEN,
		bytes.NewBuffer(bs),
	)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("new get lark user token request failed", zap.Error(err))
		return "", "", err
	}

	tt, err := GetLarkAppToken(apmCtx)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", "Bearer "+tt)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("do get lark_user_token failed", zap.Error(err))
		return "", "", err
	}

	data := new(pkg.LarkCode2TokenResp)
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("unmarshal lark_user_token failed", zap.Error(err))
		return "", "", err
	}

	if data.Data.AccessToken == "" {
		err := errors.New(data.Message)
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("lark_user_token is missing", zap.Error(err))
		return "", "", err
	}

	// cache user access token
	database.RedisClient.Set(apmCtx, common.REDIS_LARK_USER_TOKEN_KEY(data.Data.UnionID), data.Data.AccessToken, -1)

	return data.Data.UnionID, data.Data.AccessToken, nil
}

func GetLarkTenantToken(ctx context.Context) (string, error) {
	apmCtx, span := Tracer.Start(ctx, "GetLarkTenantToken")
	defer span.End()

	token, err := database.RedisClient.Get(apmCtx, common.REDIS_LARK_TENANT_TOKEN_KEY).Result()
	// meet cache
	if err == nil && token != "" {
		return token, nil
	}

	bs, err := json.Marshal(pkg.LarkTenantTokenReq{
		AppId:     conf.SSOConf.Lark.AppId,
		AppSecret: conf.SSOConf.Lark.AppSecret,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		common.LARK_TENANT_TOKEN,
		bytes.NewBuffer(bs),
	)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("new accessing lark_tenant_token request failed", zap.Error(err))
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("do accessing lark_tenant_token failed", zap.Error(err))
		return "", err
	}

	data := new(pkg.LarkTenantTokenResp)
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("unmarshal lark_tenant_token response failed", zap.Error(err))
		return "", err
	}

	if data.TenantAccessToken == "" {
		err := errors.New(data.Message)
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("get tenant access token failed", zap.Error(err))
		return "", err
	}

	database.RedisClient.Set(apmCtx, common.REDIS_LARK_TENANT_TOKEN_KEY, data.TenantAccessToken, time.Duration(data.Expire)*time.Second)

	return data.TenantAccessToken, nil
}

func GetLarkAppToken(ctx context.Context) (string, error) {
	apmCtx, span := Tracer.Start(ctx, "GetLarkAppToken")
	defer span.End()

	token, err := database.RedisClient.Get(apmCtx, common.REDIS_LARK_APP_TOKEN_KEY).Result()
	// meet cache
	if err == nil && token != "" {
		return token, nil
	}

	bs, err := json.Marshal(pkg.LarkAppTokenReq{
		AppId:     conf.SSOConf.Lark.AppId,
		AppSecret: conf.SSOConf.Lark.AppSecret,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		common.LARK_APP_TOKEN,
		bytes.NewBuffer(bs),
	)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("new accessing lark_app_token request failed", zap.Error(err))
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("do accessing lark_app_token failed", zap.Error(err))
		return "", err
	}

	data := new(pkg.LarkAppTokenResp)
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("unmarshal lark_app_token response failed", zap.Error(err))
		return "", err
	}

	if data.AppAccessToken == "" {
		err := errors.New(data.Message)
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("get app access token failed", zap.Error(err))
		return "", err
	}

	database.RedisClient.Set(apmCtx, common.REDIS_LARK_APP_TOKEN_KEY, data.AppAccessToken, time.Duration(data.Expire)*time.Second)

	return data.AppAccessToken, nil
}

func GetLarkUserBasicInfo(ctx context.Context, userToken string) (*lark.LarkUserInfo, error) {
	apmCtx, span := Tracer.Start(ctx, "GetLarkUserInfo")
	defer span.End()

	req, err := http.NewRequestWithContext(
		apmCtx,
		http.MethodGet,
		common.LARK_FETCH_USER_INFO,
		nil,
	)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("build fetch lark user info request failed", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+userToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("send fetch lark user info failed", zap.Error(err))
		return nil, err
	}

	data := new(pkg.LarkGetUserInfoResp)
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("unmarshal response failed", zap.Error(err))
		return nil, err
	}

	return &data.Data, nil
}

func GetLarkContactUserInfo(ctx context.Context, userId string) (*lark.LarkUserInfo, error) {
	apmCtx, span := Tracer.Start(ctx, "GetLarkUserContactInfo")
	defer span.End()

	req, err := http.NewRequestWithContext(
		apmCtx,
		http.MethodGet,
		common.LARK_FETCH_USER_CONTACT_INFO(userId),
		nil,
	)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("build fetch lark user contact info request failed", zap.Error(err))
		return nil, err
	}

	tt, err := GetLarkTenantToken(apmCtx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+tt)

	query := url.Values{
		"user_id_type": []string{"union_id"},
	}
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("send fetch lark user contact info failed", zap.Error(err))
		return nil, err
	}

	data := new(pkg.LarkGetContactUserInfoResp)
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("unmarshal response failed", zap.Error(err))
		return nil, err
	}

	return &data.Data.User, nil
}

func RemoveLarkUser(ctx context.Context, eid string) error {
	apmCtx, span := Tracer.Start(ctx, "RemoveLarkUser")
	defer span.End()

	req, err := http.NewRequestWithContext(
		apmCtx,
		http.MethodDelete,
		common.LARK_DELETE_USER(eid),
		nil,
	)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("build delete lark user info request failed", zap.Error(err))
		return err
	}

	tt, err := GetLarkTenantToken(apmCtx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+tt)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("send delete user contact info failed", zap.Error(err))
		return err
	}

	data := new(pkg.LarkDeleteContactUserInfoResp)
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("unmarshal response failed", zap.Error(err))
		return err
	}

	if data.Code != 0 {
		return errors.New(data.Message)
	}

	return nil
}

func MarshelLarkExternalInfo(extInfo *lark.LarkUserInfo) (*anypb.Any, error) {
	data := new(anypb.Any)
	if err := anypb.MarshalFrom(data, extInfo, proto.MarshalOptions{}); err != nil {
		return nil, err
	}
	return data, nil
}

func UnmarshalLarkExternalInfo(buf *anypb.Any) (*model.LarkExternalInfo, error) {
	larkExtInfo := new(lark.LarkUserInfo)
	if err := anypb.UnmarshalTo(buf, larkExtInfo, proto.UnmarshalOptions{}); err != nil {
		return nil, err
	}

	return &model.LarkExternalInfo{
		ExternalInfo: sso.ExternalInfo{
			EName: sso.ExternalType_Lark,
			EID:   larkExtInfo.UnionID,
		},
		LarkUserInfo: *larkExtInfo,
	}, nil
}

func LarkDepartmentIds2Groups(ids []string) ([]sso.Group, error) {
	groups := make([]sso.Group, len(ids))
	for i := range ids {
		groups[i] = conf.SSOConf.Lark.GroupIdNameMap[ids[i]]
	}
	return groups, nil
}

func LarkEmployeeType2UserType(typ int32) sso.UserType {
	userType, ok := conf.SSOConf.Lark.LarkEnumTypeMap[typ]
	if ok {
		return userType
	}
	return sso.UserType_Intern
}

func LarkGender2UserGender(gender lark.Gender) sso.Gender {
	switch gender {
	case lark.Gender_Male:
		return sso.Gender_Male
	case lark.Gender_Female:
		return sso.Gender_Female
	default:
		return sso.Gender_invalidGender
	}
}
