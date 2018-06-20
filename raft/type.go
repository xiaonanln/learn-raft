package raft

type RaftInstance interface {
	ID() int
	Recv() <-chan interface{}
}

type AppendEntriesMessage struct {
}

type AppendEntriesACKMessage struct {
}

type RequestVoteMessage struct {
}

type RequestVoteACKMessage struct {
}
