package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:root1234@tcp(127.0.0.1:3306)/wiser?charset=utf8mb4")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("failed to connect to mysql, err: ", err.Error())
		os.Exit(1) // 退出程序
	}

	// 初始化数据库

	//exitIfError(CreateTableWithSettings())
	//exitIfError(CreateTableWithDocuments())
	//exitIfError(CreateTableWithTokens())
	//exitIfError(CreateUniqueIndexBetweenTokenIndexAndTokens())
	//exitIfError(CreateUniqueIndexBetweenTitleIndexAndDocuments())
}

// ModifyDB 操作数据库
func ModifyDB(sql string, args ...interface{}) (int64, error) {
	result, err := db.Exec(sql, args...)
	if err != nil {
		fmt.Println("failed to modify db, err: ", err.Error())
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		fmt.Println("failed to modify db when rows affected, err: ", err.Error())
		return 0, nil
	}
	return count, nil
}

func CreateTableWithSettings() (err error) {
	sqlStr := `CREATE TABLE IF NOT EXISTS settings (
				  id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
                  key   TEXT,
                  value TEXT
               )`
	_, err = ModifyDB(sqlStr)
	return
}

func CreateTableWithDocuments() (err error) {
	sqlStr := `CREATE TABLE IF NOT EXISTS documents (
				  id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
				  title   TEXT NOT NULL,
                  body    TEXT NOT NULL
				)`
	_, err = ModifyDB(sqlStr)
	return
}

func CreateTableWithTokens() (err error) {
	sqlStr := `CREATE TABLE IF NOT EXISTS tokens (
				  id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
                  token      TEXT NOT NULL,
                  docs_count INT NOT NULL,
                  postings   BLOB NOT NULL
               )`
	_, err = ModifyDB(sqlStr)
	return
}

func CreateUniqueIndexBetweenTokenIndexAndTokens() (err error) {
	sqlStr := "CREATE UNIQUE INDEX token_index ON tokens(token);"
	_, err = ModifyDB(sqlStr)
	return
}

func CreateUniqueIndexBetweenTitleIndexAndDocuments() (err error) {
	sqlStr := "CREATE UNIQUE INDEX title_index ON documents(title);"
	_, err = ModifyDB(sqlStr)
	return
}

// exitIfError 发生错误就停止运行
func exitIfError(err error) {
	if err != nil {
		fmt.Printf("fatal error: %v\n", err)
		os.Exit(1)
	}
}
