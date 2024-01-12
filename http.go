package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

/*-----------------------*/
///2023-11-03
//功能记录
//1.打开开启8080端口并监听（包含/data方法）
//2.连接sqlserver数据库
//3.记录记事本日志
//4.返回json数据
//5.解决日志记录两条的问题
// Database connection string
const (
	server   = "."       //服务器地址
	port     = 1433      //SQLserver数据库端口号
	user     = "sa"      //用户名
	password = "123"     //密码
	database = "Test"    //数据库名
	encrypt  = "disable" //是否启用加密连接
)

// 日志文件名格式
const logFileFormat = "2006-01-02.log"

// Custom struct to represent the Data
type Data struct {
	Name string `json:"name"`
	Sex  string `json:"sex"`
}

var logFile *os.File

func main() {
	// 打开并创建当前日期的日志文件，文件前加./表示在当前目录下创建
	logFileName := "./log/" + time.Now().Format(logFileFormat)
	var err error
	logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	} else {
		logError("服务器启动成功，日志文件名为：%s", logFileName)
	}
	defer logFile.Close()

	// 初始化记录器以写入日志文件
	log.SetOutput(logFile)

	// 为“/data”URL路径注册数据处理程序
	http.HandleFunc("/data", dataHandler)

	// 启动HTTP服务器并侦听端口8080
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start the HTTP server: %v", err)
	}

}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// 连接数据库
	connection := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=%s", server, user, password, port, database, encrypt)
	db, err := sql.Open("mssql", connection)
	if err != nil {
		logError("Failed to connect to the database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 编写查询语句
	rows, err := db.Query("SELECT name, sex FROM C_test02")
	if err != nil {
		logError("Failed to execute the query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the results
	var results []Data

	// Process the query results
	// Process the query results
	for rows.Next() {
		var name, sex string
		if err := rows.Scan(&name, &sex); err != nil {
			logError("Failed to scan row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, Data{Name: name, Sex: sex})
	}

	if err := rows.Err(); err != nil {
		logError("Error in processing rows: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the results to JSON
	jsonData, err := json.Marshal(results)
	if err != nil {
		logError("Failed to marshal JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response content type and write the JSON data
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	// Log success message
	logError("获取数据成功: %v", results)
}

// 处理记录日志的方法被记录两次的问题，并改成一次

func logError(format string, args ...interface{}) {
	errorMsg := fmt.Sprintf(format, args...)
	//log.Printf(errorMsg)
	logFile.WriteString(time.Now().Format("2006-01-02 15:04:05") + " " + errorMsg + "\n")
}
