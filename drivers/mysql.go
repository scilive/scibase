package drivers

import (
	"time"

	"github.com/scilive/scibase/logs"

	"github.com/daqiancode/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateMysql() (*gorm.DB, error) {
	var datetimePrecision = 2
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       env.Get("MYSQL_URL", "root:123456@tcp(localhost:3306)/myiam?charset=utf8&parseTime=True&loc=Local"), // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,                                                                                                 // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,                                                                                                // disable datetime precision support, which not supported before MySQL 5.6
		DefaultDatetimePrecision:  &datetimePrecision,                                                                                  // default datetime precision
		DontSupportRenameIndex:    true,                                                                                                // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                                // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                               // smart configure based on used version
	}), &gorm.Config{})
	if err != nil {
		logs.Log.Error().Err(err).Msg("failed create mysql connection")
		return nil, err
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(env.GetIntMust("MYSQL_MAX_IDLE", 5))
	sqlDB.SetMaxOpenConns(env.GetIntMust("MYSQL_MAX_OPEN", 100))
	sqlDB.SetConnMaxLifetime(time.Duration(env.GetIntMust("MYSQL_MAX_LIFE_TIME", 30)) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(env.GetIntMust("MYSQL_MAX_IDLE_TIME", 5)) * time.Minute)
	return db, err
}

var db *gorm.DB

func GetDB() *gorm.DB {
	var err error
	if db == nil {
		db, err = CreateMysql()
		if err != nil {
			panic(err)
		}
	}
	return db
}

type NewMysqlOptions struct {
	MysqlUrl          string
	MaxIdleConns      int
	MaxOpenConns      int
	MaxLifeTime       time.Duration
	MaxIdleTime       time.Duration
	ConnMaxIdleTime   time.Duration
	DatetimePrecision int
	DefaultStringSize int
}

func defaultNewMysqlOptions(opts NewMysqlOptions) NewMysqlOptions {
	if opts.MaxIdleConns == 0 {
		opts.MaxIdleConns = 5
	}
	if opts.MaxOpenConns == 0 {
		opts.MaxOpenConns = 100
	}
	if opts.MaxLifeTime == 0 {
		opts.MaxLifeTime = 30 * time.Minute
	}
	if opts.MaxIdleTime == 0 {
		opts.MaxIdleTime = 5 * time.Minute
	}
	if opts.ConnMaxIdleTime == 0 {
		opts.ConnMaxIdleTime = 5 * time.Minute
	}
	if opts.DatetimePrecision == 0 {
		opts.DatetimePrecision = 2
	}
	if opts.DefaultStringSize == 0 {
		opts.DefaultStringSize = 256
	}
	return opts
}

func NewMysql(opts NewMysqlOptions) (*gorm.DB, error) {
	opts = defaultNewMysqlOptions(opts)
	var datetimePrecision = opts.DatetimePrecision
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       opts.MysqlUrl,                // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         uint(opts.DefaultStringSize), // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,                         // disable datetime precision support, which not supported before MySQL 5.6
		DefaultDatetimePrecision:  &datetimePrecision,           // default datetime precision
		DontSupportRenameIndex:    true,                         // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                         // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                        // smart configure based on used version
	}), &gorm.Config{})
	if err != nil {
		logs.Log.Error().Err(err).Msg("failed create mysql connection")
		return nil, err
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(opts.MaxIdleConns)
	sqlDB.SetMaxOpenConns(opts.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(opts.MaxLifeTime)
	sqlDB.SetConnMaxIdleTime(opts.ConnMaxIdleTime)
	return db, err
}
