package account

import "github.com/modernice/goes/codec"

func RegisterCommands(r codec.Registerer) {
	userRegisterCommands(r)
}
