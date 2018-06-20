package demo

import (
	"context"
	"fmt"
	"log"
)

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

func (ins *DemoRaftInstance) String() string {
	return fmt.Sprintf("RaftInstance<%d>", ins.id)
}
func (ins *DemoRaftInstance) ID() int {
	return ins.id
}

func (ins *DemoRaftInstance) Recv() <-chan interface{} {
	return ins.recvChan
}

// Broadcast sends message to all other instances
func (ins *DemoRaftInstance) Broadcast(msg interface{}) {
	log.Printf("%s BROADCAST: %+v", ins, msg)
	for _, other := range instances {
		if other == ins {
			continue
		}

		ins.recvChan <- msg
	}
}
