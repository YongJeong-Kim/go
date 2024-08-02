package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBSource   string `mapstructure:"DB_SOURCE"`
	ShardCount int    `mapstructure:"SHARD_COUNT"`
}

func LoadDBConfig(path string) *Config {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("load config file failed.", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("config unmarshal failed.", err)
	}
	return &cfg
}

func ConnDBs(cfg *Config, shardLen int) []*sqlx.DB {
	dbs := make([]*sqlx.DB, shardLen)
	for i := range shardLen {
		dbSource := fmt.Sprintf(cfg.DBSource, i)
		log.Println(dbSource)
		db, err := sqlx.Connect(cfg.DBDriver, dbSource)
		if err != nil {
			log.Fatal("db connect failed. ", err)
		}
		if err := db.Ping(); err != nil {
			log.Fatal("db ping failed. ", err)
		}
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(10)
		db.SetConnMaxLifetime(3 * time.Minute)
		dbs[i] = db
	}

	return dbs
}
