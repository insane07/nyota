package store

import (
	"database/sql"
	"goprizm/sysutils"
	golog "log"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"nyota/backend/watch"
	"os"
	"time"

	gorp "gopkg.in/gorp.v2"

	_ "github.com/lib/pq"
)

// Store is a abstraction over persistent backend databases(postgres, cassandra etc)
type Store struct {
	db      *gorp.DbMap    // db gorp handler
	Watcher *watch.Watcher //redis pubsub
}

// TODO - currently database configs are obtained from environ vars. NewStore could be modified
// in future to accept database configs as args once there is better clarity on complete config.
func New() (*Store, error) {
	// cccdb, err := setupPg("postgres://appsuperuser:@localhost/cccdb")
	nyotadb, err := setupPg("postgres://postgres:postgres@localhost/nyota?sslmode=disable")
	if err != nil {
		return nil, err
	}
	store := &Store{
		db:      nyotadb,
		Watcher: watch.New(),
	}
	addNyotaTables(store.db)
	return store, nil
}

// DB returns pg handles for read/write ops to prizmdb.
func (store *Store) DB() SqlDB {
	return newSqlDB(store.db)
}

func setupPg(connStr string) (*gorp.DbMap, error) {
	maxIdleConns := sysutils.GetenvInt("DB_MAX_IDLE_CONNS", 10)
	maxOpenConns := sysutils.GetenvInt("DB_MAX_OPEN_CONNS", 100)
	maxConnLifeTime := sysutils.GetenvInt("DB_MAX_CONN_LIFE_TIME", 60) // minutes

	traceOn := sysutils.GetenvBool("DB_TRACE", false)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Duration(maxConnLifeTime) * time.Minute)
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}, TypeConverter: gorpTypeConverter{}}

	if traceOn {
		dbMap.TraceOn("dbmap:", golog.New(os.Stdout, "", golog.LstdFlags))
	}
	return dbMap, nil
}

func addNyotaTables(db *gorp.DbMap) {
	db.AddTableWithName(config.Event{}, "events").SetKeys(true, "id")
	db.AddTableWithName(model.UserTenantAttributes{}, "user_tenant_attributes")
	db.AddTableWithName(model.UserTenantDetails{}, "user_tenant_details")
	db.CreateTablesIfNotExists()
}

// SqlDB - manages a set of gorp handles to perform database read/write operations.
type SqlDB interface {
	SelectOne(holder interface{}, query string, args ...interface{}) error
	SelectInt(query string, args ...interface{}) (int64, error)
	Select(i interface{}, query string, args ...interface{}) error
	Insert(list ...interface{}) error
	Update(list ...interface{}) (int64, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	TruncateTables() error
	Begin() (*gorp.Transaction, error)
}

// sqlDB implements SqlDB interface.
type sqlDB struct {
	db *gorp.DbMap
}

func newSqlDB(db *gorp.DbMap) *sqlDB {
	return &sqlDB{
		db: db,
	}
}

func (sqlDB *sqlDB) SelectOne(holder interface{}, query string, args ...interface{}) error {
	err := sqlDB.db.SelectOne(holder, query, args...)
	return err
}

func (sqlDB *sqlDB) SelectInt(query string, args ...interface{}) (int64, error) {
	return sqlDB.db.SelectInt(query, args...)
}

func (sqlDB *sqlDB) Select(i interface{}, query string, args ...interface{}) error {
	// If `select *` is used it could return columns which does not have mapping for given obj 'i'
	// This will result in gorp to return NoFieldInTypeError along with actual rows.
	// Return nil error for these cases.
	_, err := sqlDB.db.Select(i, query, args...)
	if _, ok := err.(*gorp.NoFieldInTypeError); ok {
		return nil
	}
	return err
}

func (sqlDB *sqlDB) Insert(list ...interface{}) error {
	return sqlDB.db.Insert(list...)
}

func (sqlDB *sqlDB) Update(list ...interface{}) (int64, error) {
	return sqlDB.db.Update(list...)
}

func (sqlDB *sqlDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return sqlDB.db.Exec(query, args...)
}

func (sqlDB *sqlDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return sqlDB.db.Query(query, args...)
}

func (sqlDB *sqlDB) TruncateTables() error {
	return sqlDB.db.TruncateTables()
}

func (sqlDB *sqlDB) Begin() (*gorp.Transaction, error) {
	return sqlDB.db.Begin()
}
