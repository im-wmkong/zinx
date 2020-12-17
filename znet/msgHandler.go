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
		TaskQueue:      make([]chan zicafe.IRequest, utils.GlobalObject.WorkerPoolSize),
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

func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan zicafe.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandle) StartOneWorker(workID int, taskQueue chan zicafe.IRequest) {
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request zicafe.IRequest) {
	workID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	mh.TaskQueue[workID] <- request
}
