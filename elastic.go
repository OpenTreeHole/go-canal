package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var client = http.Client{Timeout: time.Second * 10}

func bulkInsert(table *Table) {
	if table.Buffer.Len() == 0 {
		log.Debugf("no data in buffer %s, continue...", table.Name)
		return
	}
	table.Lock.Lock()
	resp, err := client.Post(
		fmt.Sprintf("%s/%s/_bulk", config.ElasticUrl, table.Name),
		"application/x-ndjson",
		&table.Buffer,
	)
	if err == nil && resp.StatusCode == 200 {
		table.Buffer.Reset()
		log.Infof("push data of %s success", table.Name)
	}
	table.Lock.Unlock()

	if err != nil {
		log.Error(err)
		return
	}
	if resp.StatusCode != 200 {
		//goland:noinspection GoUnhandledErrorResult
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("Read body failed:", err)
			return
		}

		log.Errorf("elastic search response error %s", body)
	}
}

func StartTimer() {
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()
	for range ticker.C {
		for _, db := range config.Schemas {
			for _, table := range db.Tables {
				go bulkInsert(table)
			}
		}
	}
}
