package db

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/p1cn/tantan-tax-tools/app/config"
)

var mysqlDB *sql.DB

// Table ...
type Table interface {
}

// MysqlDB return global mysql db
func MysqlDB() *sql.DB {
	return mysqlDB
}

// InitMysql ...
func InitMysql() (*sql.DB, error) {
	dbConfig := config.Get().Database
	var err error
	if dbConfig == nil {
		return nil, errors.New("no db config")
	}
	str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.UserName, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	fmt.Printf("connect cmd  %s\n", str)
	mysqlDB, err = sql.Open("mysql", str)
	if err != nil {
		err = mysqlDB.Ping()
	}
	return mysqlDB, err
}

// Insert insert values
func Insert(sql string, args ...interface{}) (sql.Result, error) {
	stmt, err := mysqlDB.Prepare(sql)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		log.Println(err)
	}
	return res, err
}

// Fields ...
func Fields(t interface{}) []string {
	var fields []string
	val := reflect.ValueOf(t).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		field := typeField.Name
		if t, b := tag.Lookup("mysql"); b {
			field = t
		}
		fields = append(fields, field)
	}
	return fields
}

// InsertTables ...
func InsertTables(tableName string, rows []interface{}) error {
	if len(rows) < 1 {
		return nil
	}
	var values []interface{}
	fields := Fields(rows[0])
	sql := fmt.Sprintf("Insert into %s (%s) values ", tableName, strings.Join(fields, ","))
	bracket := genBracket(len(fields), len(rows))
	sql += bracket
	for _, t := range rows {
		val := reflect.ValueOf(t).Elem()
		for i := 0; i < val.NumField(); i++ {
			valueField := val.Field(i)
			values = append(values, valueField.Interface())
		}
	}
	_, err := Insert(sql, values...)
	return err
}

func genBracket(fieldLen int, dataLen int) string {
	tmp := strings.Repeat("?,", fieldLen)
	tmp = tmp[0 : len(tmp)-1]
	tmp = "(" + tmp + "),"
	tmp = strings.Repeat(tmp, dataLen)
	return tmp[0 : len(tmp)-1]
}

// BackupTable back up table to new table
func BackupTable(tableName string, newName string) error {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s LIKE %s", newName, tableName)
	_, err := mysqlDB.Query(sql)
	if err != nil {
		return err
	}
	err = TruncTable(newName)
	if err != nil {
		return err
	}
	sql = fmt.Sprintf("INSERT INTO %s SELECT * FROM %s", newName, tableName)
	_, err = mysqlDB.Query(sql)
	if err != nil {
		return err
	}
	return nil
}

// TruncTable a table
func TruncTable(tableName string) error {
	sql := fmt.Sprintf("TRUNCATE TABLE %s", tableName)
	_, err := mysqlDB.Query(sql)
	if err != nil {
		return err
	}
	return nil
}

// ParseFile parse text from file to userdomain PrivateQuestion
func ParseFile(fileName string, delimiter string) ([][]string, error) {

	rw, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer rw.Close()

	rb := bufio.NewReader(rw)
	i := 1
	var ret [][]string
	for {
		lineStr, err := readline(rb)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if lineStr == "" {
			continue
		}
		if i == 1 && strings.HasPrefix(lineStr, "\ufeff") {
			//remove BOM tag
			lineStr = lineStr[3:]
		}
		// line format
		values := strings.Split(lineStr, delimiter)
		ret = append(ret, values)
		i++

	}
	fmt.Printf("file: %s  lines :%d\n", fileName, i-1)
	return ret, nil
}

func readline(rb *bufio.Reader) (string, error) {
	line, _, err := rb.ReadLine()
	if err != nil {
		return "", err
	}
	lineStr := string(line)
	lineStr = strings.TrimRight(lineStr, "\n\r")
	return lineStr, nil
}
