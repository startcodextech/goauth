package users

//go:generate buf generate

//go:generate mockery --quiet --dir ./services/users -r --all --inpackage --case underscore

//go:generate swag init
