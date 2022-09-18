package types

import (
	pe "pg_permissions/errors"
	"pg_permissions/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"fmt"
)

type Client struct {
	Database *sqlx.DB
	Config   *Config
}

func (c *Client) initDatabase() {
	log.Infoln("Initializing database object")
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=require",
		c.Config.DatabaseUser,
		c.Config.DatabasePassword,
		c.Config.DatabaseHost,
		c.Config.DatabasePort,
		c.Config.DatabaseName,
	)
	connection, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(pe.GenericError{Code: pe.ErrDBConnect})
	}
	c.Database = connection
}

func (c *Client) initLogger() {
	log.Infoln("Initializing log object")
	log.SetLogLevel(&c.Config.LoggingLevel)
}

type Config struct {
	LoggingLevel        string
	DatabaseHost        string
	DatabasePort        string
	DatabaseName        string
	DatabaseUser        string
	DatabasePassword    string
	DatabasePermissions string
	MaxRequests         int
	ReportHTML          bool
	ReportJSON          bool
}

func setDefault(v *viper.Viper) {
	v.SetDefault("logging.level", "info")
	v.SetDefault("database.host", "127.0.0.1")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.name", "test")
	v.SetDefault("database.user", "test")
	v.SetDefault("database.permissions", "ALL")
	v.SetDefault("performance.maxrequests", 10)
	v.SetDefault("report.html", true)
	v.SetDefault("report.json", true)
}

func readConfig(v *viper.Viper) *Config {
	config := &Config{
		LoggingLevel:        v.GetString("logging.level"),
		DatabaseHost:        v.GetString("database.host"),
		DatabasePort:        v.GetString("database.port"),
		DatabaseName:        v.GetString("database.name"),
		DatabaseUser:        v.GetString("database.user"),
		DatabasePassword:    v.GetString("database.password"),
		DatabasePermissions: v.GetString("database.permissions"),
		MaxRequests:         v.GetInt("performance.maxrequests"),
		ReportHTML:          v.GetBool("report.html"),
		ReportJSON:          v.GetBool("report.json"),
	}
	return config
}

func NewClient(v *viper.Viper) *Client {
	c := &Client{}
	setDefault(v)
	c.Config = readConfig(v)
	c.initLogger()
	c.initDatabase()
	return c
}
