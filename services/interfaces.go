package services

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Serviceable interface {
	Start(context.Context, *errgroup.Group) error
	Stop() error
}
