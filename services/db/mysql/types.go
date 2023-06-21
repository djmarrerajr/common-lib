package cockroach

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/djmarrerajr/common-lib/services/db"
)

var _ db.Adapter = new(MySqlDB)

type MySqlDB struct{}

func (d *MySqlDB) Connect()                             {}
func (d *MySqlDB) Disconnect()                          {}
func (d *MySqlDB) CreateAccount(acct *db.Account) error { return nil }
func (d *MySqlDB) GetAccount(acct *db.Account) error    { return nil }
func (d *MySqlDB) UpdateAccount(acct *db.Account) error { return nil }
func (d *MySqlDB) DeleteAccount(acct *db.Account) error { return nil }

func (d *MySqlDB) Start(ctx context.Context, grp *errgroup.Group) error { return nil }
func (d *MySqlDB) Stop() error                                          { return nil }
