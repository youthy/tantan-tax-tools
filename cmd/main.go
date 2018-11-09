package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/p1cn/tantan-tax-tools/app/config"
	"github.com/p1cn/tantan-tax-tools/app/db"
)

var configPath = flag.String("config", "./conf", "path to service config directory")
var nsrfile = flag.String("nsr", filepath.Join(*configPath, string(filepath.Separator), "nsrjcxx.csv"), "nsr file name")
var lwsdfile = flag.String("lwsd", filepath.Join(*configPath, string(filepath.Separator), "lwbcsd.csv"), "lwsd file name")
var delimeter = flag.String("dmt", "\t", "file delimeter")

func main() {
	flag.Parse()
	config.Init(*configPath)
	dbh, err := db.InitMysql()
	if err != nil {
		log.Fatalf("init db failed:%v", err)
	}
	defer dbh.Close()
	var data []interface{}
	nsrs, err := db.ParseNsrjcxxFile(*nsrfile, *delimeter)
	if err != nil {
		log.Fatalf("parse nsr file failed :%v", err)
	}
	lwsd, err := db.ParseLwbcsdFile(*lwsdfile, *delimeter)
	if err != nil {
		log.Fatalf("parse lwsd file failed :%v", err)
	}
	for _, n := range nsrs {
		data = append(data, n)
	}
	db.BackupTable(db.NsrTable, db.NsrTable+"_bak")
	db.TruncTable(db.NsrTable)
	err = db.InsertTables(db.NsrTable, data)
	if err != nil {
		log.Fatalf("insert nsrjcxx failed :%v", err)
	}
	data = data[:0]
	for _, n := range lwsd {
		data = append(data, n)
	}
	db.BackupTable(db.LwbcsdTable, db.LwbcsdTable+"_bak")
	db.TruncTable(db.LwbcsdTable)
	err = db.InsertTables(db.LwbcsdTable, data)
	if err != nil {
		log.Fatalf("insert lwbcsd failed :%v", err)
	}

	fmt.Println("end")
}
