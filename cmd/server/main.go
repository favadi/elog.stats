package main

import (
	"net"

	"go.uber.org/zap"

	"github.com/txchuyen/elog.stats/elog"
	"github.com/txchuyen/elog.stats/elog/store"
	"github.com/txchuyen/elog.stats/middleware/log"
	pb "github.com/txchuyen/elog.stats/pb/elog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// unmarshal any.Anty into proto.Message ...
	// TODO: why?
	_ "github.com/golang/protobuf/ptypes/any"
	_ "github.com/golang/protobuf/ptypes/duration"
	_ "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/golang/protobuf/ptypes/struct"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/golang/protobuf/ptypes/wrappers"
)

func main() {
	conf := loadConfig()
	db := connectDb(conf)
	defer closeDb(db)

	conn, err := net.Listen("tcp", conf.Service.Addr())
	if err != nil {
		panic("net: unable to listen: " + err.Error())
	}
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic("zap: unable to create logger " + err.Error())
	}
	var (
		opts = []grpc.ServerOption{
			grpc.UnaryInterceptor(log.ZapUnaryServerInterceptor(zapLogger)),
			grpc.StreamInterceptor(log.ZapStreamServerInterceptor(zapLogger)),
		}
		s     = grpc.NewServer(opts...)
		store = store.NewEvent(db)
		srv   = elog.NewServer(store)
	)
	pb.RegisterElogServer(s, srv)
	reflection.Register(s)

	if err := s.Serve(conn); err != nil {
		panic("service: unable to serve: " + err.Error())
	}
}
