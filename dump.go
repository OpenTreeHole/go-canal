package main

import (
	"encoding/json"
	"fmt"
	mysqlClient "github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	"strings"
)

func getClient() *mysqlClient.Conn {
	connection, err := mysqlClient.Connect(config.Address, config.User, config.Password, "")
	if err != nil {
		panic(err)
	}

	err = connection.Ping()
	if err != nil {
		panic(err)
	}

	return connection
}

func Dump() {
	connection := getClient()

	log.Info("-dump is true, start dumping...")
	for _, db := range config.Schemas {
		for _, table := range db.Tables {
			log.Infof("start dumping for %s.%s", db.Name, table.Name)

			n := 0

			var result mysql.Result
			err := connection.ExecuteSelectStreaming(
				fmt.Sprintf(
					"SELECT %s FROM `%s`.`%s`",
					strings.Join(table.ColumnNames, ","),
					db.Name,
					table.Name,
				),
				&result,
				func(row []mysql.FieldValue) error {
					data := map[string]any{}
					for i, field := range row {

						if field.Type == mysql.FieldValueTypeString {
							data[table.ColumnNames[i]] = string(field.AsString())
						} else {
							data[table.ColumnNames[i]] = field.Value()
						}
					}
					bytes, err := json.Marshal(data)
					if err != nil {
						log.Error("failed to parse row data to json", err)
					}
					table.Buffer.WriteString(fmt.Sprintf("{\"index\":{\"_id\": %d}}\n", data["id"]))
					table.Buffer.WriteString(string(bytes))
					table.Buffer.WriteString("\n")
					n++
					if n == 100000 {
						BulkInsert(table)
						n = 0
					}
					return nil
				}, nil)
			if err != nil {
				panic(err)
			}

			BulkInsert(table)

			log.Infof("dump %s finished, tables: %s", db.Name, table.Name)
		}
	}
}
