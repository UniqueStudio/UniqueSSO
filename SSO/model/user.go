package model

import (
	"time"

	"github.com/UniqueStudio/UniqueSSO/pb/lark"
	"github.com/UniqueStudio/UniqueSSO/pb/sso"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type BasicUserInfo struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	*sso.User
}

type UserGroup struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID string        `gorm:"column:uid;primaryKey"`
	Groups pq.Int64Array `gorm:"type:integer[]"`
}

type RolePermission struct {
	gorm.Model `json:"-"`

	Role sso.Role `gorm:"column:role"`
	*sso.Permission
}

type LarkExternalInfo struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	*sso.ExternalInfo
	*lark.LarkUserInfo
}

type UserExternalInfo struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID      string         `gorm:"column:uid;primaryKey"`
	ExternalIDs pq.StringArray `gorm:"column:eids;type:text[]"`
}

func (BasicUserInfo) TableName() string {
	return "user"
}

func (UserGroup) TableName() string {
	return "user_groups"
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (LarkExternalInfo) TableName() string {
	return "lark_external_info"
}

func (UserExternalInfo) TableName() string {
	return "user_external_ids"
}
