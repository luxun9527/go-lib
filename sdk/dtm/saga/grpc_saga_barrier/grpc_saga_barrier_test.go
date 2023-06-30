package busi

import (
	"fmt"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/lithammer/shortuuid/v3"
	"testing"
)

var BusiGrpc = fmt.Sprintf("localhost:%d", 58081)

func TestGrpcSagaBarrier(t *testing.T) {
	req := &ReqGrpc{Amount: 30}
	gid := shortuuid.New()
	saga := dtmgrpc.NewSagaGrpc("localhost:36790", gid).
		Add(BusiGrpc+"/busi.Busi/TransOutBSaga", BusiGrpc+"/busi.Busi/TransOutRevertBSaga", req).
		Add(BusiGrpc+"/busi.Busi/TransInBSaga", BusiGrpc+"/busi.Busi/TransInRevertBSaga", req)
	_ = saga.Submit()
}
