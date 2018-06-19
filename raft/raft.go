package raft

import (
	"context"
	"fmt"
	"log"
	"time"
)

type workMode int

const (
	followerMode workMode = iota
	candidateMode
	leaderMode
)

const (
	electionTimeout = time.Millisecond * 500
)

type Raft struct {
	ins                       RaftInstance
	ctx                       context.Context
	instanceNum               int
	mode                      workMode
	lastAppendEntriesRecvTime time.Time

	// raft states
	currentTerm int
	votedFor    []int
}

func NewRaft(ctx context.Context, instanceNum int, ins RaftInstance) *Raft {
	if instanceNum <= 0 {
		log.Fatalf("instanceNum should be larger than 0")
	}

	if ins.ID() >= instanceNum {
		log.Fatalf("instance ID should be smaller than %d", instanceNum)
	}

	raft := &Raft{
		ins:         ins,
		ctx:         ctx,
		instanceNum: instanceNum,
		mode:        followerMode,
		lastAppendEntriesRecvTime: time.Now(),
		// init raft states
		currentTerm: 0,
		votedFor:    make([]int, instanceNum),
	}

	raft.clearVotedFor()
	go raft.routine()

	return raft
}

func (r *Raft) ID() int {
	return r.ins.ID()
}

func (r *Raft) String() string {
	return fmt.Sprintf("Raft<%d>", r.ins.ID())
}
func (r *Raft) clearVotedFor() {
	for i := 0; i < r.instanceNum; i++ {
		r.votedFor[i] = -1
	}
}

func (r *Raft) routine() {
	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()
	log.Printf("%s start running ...", r)

forloop:
	for {
		select {
		case <-ticker.C:
			switch r.mode {
			case leaderMode:
			case candidateMode:
			case followerMode:
				r.followerTick()
			default:
				log.Fatalf("invalid mode: %d", r.mode)
			}

		case <-r.ins.Recv():

		case <-r.ctx.Done():
			log.Printf("%s stopped.", r)
			break forloop
		}
	}
}
func (r *Raft) followerTick() {
	now := time.Now()
	if now.Sub(r.lastAppendEntriesRecvTime) > electionTimeout {
		r.enterMode(candidateMode)
	}
}

func (r *Raft) enterMode(mode workMode) {
	if r.mode == mode {
		log.Panicf("%s: already in mode %d", r, mode)
	}

	r.mode = mode
	log.Printf("%s enter mode %d", r, r.mode)
	switch r.mode {
	case followerMode:
	case candidateMode:
	case leaderMode:
	}
}
