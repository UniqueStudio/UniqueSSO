package pkg

import "context"

// now there are lark user maintainer and qq group maintainer
type Kicker interface {
	KickUser(ctx context.Context, userId, groupId string) error
}

type UserMaintainer interface {
	RegisterKicker(Kicker) error
}
