package model

import (
	"time"

	"github.com/UniqueStudio/UniqueSSO/pb/lark"
	"github.com/UniqueStudio/UniqueSSO/pb/sso"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type BasicUserInfo struct {
	CreateAt time.Time      `json:"-"`
	UpdateAt time.Time      `json:"-"`
	DeleteAt gorm.DeletedAt `json:"-" gorm:"index"`

	sso.User
}

type UserGroup struct {
	CreateAt time.Time      `json:"-"`
	UpdateAt time.Time      `json:"-"`
	DeleteAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID string        `gorm:"column:uid;primaryKey"`
	Groups pq.Int64Array `gorm:"type:integer[]"`
}

type UserPermission struct {
	gorm.Model
	UserID string `gorm:"column:uid"`
	sso.Permission
}

type LarkExternalInfo struct {
	CreateAt time.Time      `json:"-"`
	UpdateAt time.Time      `json:"-"`
	DeleteAt gorm.DeletedAt `json:"-" gorm:"index"`

	sso.ExternalInfo
	lark.LarkUserInfo
}

type UserExternalInfo struct {
	CreateAt time.Time      `json:"-"`
	UpdateAt time.Time      `json:"-"`
	DeleteAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID      string         `gorm:"column:uid;primaryKey"`
	ExternalIDs pq.StringArray `gorm:"column:eids;type:text[]"`
}

func (BasicUserInfo) TableName() string {
	return "user"
}

func (UserGroup) TableName() string {
	return "user_groups"
}

func (UserPermission) TableName() string {
	return "user_permissions"
}

func (LarkExternalInfo) TableName() string {
	return "lark_external_info"
}

func (UserExternalInfo) TableName() string {
	return "user_external_ids"
}
