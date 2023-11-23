package main

import (
	"github.com/startcodextech/goauth/internal/application"
	"github.com/startcodextech/goauth/internal/domain/account"
)

func main() {

	setup := application.New("")

	defer setup.Cancel()

	setup.Events()
	defer setup.Disconnect()

	setup.Commands()

	repo := setup.Aggregates()

	commandErrors := account.UserHandleCommands(setup.Context(), setup.CommandBus(), repo)

	setup.Rest()

	application.LogErrors(setup.Context(), commandErrors)

}
