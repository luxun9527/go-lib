package busi

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dtm-labs/client/dtmgrpc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lithammer/shortuuid/v3"
	"go-lib/sdk/dtm/saga/grpc_saga_barrier/busi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"testing"
)

// 定义一个全局对象db
var db *sql.DB

// 定义一个初始化数据库的函数
func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:root@tcp(192.168.2.99:3306)/dtm_busi?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

var BusiGrpc = fmt.Sprintf("192.168.2.138:%d", 58081)

func TestGrpcSagaBarrier(t *testing.T) {
	req := &busi.ReqGrpc{Amount: 30}
	gid := shortuuid.New()
	saga := dtmgrpc.NewSagaGrpc("192.168.2.99:36790", gid).
		Add(BusiGrpc+"/busi.Busi/TransOutBSaga", BusiGrpc+"/busi.Busi/TransOutRevertBSaga", req).
		Add(BusiGrpc+"/busi.Busi/TransInBSaga", BusiGrpc+"/busi.Busi/TransInRevertBSaga", req).
		Add(BusiGrpc+"/busi.Busi/PayCommissions", BusiGrpc+"/busi.Busi/PayCommissionsRevert", req)
	saga.WaitResult = true
	err := saga.Submit()
	if err != nil {
		log.Println(err)
	}

	log.Println(saga.WaitResult)

}

func TestNewGrpcServer(t *testing.T) {
	listener, err := net.Listen("tcp", "0.0.0.0:58081")
	if err != nil {
		log.Println("net listen err ", err)
		return
	}

	s := grpc.NewServer()
	busi.RegisterBusiServer(s, new(GrpcSagaServer))
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}
}

// GrpcSagaServer must be embedded to have forward compatible implementations.
type GrpcSagaServer struct {
	busi.UnimplementedBusiServer
}

func (GrpcSagaServer) TransIn(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransIn not implemented")
}
func (GrpcSagaServer) TransOut(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransOut not implemented")
}
func (GrpcSagaServer) TransInRevert(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransInRevert not implemented")
}
func (GrpcSagaServer) TransOutRevert(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransOutRevert not implemented")
}
func (GrpcSagaServer) TransInConfirm(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransInConfirm not implemented")
}
func (GrpcSagaServer) TransOutConfirm(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransOutConfirm not implemented")
}
func (GrpcSagaServer) XaNotify(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method XaNotify not implemented")
}
func (GrpcSagaServer) TransInXa(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransInXa not implemented")
}
func (GrpcSagaServer) TransOutXa(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransOutXa not implemented")
}
func (GrpcSagaServer) TransInTcc(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransInTcc not implemented")
}
func (GrpcSagaServer) TransOutTcc(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransOutTcc not implemented")
}
func (GrpcSagaServer) TransInTccNested(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransInTccNested not implemented")
}
func (GrpcSagaServer) TransInBSaga(ctx context.Context, req *busi.ReqGrpc) (*emptypb.Empty, error) {
	log.Println("TransInBSaga")
	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		return nil, err
	}
	if err := initDB(); err != nil {
		return nil, err
	}
	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		_, err = db.Exec("update user_account set balance= balance+? where id =2", 10)
		return err
	}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (GrpcSagaServer) TransOutBSaga(ctx context.Context, req *busi.ReqGrpc) (*emptypb.Empty, error) {
	log.Println("TransOutBSaga")
	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		return nil, err
	}
	if err := initDB(); err != nil {
		return nil, err
	}
	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		_, err = db.Exec("update user_account set balance= balance-? where id =1", 10)
		return err
	}); err != nil {
		//	return nil, status.Error(codes.Aborted, "test")
		return nil, status.Error(codes.Aborted, "test232323")
	}
	return &emptypb.Empty{}, nil
}
func (GrpcSagaServer) TransInRevertBSaga(ctx context.Context, req *busi.ReqGrpc) (*emptypb.Empty, error) {
	log.Println("TransInRevertBSaga")
	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		return nil, err
	}
	if err := initDB(); err != nil {
		return nil, err
	}
	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		_, err = db.Exec("update user_account set balance= balance+? where id =1", 10)
		return err
	}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (GrpcSagaServer) TransOutRevertBSaga(ctx context.Context, req *busi.ReqGrpc) (*emptypb.Empty, error) {
	log.Println("TransOutRevertBSaga")
	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		return nil, err
	}
	if err := initDB(); err != nil {
		return nil, err
	}
	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		_, err = db.Exec("update user_account set balance= balance-? where id =1", 10)
		return err
	}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (GrpcSagaServer) PayCommissions(ctx context.Context, req *busi.ReqGrpc) (*empty.Empty, error) {
	log.Println("PayCommissions")
	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		return nil, err
	}
	if err := initDB(); err != nil {
		return nil, err
	}
	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		_, err = db.Exec("update user_account1 set balance= balance+? where id =3", 10)
		return err
	}); err != nil {
		return nil, status.Error(codes.Aborted, "test11212")
	}
	return &emptypb.Empty{}, nil
}
func (GrpcSagaServer) PayCommissionsRevert(ctx context.Context, req *busi.ReqGrpc) (*empty.Empty, error) {
	log.Println("PayCommissionsRevert")
	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		return nil, err
	}
	if err := initDB(); err != nil {
		return nil, err
	}
	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		_, err = db.Exec("update user_account set balance= balance-? where id =3", 10)
		return err
	}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (GrpcSagaServer) TransOutHeaderYes(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransOutHeaderYes not implemented")
}
func (GrpcSagaServer) TransOutHeaderNo(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransOutHeaderNo not implemented")
}
func (GrpcSagaServer) TransInRedis(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransInRedis not implemented")
}
func (GrpcSagaServer) TransOutRedis(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransOutRedis not implemented")
}
func (GrpcSagaServer) TransInRevertRedis(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransInRevertRedis not implemented")
}
func (GrpcSagaServer) TransOutRevertRedis(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransOutRevertRedis not implemented")
}
func (GrpcSagaServer) QueryPrepared(context.Context, *busi.ReqGrpc) (*busi.BusiReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryPrepared not implemented")
}
func (GrpcSagaServer) QueryPreparedB(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryPreparedB not implemented")
}
func (GrpcSagaServer) QueryPreparedRedis(context.Context, *busi.ReqGrpc) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryPreparedRedis not implemented")
}
