package internal

import "github.com/UniqueStudio/UniqueSSO/pkg"

type ExternalUserMaintainer struct {
	kickers []pkg.Kicker
}

var _ pkg.UserMaintainer = &ExternalUserMaintainer{}

var Maintainer pkg.UserMaintainer

func SetupUserMaintainer() error {
	Maintainer = &ExternalUserMaintainer{
		kickers: make([]pkg.Kicker, 0, 2),
	}
	return nil
}

func (m *ExternalUserMaintainer) RegisterKicker(k pkg.Kicker) error {
	m.kickers = append(m.kickers, k)
	return nil
}
