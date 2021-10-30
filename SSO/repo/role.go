package repo

import (
	"context"
	"fmt"

	"github.com/UniqueStudio/UniqueSSO/database"
	"github.com/UniqueStudio/UniqueSSO/model"
	"github.com/UniqueStudio/UniqueSSO/pb/sso"
	"github.com/UniqueStudio/UniqueSSO/util"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

func RoleHasPermission(ctx context.Context, role sso.Role, perm *sso.Permission) error {
	apmCtx, span := util.Tracer.Start(ctx, "RoleHasPermission")
	defer span.End()

	rp := model.RolePermission{}
	newTx := database.DB.WithContext(apmCtx).Table(model.RolePermission{}.TableName()).Where(&model.RolePermission{
		Role:       role,
		Permission: perm,
	}).Scan(&rp)
	if newTx.RowsAffected == 0 {
		zapx.WithContext(apmCtx).Error("can't find permission for role", zap.String("role", sso.Role_name[int32(role)]), zap.Any("perm", perm))
		return fmt.Errorf("no record found. %v", newTx.Error)
	}
	return nil
}

func UpdateUserRole(ctx context.Context, uid string, role sso.Role) error {
	apmCtx, span := util.Tracer.Start(ctx, "UpdateUserRole")
	defer span.End()

	newTx := database.DB.WithContext(apmCtx).Where("uid = ?", uid).Update("role", role)
	if newTx.RowsAffected == 0 {
		return fmt.Errorf("update user role failed. %v", newTx.Error)
	}
	return nil
}

func AddRolePermission(ctx context.Context, role sso.Role, perm *sso.Permission) error {
	apmCtx, span := util.Tracer.Start(ctx, "AddRolePermission")
	defer span.End()

	newTx := database.DB.WithContext(apmCtx).Create(&model.RolePermission{
		Role:       role,
		Permission: perm,
	})

	if newTx.RowsAffected == 0 {
		return fmt.Errorf("add permission for role failed. %v", newTx.Error)
	}
	return nil
}

func DeleteRolePermission(ctx context.Context, role sso.Role, perm *sso.Permission) error {
	apmCtx, span := util.Tracer.Start(ctx, "DeleteRolePermission")
	defer span.End()

	newTx := database.DB.WithContext(apmCtx).
		Where("role = ? AND action = ? And resource = ?", role, perm.Action, perm.Resource).
		Delete(&model.RolePermission{})

	if newTx.RowsAffected == 0 {
		return fmt.Errorf("delete permission for role failed. ", newTx.Error)
	}

	return nil
}

func GetRoleAllPermissions(ctx context.Context, role sso.Role) (*[]model.RolePermission, error) {
	apmCtx, span := util.Tracer.Start(ctx, "GetRoleAllPermissions")
	defer span.End()

	rps := new([]model.RolePermission)
	newTx := database.DB.WithContext(apmCtx).Table(model.RolePermission{}.TableName()).Where("role = ?", role).Scan(rps)
	if newTx.Error != nil {
		return nil, newTx.Error
	}
	return rps, nil
}
