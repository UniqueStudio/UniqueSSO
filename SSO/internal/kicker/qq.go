package kicker

import (
	"context"

	"github.com/UniqueStudio/UniqueSSO/pkg"
)

type QQKicker struct {
}

var _ pkg.Kicker = &QQKicker{}

func (k *QQKicker) KickUser(ctx context.Context, userId, groupId string) error {
	return nil
}

func NewQQKicker() pkg.Kicker {
	return &QQKicker{}
}
