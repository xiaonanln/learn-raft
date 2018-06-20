package demo

import "context"

type DemoRaftInstance struct {
	ctx      context.Context
	id       int
	recvChan chan interface{}
}

var (
	instances = map[int]*DemoRaftInstance{}
)

func NewDemoRaftInstance(ctx context.Context, id int) *DemoRaftInstance {
	ins := &DemoRaftInstance{
		ctx:      ctx,
		id:       id,
		recvChan: make(chan interface{}, 1000),
	}

	instances[id] = ins
	return ins
}

func (ins *DemoRaftInstance) ID() int {
	return ins.id
}

func (ins *DemoRaftInstance) Recv() <-chan interface{} {
	return ins.recvChan
}
