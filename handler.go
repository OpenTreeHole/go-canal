package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
)

type MyEventHandler struct {
	canal.DummyEventHandler
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	config.Initialize(e.Table)

	log.Debugf("received a row event %v %v", e.Table, e.Action)
	db, ok := config.Schemas[e.Table.Schema]
	if !ok {
		return nil
	}
	table, ok := db.Tables[e.Table.Name]
	if !ok {
		return nil
	}

	for i, row := range e.Rows {
		if e.Action == canal.UpdateAction && i%2 == 0 {
			continue // ignore the even row because it's the row before update
		}
		receiveRow(table, row)
	}

	return nil
}

func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}

func receiveRow(table *Table, row []any) {
	data := map[string]any{}
	for _, name := range table.ColumnNames {
		value := row[table.Columns[name].Index]
		if v, ok := value.([]byte); ok {
			data[name] = string(v)
		} else {
			data[name] = value
		}
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Error("failed to parse row data to json", err)
	}
	log.Debug(string(bytes))
	table.Lock.Lock()
	table.Buffer.WriteString(fmt.Sprintf("{\"index\":{\"_id\": %d}}\n", data["id"]))
	table.Buffer.WriteString(string(bytes))
	table.Buffer.WriteString("\n")
	table.Lock.Unlock()
}
