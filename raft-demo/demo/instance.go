package demo

import (
	"context"
	"fmt"
	"log"

	"github.com/xiaonanln/learn-raft/raft"
)

type DemoRaftInstance struct {
	ctx      context.Context
	id       int
	recvChan chan raft.RPCMessage
}

var (
	instances = map[int]*DemoRaftInstance{}
)

func NewDemoRaftInstance(ctx context.Context, id int) *DemoRaftInstance {
	ins := &DemoRaftInstance{
		ctx:      ctx,
		id:       id,
		recvChan: make(chan raft.RPCMessage, 1000),
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

func (ins *DemoRaftInstance) Recv() <-chan raft.RPCMessage {
	return ins.recvChan
}

func (ins *DemoRaftInstance) Send(insID int, msg raft.RPCMessage) {
	instances[insID].recvChan <- msg
}

// Broadcast sends message to all other instances
func (ins *DemoRaftInstance) Broadcast(msg raft.RPCMessage) {
	log.Printf("%s BROADCAST: %+v", ins, msg)
	for _, other := range instances {
		if other == ins {
			continue
		}

		other.recvChan <- msg
	}
}
