package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter // apis[id]=router
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("DoMsgHandler panic: %v", r)
			return
		}
	}()
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("api msgId = %d , is not found", request.GetMsgID())
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := m.Apis[msgId]; ok {
		panic("repeated api,msgId = " + strconv.Itoa(int(msgId)))
	}
	m.Apis[msgId] = router
	fmt.Printf("add api mshid %d", msgId)
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandler) StartOneWorker(workId int, taskQueue chan ziface.IRequest) {
	fmt.Printf("Work id :%d ,is started", workId)
	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

func (m *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := uint32(len(m.TaskQueue)) % m.WorkerPoolSize
	fmt.Printf("Add ConnID=%d request msgID = %d to workerID = %d", request.GetConnection().GetConnID(), request.GetMsgID(), workerID)
	m.TaskQueue[workerID] <- request
}
