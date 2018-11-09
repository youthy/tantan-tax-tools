package main

import (
	"fmt"
	//"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/p1cn/tantan-tax-tools/app/config"
	"github.com/p1cn/tantan-tax-tools/app/db"
)

func TestMysqlInsert(t *testing.T) {
	config.Init("../conf")
	_, err := db.InitMysql()
	if err != nil {
		t.Fatalf("db init failed")
	}
	defer db.MysqlDB().Close()
	data := []interface{}{}
	a := &db.Nsrjcxx{
		Xm:       "test1",
		Zzlx:     "身份证",
		Zzhm:     "110226199204282611",
		Gj:       "中国",
		Xb:       "男",
		Nsrzt:    "正常",
		Sfgd:     "否",
		Sfgy:     "否",
		Lxdh:     "15011536289",
		Sftdhy:   "否",
		Sftstzr:  "否",
		Bmbh:     "001",
		Qynsrsbh: "91110105327178273U",
	}

	b := &db.Lwbcsd{
		XM:       "test1",
		ZZLX:     "身份证",
		ZZHM:     "110226199204282611",
		SRE:      "1000",
		NSRSBH:   "001",
		QYNSRSBH: "91110105327178273U",
	}

	data = append(data, a)
	err = db.InsertTables("nsrjcxx", data)
	if err != nil {
		t.Fatalf("insert failed:%v", err)
	}
	data = append(data[:0], b)
	err = db.InsertTables("lwbcsd", data)
	if err != nil {
		t.Fatalf("insert failed:%v", err)
	}

	//err = db.BackupTable("nsrjcxx", "haha")
	//if err != nil {
	//	t.Fatal(err)
	//}

	//stmt, err := db.MysqlDB().Prepare("Insert into nsrjcxx (xm,zzlx,zzhm,gj,xb,nsrzt,sfgd,sfgy,lxdh,sftdhy,SFTSTZR,BMBH,QYNSRSBH) values (?,?,?,?,?,?,?,?,?,?,?,?,?)")
	//defer stmt.Close()

	//if err != nil {
	//	t.Fatalf("err :%v", err)
	//	//return nil, err
	//}
	//args := []interface{}{}
	//for i := 0; i < 13; i++ {
	//	args = append(args, "haha")
	//}
	////fmt.Printf("stmt:%v\n", stmt)
	//_, err = stmt.Exec(args...)
	//if err != nil {
	//	t.Fatalf("err :%v", err)
	//}
}

func TestReadFile(t *testing.T) {
	nsrs, err := db.ParseNsrjcxxFile("./tmp.csv", "\t")
	if err != nil {
		t.Fatalf("parse nsrjcxx failed: %v", err)
	}
	t.Logf("nsrs %v", nsrs)
	fmt.Printf("nsr %+v", nsrs)
}
