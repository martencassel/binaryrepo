package main

import (
	postgres "github.com/martencassel/binaryrepo/pkg/postgres"
	"github.com/spf13/viper"
)

type binaryrepoConfig struct {
	Database postgres.Config

}

func configFromArgv() *binaryrepoConfig {
	return &binaryrepoConfig{
		Database: parseDBOpts(),
	}
}

func parseDBOpts() postgres.Config {
	return postgres.Config{
		Host: viper.GetString("db-host"),
		Port: viper.GetInt("db-port"),
		DBName: viper.GetString("db-name"),
		User: viper.GetString("db-user"),
		Password: viper.GetString("db-password"),
		MaxOpenConnections: viper.GetInt32("db-max-open-connections"),
		MaxIdleConnections: viper.GetInt32("db-max-idle-connections"),
		ConnMaxLifetime: viper.GetDuration("db-conn-max-lifetime"),
		CreateDB: viper.GetBool("db-create-db"),
	}
}
