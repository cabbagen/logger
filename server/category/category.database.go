package category

import (
	"log"
	"fmt"
	"time"
	"strings"
	"database/sql"
	"logger/server/cache"
	_ "github.com/go-sql-driver/mysql"
)

const (
	database = "mysql"
	connection = "root:root@tcp(localhost:3306)/logger"
)

type DatabaseLogger struct {
	LoggerStruct
	
	Database      *sql.DB
}

func NewDatabaseLogger(appid, level string) DatabaseLogger {
	databaseLogger := DatabaseLogger{}
	
	databaseLogger.Appid = appid
	databaseLogger.Level = level
	databaseLogger.IsCache = true
	databaseLogger.LimitTime = 1 * time.Minute
	databaseLogger.Cacher = cache.NewRedisCacher()
	
	databaseLogger.Database = databaseLogger.Connect()
	
	return databaseLogger
}

func (dl DatabaseLogger) Connect() *sql.DB {
	db, error := sql.Open(database, connection)
	
	if error != nil {
		log.Fatal("数据库链接失败：", error)
	}
	
	return db
}

func (dl DatabaseLogger) Log(message string) {
	dl.AppendToCache(message)
	dl.ListenSyncTask()
}

func (dl DatabaseLogger) AppendToCache(message string) {
	if dl.Cacher == nil {
		log.Fatal("缓存服务连接异常")
	}
	
	now := time.Now()
	
	prefix := fmt.Sprintf("%d-%d-%d %d:%d:%d ", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	
	dl.Cacher.Push(dl.Appid, prefix + "@" + message)
}

func (dl DatabaseLogger) ListenSyncTask() {
	
	select {
	
		case <-time.After(dl.LimitTime):
			log.Println("state record database")
			dl.ExecCacheRecord()
	}
}

func (dl DatabaseLogger) ExecCacheRecord() {
	records, error := dl.Cacher.Flush(dl.Appid)
	
	if error != nil {
		log.Fatal("get cache records error: ", error)
	}
	
	stmt, error := dl.Database.Prepare("insert into loggers(logger_app, logger_lever, logger_message, logger_time) values(?, ?, ?, ?)")
	
	if error != nil {
		log.Fatal("database prepare error: ", error)
	}

	
	for _, record := range records {
		timestring, message := dl.GetLoggerMessageInfo(record)
		
		if _, error := stmt.Exec(dl.Appid, dl.Level, message, timestring); error != nil {
			log.Fatal("database insert into error: ", error)
		}
	}
}

func (dl DatabaseLogger) GetLoggerMessageInfo(record string) (string, string) {
	sslice := strings.Split(record, "@")
	
	return sslice[0], strings.Join(sslice[1:], "")
}



