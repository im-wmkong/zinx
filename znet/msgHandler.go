package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/zicafe"
)

type MsgHandle struct {
	Apis           map[uint32]zicafe.IRouter
	TaskQueue      []chan zicafe.IRequest
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]zicafe.IRouter),
		TaskQueue:      make([]chan zicafe.IRequest, utils.GlobalObject.MaxWorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

func (mh *MsgHandle) DoMsgHandler(request zicafe.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Printf("msg handle is not found, id: %d\n", request.GetMsgId())
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}
func (mh *MsgHandle) AddRouter(msgID uint32, router zicafe.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		panic("msg handle is added, id: " + strconv.Itoa(int(msgID)))
	}
	mh.Apis[msgID] = router
	fmt.Printf("Add msg handle is success, id: %d\n", msgID)
}
