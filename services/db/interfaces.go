package db

import "github.com/djmarrerajr/common-lib/services"

type Adapter interface {
	services.Serviceable

	CreateAccount(*Account) error
	GetAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(*Account) error
}
