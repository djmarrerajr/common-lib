package cockroach

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	"github.com/djmarrerajr/common-lib/services/db"
	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

var _ db.Adapter = new(CockroachDB)

const connectionString = "postgresql://%s@%s:%s/%s?sslcert=%s&sslkey=%s&sslmode=verify-full&sslrootcert=%s"

type CockroachDB struct {
	conn *gorm.DB

	AppCtx shared.ApplicationContext
	logger utils.Logger

	host     string
	port     string
	user     string
	database string

	ca   string
	cert string
	key  string

	maxConn  int
	idleConn int
	maxTime  time.Duration
	idleTime time.Duration
}

func (d *CockroachDB) CreateAccount(acct *db.Account) error { return d.conn.Create(acct).Error }
func (d *CockroachDB) GetAccount(acct *db.Account) error    { return d.conn.First(acct).Error }
func (d *CockroachDB) UpdateAccount(acct *db.Account) error { return d.conn.Save(acct).Error }
func (d *CockroachDB) DeleteAccount(acct *db.Account) error { return d.conn.Delete(acct).Error }

func (d *CockroachDB) Start(ctx context.Context, grp *errgroup.Group) error {
	url := fmt.Sprintf(connectionString, d.user, d.host, d.port, d.database, d.cert, d.key, d.ca)

	conn, err := gorm.Open(postgres.Open(url),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: logger.New(
				NewGormLogger(d.logger.WithCtx(d.AppCtx.RootCtx)),
				logger.Config{
					LogLevel: logger.Error,
				},
			),
		})
	if err != nil {
		d.logger.Errorf("unable to connect to database: %s", err)
		return err
	}

	err = conn.Use(dbresolver.Register(
		dbresolver.Config{
			Sources: []gorm.Dialector{
				postgres.Open(url),
			},
		}).
		SetMaxOpenConns(d.maxConn).
		SetMaxIdleConns(d.idleConn).
		SetConnMaxLifetime(d.maxTime).
		SetConnMaxIdleTime(d.idleTime),
	)
	if err != nil {
		d.logger.Errorf("unable to configure database options: %s", err)
		return err
	}

	// don't REALLY need this...
	grp.Go(func() error {
		<-ctx.Done()
		return nil
	})

	d.conn = conn
	return nil
}

func (d *CockroachDB) Stop() error {
	d.logger.Infof("database connection closed")
	return nil
}
