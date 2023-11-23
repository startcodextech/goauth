package account

import "github.com/modernice/goes/codec"

func RegisterEvents(r codec.Registerer) {
	userRegisterEvents(r)
}
