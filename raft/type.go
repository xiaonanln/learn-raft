package raft

type RaftInstance interface {
	ID() int
	Recv() <-chan interface{}
	Broadcast(msg interface{})
}

type AppendEntriesMessage struct {
}

type AppendEntriesACKMessage struct {
}

type RequestVoteMessage struct {
	// term candidate’s term
	term int
	// candidateId candidate requesting vote
	candidateId int
	// lastLogIndex index of candidate’s last log entry
	lastLogIndex int64
	//lastLogTerm term of candidate’s last log entry
	lastLogTerm int
}

type RequestVoteACKMessage struct {
}

type Log struct {
	term     int
	logIndex int64
}
