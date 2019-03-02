package server

import (
	"errors"
	"context"
	"logger/server/category"
	pb "logger/proto"
)

type server struct {}

func NewRPCServer() *server {
	return &server{}
}

func (s *server) AllLog(ctx context.Context, in *pb.LoggerRequest) (*pb.LoggerReply, error) {
	return generateLogRecord("ALL", in)
}

func (s *server) InfoLog(ctx context.Context, in *pb.LoggerRequest) (*pb.LoggerReply, error) {
	return generateLogRecord("INFO", in)
}

func (s *server) WarnLog(ctx context.Context, in *pb.LoggerRequest) (*pb.LoggerReply, error) {
	return generateLogRecord("WARN", in)
}

func (s *server) DebugLog(ctx context.Context, in *pb.LoggerRequest) (*pb.LoggerReply, error) {
	return generateLogRecord("DEBUG", in)
}

func (s *server) ErrorLog(ctx context.Context, in *pb.LoggerRequest) (*pb.LoggerReply, error) {
	return generateLogRecord("ERROR", in)
}

func generateLogRecord(level string, in *pb.LoggerRequest) (*pb.LoggerReply, error) {
	
	logger := createLogger(level, in).(category.Logger)
	
	if logger == nil {
		return nil, errors.New("创建 logger 对象失败")
	}
	
	logger.Log(*in.Message)
	
	code := "200"
	msg := "logger ok ..."
	
	return &pb.LoggerReply {
		Code: &code,
		Message: &msg,
	}, nil
}

func createLogger(level string, in *pb.LoggerRequest) interface{} {
	
	if *in.Category == "console" {
		return category.NewConsoleLogger(*in.Appid, level)
	}
	
	if *in.Category == "file" {
		return category.NewFileLogger(*in.Appid, level)
	}
	
	if *in.Category == "database" {
		return category.NewDatabaseLogger(*in.Appid, level)
	}
	
	return nil
}














