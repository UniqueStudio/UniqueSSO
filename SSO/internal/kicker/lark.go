package kicker

import (
	"context"

	"github.com/UniqueStudio/UniqueSSO/pkg"
	"github.com/UniqueStudio/UniqueSSO/util"
)

type LarkKicker struct {
}

var _ pkg.Kicker = &LarkKicker{}

func (k *LarkKicker) KickUser(ctx context.Context, userId, groupId string) error {
	return util.RemoveLarkUser(ctx, userId)
}

func NewLarkKicker() pkg.Kicker {
	return &LarkKicker{}
}
