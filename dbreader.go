package helper

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"fmt"
)

//数据库层结构，包含数据库连接和日志
type DBLayer struct {
	db *sql.DB
	logger *log.Logger
}

type dbRawRow sql.RawBytes

//结果集复制出来的大slice
var rawRows [][]string

//数据库连接初始化
func Init(user, password, host, port, database, charset string) *DBLayer {
	logFile, err := os.Create("db.log")
	if err != nil {
		panic(error(err))
	}
	logger := log.New(logFile, "[db]", log.LstdFlags|log.Llongfile)
	dnsname := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", user, password, host, port, database, charset)
	logger.Println(dnsname)
	db, err := sql.Open("mysql", dnsname)
	if err != nil {
		logger.Fatal(err)
		return nil
	}
	return &DBLayer{db, logger}
}

//从结果集中获取一行数据
func (con *DBLayer) FetchOne(sql string) []string {
	rows, err := con.db.Query(sql)
	if err != nil {
		con.logger.Fatal(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		con.logger.Fatal(err)
	}
	values := make([]dbRawRow, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	rows.Next()
	err = rows.Scan(scanArgs...)
	if err != nil {
		con.logger.Fatal(err)
	}
	rowString := make([]string, len(columns))
	for i, eachCol := range values {
		if eachCol != nil {
			rowString[i] = eachCol
		}else{
			rowString[i] = ""
		}
	}
	return rowString
}

//从结果集中获取所有数据
func (con *DBLayer) FetchAll(sql string) [][]string {
	rows, err := con.db.Query(sql)
	if err != nil {
		con.logger.Fatal(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		con.logger.Fatal(err)
	}
	values := make([]dbRawRow, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		rowString := make([]string, len(columns))
		err = rows.Scan(scanArgs...)
		if err != nil {
			con.logger.Fatal(err)
		}
		for i, eachCol := range values {
			if eachCol != nil {
				rowString[i] = string(eachCol)
			}else{
				rowString[i] = ""
			}
		}
		rawRows = append(rawRows, rowString)
	}
	return rawRows
}

//关闭数据库连接
func (con *DBLayer) Close() {
	con.db.Close()
}
