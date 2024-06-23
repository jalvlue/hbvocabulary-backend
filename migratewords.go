package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/tealeg/xlsx"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/db_course_design?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	excelFileName := "words.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		panic(err)
	}

	sheet := xlFile.Sheets[0] // 假设数据在第一个工作表中

	for _, row := range sheet.Rows {
		cell := row.Cells[0] // 假设数据在第一列中
		word := cell.String()

		// 执行插入数据库的操作
		_, err := db.Exec("INSERT INTO t_vocabularies (word) VALUES (?)", word)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("数据插入完成")
}
