package xbgorm

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/starryck/strk-tc-x-lib-go/source/core/base/xbcfg"
	"github.com/starryck/strk-tc-x-lib-go/source/core/utility/xblogger"
)

var mPostgresClient *PostgresClient

type PostgresClient = Client

func GetPostgresClient() *PostgresClient {
	if mPostgresClient == nil {
		mPostgresClient = NewPostgresClient(nil)
	}
	return mPostgresClient
}

func NewPostgresClient(options *PostgresClientOptions) *PostgresClient {
	client := (&postgresClientBuilder{options: options}).
		initialize().
		setDSN().
		setHost().
		setPort().
		setName().
		setUser().
		setPassword().
		setClient().
		initClient().
		build()
	return client
}

type postgresClientBuilder struct {
	client  *PostgresClient
	configs *postgresClientConfigs
	options *PostgresClientOptions
}

type PostgresClientOptions struct {
	DSN      *string
	Host     *string
	Port     *string
	Name     *string
	User     *string
	Password *string
}

func (builder *postgresClientBuilder) build() *PostgresClient {
	return builder.client
}

func (builder *postgresClientBuilder) initialize() *postgresClientBuilder {
	builder.configs = &postgresClientConfigs{}
	if builder.options == nil {
		builder.options = &PostgresClientOptions{}
	}
	return builder
}

func (builder *postgresClientBuilder) setDSN() *postgresClientBuilder {
	dsn := builder.options.DSN
	if dsn != nil {
		builder.configs.dsn = *dsn
	}
	return builder
}

func (builder *postgresClientBuilder) setHost() *postgresClientBuilder {
	host := builder.options.Host
	if host != nil {
		builder.configs.host = *host
	} else {
		builder.configs.host = xbcfg.GetPostgresHost()
	}
	return builder
}

func (builder *postgresClientBuilder) setPort() *postgresClientBuilder {
	port := builder.options.Port
	if port != nil {
		builder.configs.port = *port
	} else {
		builder.configs.port = xbcfg.GetPostgresPort()
	}
	return builder
}

func (builder *postgresClientBuilder) setName() *postgresClientBuilder {
	name := builder.options.Name
	if name != nil {
		builder.configs.name = *name
	} else {
		builder.configs.name = xbcfg.GetPostgresName()
	}
	return builder
}

func (builder *postgresClientBuilder) setUser() *postgresClientBuilder {
	user := builder.options.User
	if user != nil {
		builder.configs.user = *user
	} else {
		builder.configs.user = xbcfg.GetPostgresUser()
	}
	return builder
}

func (builder *postgresClientBuilder) setPassword() *postgresClientBuilder {
	password := builder.options.Password
	if password != nil {
		builder.configs.password = *password
	} else {
		builder.configs.password = xbcfg.GetPostgresPassword()
	}
	return builder
}

func (builder *postgresClientBuilder) setClient() *postgresClientBuilder {
	config := builder.configs.getConfig()
	dialector := builder.configs.getDialector()
	if client, err := gorm.Open(dialector, config); err != nil {
		panic(err)
	} else {
		builder.client = client
	}
	return builder
}

func (builder *postgresClientBuilder) initClient() *postgresClientBuilder {
	client := builder.client
	var db *sql.DB
	if mDB, err := client.DB(); err != nil {
		panic(err)
	} else {
		db = mDB
	}
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10 * time.Minute)
	if xblogger.IsDebugLevel() {
		client = client.Debug()
	}
	builder.client = client
	return builder
}

type postgresClientConfigs struct {
	dsn      string
	host     string
	port     string
	name     string
	user     string
	password string
}

func (configs *postgresClientConfigs) getConfig() *Config {
	config := &Config{
		Logger: logger.New(xblogger.GetLogger(), logger.Config{
			LogLevel:      logger.Warn,
			SlowThreshold: 200 * time.Millisecond,
		}),
		NamingStrategy: &schema.NamingStrategy{
			SingularTable: true,
		},
	}
	return config
}

func (configs *postgresClientConfigs) getDialector() Dialector {
	dialector := postgres.Open(configs.makeDSN())
	return dialector
}

func (configs *postgresClientConfigs) makeDSN() string {
	dsn := configs.dsn
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
			configs.host, configs.port, configs.name, configs.user, configs.password)
	}
	return dsn
}
