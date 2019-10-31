package main

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/robfig/cron/v3"
	"os"
	"sync"
)

var wg sync.WaitGroup
var db *gorm.DB
var err error
var rows *sql.Rows

func main() {

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=test password=1234 sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	c := cron.New()
	c.AddFunc("@every 300ms", jobTask)
	c.Start()
	defer c.Stop()
	sig := make(chan os.Signal)
	<-sig

}

func jobTask() {
	rows, _ = db.Raw("select pid from pg_stat_activity where cardinality(pg_blocking_pids(pid)) > 0").Rows()
	defer rows.Close()

	for rows.Next() {
		var result int64
		rows.Scan(&result)
		db.Exec("SELECT pg_cancel_backend(?)", []int64{result})
		fmt.Println(result)
	}
}
