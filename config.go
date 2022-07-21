package main

import (
	"bytes"
	"flag"
	"github.com/go-mysql-org/go-mysql/schema"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"sync"
)

type Column struct {
	Index int // index in the row
}

type Table struct {
	Name        string             `yaml:"name"`
	Columns     map[string]*Column `yaml:"-"`
	ColumnNames []string           `yaml:"columns"`
	Initialized bool
	Buffer      bytes.Buffer
	Lock        sync.Mutex
}

type Schema struct {
	Name        string            `yaml:"name"`
	Tables      map[string]*Table `yaml:"tables"`
	Initialized bool
}

func (db *Schema) Initialize(eventTable *schema.Table) bool {
	if db.Initialized {
		return true
	}

	table, ok := db.Tables[eventTable.Name]
	if !ok {
		return false
	}

	for k := 0; k < len(eventTable.Columns); k++ {
		column, ok := table.Columns[eventTable.Columns[k].Name]
		if !ok {
			continue
		}
		column.Index = k
	}
	table.Initialized = true

	initialized := true
	for _, table = range db.Tables {
		initialized = initialized && table.Initialized
	}
	db.Initialized = initialized
	return initialized
}

type Config struct {
	Address     string             `yaml:"address"`
	User        string             `yaml:"user"`
	Password    string             `yaml:"password"`
	ElasticUrl  string             `yaml:"elastic_url"`
	Schemas     map[string]*Schema `yaml:"schemas"`
	Initialized bool               `yaml:"-"`
}

func (config *Config) Initialize(eventTable *schema.Table) {
	if config.Initialized {
		return
	}
	log.Info("initializing config")
	db, ok := config.Schemas[eventTable.Schema]
	if !ok {
		return
	}

	if db.Initialize(eventTable) {
		initialized := true
		for _, db = range config.Schemas {
			initialized = initialized && db.Initialized
		}
		config.Initialized = initialized
		if initialized {
			log.Info("config initialized")
		}
	}
}

var config Config

func init() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.yaml", "config file path")
	flag.Parse()

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	for name, db := range config.Schemas {
		if db.Name == "" {
			db.Name = name
		}
		for tableName, table := range db.Tables {
			if table.Name == "" {
				table.Name = tableName
			}
			table.Columns = make(map[string]*Column)
			for _, columnName := range table.ColumnNames {
				table.Columns[columnName] = &Column{}
			}
			if !slices.Contains(table.ColumnNames, "id") {
				panic("table " + tableName + " must have an id column")
			}
		}
	}
}
