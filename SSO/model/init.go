package model

import (
	"github.com/UniqueStudio/UniqueSSO/database"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

func InitTables() (err error) {
	err = database.DB.AutoMigrate(&BasicUserInfo{})
	if err != nil {
		zapx.Error("migrate basic user info failed", zap.Error(err))
		return err
	}

	err = database.DB.AutoMigrate(&UserGroup{})
	if err != nil {
		zapx.Error("migrate user group failed", zap.Error(err))
		return err
	}

	err = database.DB.AutoMigrate(&RolePermission{})
	if err != nil {
		zapx.Error("migrate role permission failed", zap.Error(err))
		return err
	}

	err = database.DB.AutoMigrate(&LarkExternalInfo{})
	if err != nil {
		zapx.Error("migrate lark external info failed", zap.Error(err))
		return err
	}

	err = database.DB.AutoMigrate(&UserExternalInfo{})
	if err != nil {
		zapx.Error("migrate user external info failed", zap.Error(err))
		return err
	}

	return nil
}
