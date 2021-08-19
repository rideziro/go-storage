package storage

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type Database struct {
	Driver       string
	Username     string
	Password     string
	Host         string
	ReplicaHost  string
	Port         string
	DatabaseName string

	DisableLog bool
}

func (d *Database) New() (*gorm.DB, error) {
	config := &gorm.Config{}
	if d.DisableLog {
		config.Logger = logger.Discard
	}
	db, err := gorm.Open(mysql.Open(d.DSN()), config)
	if err != nil {
		return nil, fmt.Errorf("database.New() init: %w", err)
	}

	if replica, err := d.DSNReplica(); err == nil {
		err := db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{mysql.Open(replica)},
		}))
		if err != nil {
			return nil, fmt.Errorf("database.New() create replica: %w", err)
		}
	}
	return db, nil
}

func (d *Database) DSN() string {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.DatabaseName,
	)
	return dsn
}

func (d *Database) DSNReplica() (string, error) {
	if d.ReplicaHost == "" {
		return "", errors.New("no replica found")
	}
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		d.Username,
		d.Password,
		d.ReplicaHost,
		d.Port,
		d.DatabaseName,
	)
	return dsn, nil
}

func (d *Database) DSNMigrate() string {
	dsn := fmt.Sprintf("%s://%s:%s@(%s:%s)/%s?charset=utf8mb4",
		d.Driver,
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.DatabaseName,
	)
	return dsn
}

func DefaultDatabase() Database {
	database := Database{
		Driver:       viper.GetString("DB_DRIVER"),
		Username:     viper.GetString("DB_USERNAME"),
		Password:     viper.GetString("DB_PASSWORD"),
		Host:         viper.GetString("DB_HOST"),
		ReplicaHost:  viper.GetString("DB_REPLICA_HOST"),
		Port:         viper.GetString("DB_PORT"),
		DatabaseName: viper.GetString("DB_NAME"),
	}

	if database.Driver == "" {
		database.Driver = "mysql"
	}

	return database
}
