package main

import (
	"context"

	"github.com/xiaonanln/learn-raft/raft"
	"github.com/xiaonanln/learn-raft/raft-demo/demo"
)

const (
	INSTANCE_NUM = 3
)

func main() {
	ctx := context.Background()
	for i := 0; i < INSTANCE_NUM; i++ {
		ins := demo.NewDemoRaftInstance(ctx, i)
		raft.NewRaft(ctx, INSTANCE_NUM, ins)
	}

	<-ctx.Done()
}
