package raft

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type workMode int

func (m workMode) String() string {
	return [3]string{"follower", "candidate", "leader"}[m]
}

const (
	followerMode workMode = iota
	candidateMode
	leaderMode
)

const (
	electionTimeout       = time.Millisecond * 3000
	maxStartElectionDelay = electionTimeout / 2
)

type Raft struct {
	ins                      RaftInstance
	ctx                      context.Context
	instanceNum              int
	mode                     workMode
	resetElectionTimeoutTime time.Time

	// raft states
	currentTerm int
	votedFor    int

	// candidate mode
	startElectionTime time.Time
	electionStarted   bool
	log               []Log
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
		resetElectionTimeoutTime: time.Now(),
		// init raft states
		currentTerm: 0,
		votedFor:    -1,
		log:         []Log{},
	}

	go raft.routine()

	return raft
}

func (r *Raft) ID() int {
	return r.ins.ID()
}

func (r *Raft) String() string {
	return fmt.Sprintf("Raft<%d>", r.ins.ID())
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
				r.candidateTick()
			case followerMode:
				r.followerTick()
			default:
				log.Fatalf("invalid mode: %d", r.mode)
			}

		case msg := <-r.ins.Recv():
			log.Printf("%s received msg: %+v", r, msg)
			r.handleMsg(msg)

		case <-r.ctx.Done():
			log.Printf("%s stopped.", r)
			break forloop
		}
	}
}

func (r *Raft) handleMsg(_msg interface{}) {
	//Current terms are exchanged whenever servers communicate; if one server’s current
	//term is smaller than the other’s, then it updates its current
	//term to the larger value. If a candidate or leader discovers
	//that its term is out of date, it immediately reverts to follower
	//state. If a server receives a request with a stale term
	//number, it rejects the request.
	switch msg := _msg.(type) {
	case *RequestVoteMessage:
		r.handleRequestVote(msg)
	default:
		log.Fatalf("unexpected _msg type: %T", _msg)
	}
}

func (r *Raft) handleRequestVote(msg *RequestVoteMessage) {
	if msg.term > r.currentTerm {
		// RequestVote with higher term, fall back to follower mode
		if r.mode != followerMode {
			r.enterMode(followerMode)
		}
	}

	grantVote := r._handleRequestVote(msg)
	log.Printf("%s grant vote %+v: %v", r, msg, grantVote)
}

func (r *Raft) _handleRequestVote(msg *RequestVoteMessage) bool {
	//1. Reply false if term < currentTerm (§5.1)
	//2. If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)

	if msg.term < r.currentTerm {
		return false
	}

	grantVote := (r.votedFor == -1 || r.votedFor == msg.candidateId) && r.isLogUpToDate(msg.lastLogTerm, msg.lastLogIndex)
	return grantVote
}

// isLogUpToDate determines if the log of specified term and index is at least as up-to-date as r.log
func (r *Raft) isLogUpToDate(term int, logIndex int64) bool {
	myterm := r.lastLogTerm()
	if term < myterm {
		return false
	} else if term > myterm {
		return true
	}

	return logIndex >= r.lastLogIndex()
}

func (r *Raft) followerTick() {
	now := time.Now()
	if now.Sub(r.resetElectionTimeoutTime) > electionTimeout {
		r.enterMode(candidateMode)
	}
}

func (r *Raft) candidateTick() {
	now := time.Now()

	if now.Sub(r.resetElectionTimeoutTime) > electionTimeout {
		log.Printf("%s: election timeout in candidate mode, restarting election ...", r)
		r.prepareElection()
		return
	}

	if !r.electionStarted {
		if now.After(r.startElectionTime) {
			r.startElection()
		}
	}

}

func (r *Raft) resetElectionTimeout() {
	r.resetElectionTimeoutTime = time.Now()
}

func (r *Raft) prepareElection() {
	r.assureInMode(candidateMode)

	log.Printf("%s prepare election ...", r)
	r.startElectionTime = time.Now().Add(time.Duration(rand.Int63n(int64(maxStartElectionDelay))))
	r.electionStarted = false
	r.votedFor = -1
	r.resetElectionTimeout()
	log.Printf("%s set start election time = %s", r, r.startElectionTime)
}

func (r *Raft) startElection() {
	r.assureInMode(candidateMode)

	if r.electionStarted {
		log.Panicf("election already started")
	}
	//On conversion to candidate, start election:
	//?Increment currentTerm
	//?Vote for self
	//?Reset election timer
	//?Send RequestVote RPCs to all other servers
	log.Printf("%s start election ...", r)
	r.currentTerm += 1
	r.electionStarted = true
	r.votedFor = r.ID() // vote for self
	r.sendRequestVote()
}

func (r *Raft) sendRequestVote() {
	//Arguments:
	//term candidate’s term
	//candidateId candidate requesting vote
	//lastLogIndex index of candidate’s last log entry (§5.4)
	//lastLogTerm term of candidate’s last log entry (§5.4)
	msg := &RequestVoteMessage{
		term:         r.currentTerm,
		candidateId:  r.ID(),
		lastLogIndex: r.lastLogIndex(),
		lastLogTerm:  r.lastLogTerm(),
	}
	r.ins.Broadcast(msg)
}

func (r *Raft) assureInMode(mode workMode) {
	if r.mode != mode {
		log.Fatalf("%s should in %s mode, but in %s mode", r, mode, r.mode)
	}
}

func (r *Raft) enterMode(mode workMode) {
	if r.mode == mode {
		log.Panicf("%s: already in mode %d", r, mode)
	}

	r.mode = mode
	log.Printf("%s enter %s mode", r, r.mode)
	switch r.mode {
	case followerMode:
	case candidateMode:
		r.prepareElection()
	case leaderMode:
	}
}
func (r *Raft) lastLogIndex() int64 {
	if len(r.log) > 0 {
		return r.log[len(r.log)-1].logIndex
	} else {
		return 0
	}
}

func (r *Raft) lastLogTerm() int {
	if len(r.log) > 0 {
		return r.log[len(r.log)-1].term
	} else {
		return 0
	}
}
