package repo

import (
	"context"
	"errors"

	"github.com/UniqueStudio/UniqueSSO/database"
	"github.com/UniqueStudio/UniqueSSO/model"
	"github.com/UniqueStudio/UniqueSSO/pb/sso"
	"github.com/UniqueStudio/UniqueSSO/util"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetUserIDByEid(ctx context.Context, eid string) (string, error) {
	return getUserIDByEid(ctx, database.DB, eid)
}

func GetBasicUserById(ctx context.Context, uid string) (*sso.User, error) {
	user, err := getBasicUserInfoById(ctx, database.DB, uid)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return &user.User, nil
}

func GetBasicUserByEmail(ctx context.Context, email string) (*sso.User, error) {
	apmCtx, span := util.Tracer.Start(ctx, "GetBasicUserByEmail")
	defer span.End()

	user := model.BasicUserInfo{}
	tx := database.DB.Where(&model.BasicUserInfo{
		User: sso.User{
			Email: email,
		},
	}).First(&user)
	if tx.RowsAffected == 0 {
		err := tx.Error
		if err == nil {
			err = gorm.ErrRecordNotFound
		}
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("get basic user by email failed", zap.Error(err))
		return nil, err
	}

	user.Password = ""
	return &user.User, nil
}

func GetBasicUserByPhone(ctx context.Context, phone string) (*sso.User, error) {
	apmCtx, span := util.Tracer.Start(ctx, "GetBasicUserByPhone")
	defer span.End()

	user := model.BasicUserInfo{}
	tx := database.DB.Where(&model.BasicUserInfo{
		User: sso.User{
			Phone: phone,
		},
	}).First(&user)
	if tx.RowsAffected == 0 {
		err := tx.Error
		if err == nil {
			err = gorm.ErrRecordNotFound
		}
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("get basic user by phone failed", zap.Error(err))
		return nil, err
	}

	user.Password = ""
	return &user.User, nil
}

func GetBasicUserByEID(ctx context.Context, eid string) (*sso.User, error) {
	apmCtx, span := util.Tracer.Start(ctx, "GetBasicUserByEID")
	defer span.End()

	tx := database.DB.Begin()
	uid, err := getUserIDByEid(ctx, tx, eid)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user, err := getBasicUserInfoById(ctx, tx, uid)
	if err != nil {
		if err == nil {
			err = gorm.ErrRecordNotFound
		}
		zapx.WithContext(apmCtx).Info("can't get basic user info by UID", zap.Error(err))
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	user.Password = ""
	return &user.User, nil
}

func GetBasicUserWithGroupById(ctx context.Context, uid string) (*sso.User, error) {
	apmCtx, span := util.Tracer.Start(ctx, "GetBasicUserWithGroupById")
	defer span.End()

	tx := database.DB.Begin()
	user := model.BasicUserInfo{}
	newTx := tx.Where(&model.BasicUserInfo{
		User: sso.User{
			UID: uid,
		},
	}).First(&user)
	if newTx.RowsAffected == 0 {
		err := newTx.Error
		if err == nil {
			err = gorm.ErrRecordNotFound
		}
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("get basic user by uid failed", zap.Error(err))
		tx.Rollback()
		return nil, err
	}

	group := model.UserGroup{}
	newTx = tx.Where(&model.UserGroup{
		UserID: uid,
	}).First(&group)
	if newTx.RowsAffected == 0 {
		err := newTx.Error
		if err == nil {
			err = gorm.ErrRecordNotFound
		}
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("get user groups by uid failed", zap.Error(err))
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	ssogroups := make([]sso.Group, len(group.Groups))
	for i := range ssogroups {
		ssogroups[i] = sso.Group(group.Groups[i])
	}
	user.Group = ssogroups

	user.Password = ""
	return &user.User, nil
}

// generate uid and totp secret while
func SaveUser(ctx context.Context, user *sso.User) error {
	apmCtx, span := util.Tracer.Start(ctx, "SaveUser")
	defer span.End()

	uid, err := uuid.NewUUID()
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("genrate uuid failed", zap.Error(err))
		return err
	}
	user.UID = uid.String()

	totpSecret, err := util.GenerateTOTPSharedKey()
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("generate totp shared key failed", zap.Error(err))
		return err
	}
	user.TOTPSecret = totpSecret

	// group table data
	groups := make(pq.Int64Array, len(user.Group))
	for i := range groups {
		groups[i] = int64(user.GetGroup()[i])
	}

	// permission table data
	perms := make([]model.UserPermission, len(user.Permissions))
	for i := range perms {
		perms[i] = model.UserPermission{
			UserID:     user.UID,
			Permission: *user.Permissions[i],
		}
	}

	// lark external user info data
	// ! handle external info
	// ! FIXME: define external name in IDL
	eids := make([]string, len(user.ExternalInfos))
	extInfos := make([]interface{}, len(user.ExternalInfos))
	for i := range user.ExternalInfos {
		eids[i] = user.ExternalInfos[i].EID
		switch user.ExternalInfos[i].EName {
		case sso.ExternalType_Lark:
			larkExtInfo, err := util.UnmarshalLarkExternalInfo(user.ExternalInfos[i].Detail)
			if err != nil {
				span.RecordError(err)
				zapx.WithContext(apmCtx).Error("unmarshal larkExtInfo failed", zap.Error(err))
				return err
			}
			extInfos[i] = larkExtInfo

		default:
			err := errors.New("unsupported external name")
			span.RecordError(err)
			zapx.WithContext(apmCtx).Error("unsupported external name")
			return err
		}
	}

	tx := database.DB.Begin()

	//insert basic user info
	newTx := tx.Create(&model.BasicUserInfo{User: *user})
	if newTx.RowsAffected == 0 {
		err := newTx.Error
		if err == nil {
			err = errors.New("no basic user info records were saved")
		}
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("no records were saved")
		tx.Rollback()
		return err
	}

	// insert user group information
	newTx = tx.Create(&model.UserGroup{
		UserID: user.UID,
		Groups: groups,
	})
	if newTx.RowsAffected == 0 {
		err := newTx.Error
		if err == nil {
			err = errors.New("no user group records were saved")
		}
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("no records were saved")
		tx.Rollback()
		return err
	}

	// insert user permission information
	newTx = tx.Create(&perms)
	if newTx.RowsAffected == 0 {
		err := newTx.Error
		if err == nil {
			err = errors.New("no user permissions records were saved")
		}
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("no records were saved")
		tx.Rollback()
		return err
	}

	// insert external information
	for i := range extInfos {
		newTx := tx.Create(extInfos[i])
		// newTx = tx.Table(extInfos[i].tableName).Create(i)
		if newTx.RowsAffected == 0 {
			err := newTx.Error
			if err == nil {
				err = errors.New("no external info records were saved")
			}
			span.RecordError(err)
			zapx.WithContext(apmCtx).Error("no records were saved")
			tx.Rollback()
			return err
		}
	}

	newTx = tx.Create(&model.UserExternalInfo{
		UserID:      user.UID,
		ExternalIDs: eids,
	})
	if newTx.RowsAffected == 0 {
		err := newTx.Error
		if err == nil {
			err = errors.New("no user external id relation records were saved")
		}
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("no records were saved")
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func getBasicUserInfoById(ctx context.Context, db *gorm.DB, uid string) (*model.BasicUserInfo, error) {
	apmCtx, span := util.Tracer.Start(ctx, "GetBasicUserById")
	defer span.End()

	user := model.BasicUserInfo{}
	tx := db.Where(&model.BasicUserInfo{
		User: sso.User{
			UID: uid,
		},
	}).First(&user)
	if tx.RowsAffected == 0 {
		err := tx.Error
		if err == nil {
			err = gorm.ErrRecordNotFound
		}
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("get basic user by uid failed", zap.Error(err))
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

func getUserIDByEid(ctx context.Context, db *gorm.DB, eid string) (string, error) {
	apmCtx, span := util.Tracer.Start(ctx, "getUserIDByEID")
	defer span.End()

	userExtInfo := model.UserExternalInfo{}
	newTx := db.Raw("SELECT * FROM user_external_ids WHERE ? = ANY (eids)", eid).Scan(&userExtInfo)
	if newTx.RowsAffected == 0 {
		err := newTx.Error
		if err == nil {
			err = gorm.ErrRecordNotFound
		}
		zapx.WithContext(apmCtx).Info("can't get eid info", zap.Error(err))
		return "", err
	}

	return userExtInfo.UserID, nil
}
