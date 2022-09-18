package main

import (
	"pg_permissions/cmd"
	"pg_permissions/errors"
	"pg_permissions/log"

	"github.com/alecthomas/kong"
	"github.com/spf13/viper"
)

var (
	v *viper.Viper
)

var cli struct {
	ConfigPath string `name:"config" short:"c" default:"." help:"Directory to search for config.yaml file" type:"path"`

	Postgres cmd.PostgresClientCmd `cmd:"" postgres:"gather permission details"`
}

func initConfig(path string) *viper.Viper {
	v = viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(path)
	v.AutomaticEnv()
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalln(errors.ErrConfigFileRead)
	}
	return v
}

func main() {
	log.Infoln("Starting main routine")
	ctx := kong.Parse(&cli, kong.UsageOnError())

	v = initConfig(cli.ConfigPath)
	ctx.Bind(v)

	err := ctx.Run()
	if err != nil {
		log.Infoln(err)
		log.Fatalf(errors.ErrContextRun)
	}
}
