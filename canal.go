package main

import (
	"github.com/go-mysql-org/go-mysql/canal"
)

func NewConfig() *canal.Config {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = config.Address
	cfg.User = config.User
	cfg.Password = config.Password

	return cfg
}

func Dump() {
	for _, db := range config.Schemas {
		cfg := NewConfig()
		cfg.Dump.TableDB = db.Name
		for tableName := range db.Tables {
			cfg.Dump.Tables = append(cfg.Dump.Tables, tableName)
		}
		c, err := canal.NewCanal(NewConfig())
		if err != nil {
			panic(err)
		}
		err = c.Run()
		if err != nil {
			panic(err)
		}
	}
}
