package category

import (
	"os"
	"fmt"
	"log"
	"time"
	"strings"
	"io/ioutil"
	"logger/server/cache"
)

type FileLogger struct {
	LoggerStruct
	
	CurrentFileName   string
}

const FilePrefix = "logfile-"

func NewFileLogger(appid, level string) FileLogger {
	fileLogger := FileLogger{}
	
	fileLogger.Appid = appid
	fileLogger.Level = level
	fileLogger.IsCache = true
	fileLogger.LimitTime = 1 * time.Minute
	fileLogger.Cacher = cache.NewRedisCacher()
	fileLogger.CurrentFileName = ""
	
	return fileLogger
}

func (fl FileLogger) Log(message string) {	
	fl.AppendToCache(message)
	fl.ListenSyncTask()
}

func (fl FileLogger) AppendToCache(message string) {
	if fl.Cacher == nil {
		log.Fatal("缓存服务连接异常")
	}
	
	fl.Cacher.Push(FilePrefix + fl.Appid, fl.adapteMessageFormat(message))
}

func (fl FileLogger) adapteMessageFormat(message string) string {
	now := time.Now()
	
	prefix := fmt.Sprintf("[%d-%d-%d %d-%d-%d] ", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	
	return fmt.Sprintf("%s app[%s] [%s] %s", prefix, fl.Appid, fl.Level, message)
}

func (fl FileLogger) ListenSyncTask() {
	fl.CurrentFileName = fl.createTodayLoggerFiles()
	
	select {
	
		case <-time.After(fl.LimitTime):
			log.Println("state record file")
			fl.ExecCacheRecord()
	}
}

func (fl FileLogger) createTodayLoggerFiles() string {
	now := time.Now()
	
	dirname := fmt.Sprintf("./files/%s/%d-%d-%d/", fl.Appid, now.Year(), now.Month(), now.Day())
	
	isExistFiles := fl.isExistTodayLoggerFiles(dirname)
	
	if isExistFiles {
		return dirname + strings.ToLower(fl.Level) + ".logger.log"
	}
	
	if error := os.MkdirAll(dirname, os.ModePerm); error != nil {
		log.Fatal("mkdir error: ", error)
	}
	
	if error := ioutil.WriteFile(dirname + strings.ToLower(fl.Level) + ".logger.log", []byte(""), os.ModePerm); error != nil {
		log.Fatal("mkfile error: ", error)
	}
	
	return dirname + strings.ToLower(fl.Level) + ".logger.log"
}

func (fl FileLogger) isExistTodayLoggerFiles(filepath string) bool {
	_, error := os.Stat(filepath)
	
	if error != nil {
		if os.IsExist(error) {
			return true
		}
		return false
	}
	
	return true
}


func (fl FileLogger) ExecCacheRecord() {
	file, error := os.OpenFile(fl.CurrentFileName, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	
	if error != nil {
		log.Fatal("open file error: ", error)
	}
	
	records, error := fl.Cacher.Flush(FilePrefix + fl.Appid)
	
	if error != nil {
		log.Fatal("get cache records error: ", error)
	}
	
	messages := ""
	
	for _, record := range records {
		messages += (record + "\n")
	}
	
	if _, error := file.Write([]byte(messages)); error != nil {
		log.Fatal("record file error: ", error)
	}
	
	if error := file.Close(); error != nil {
		log.Fatal(error)
	}
}
