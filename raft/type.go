package raft

type Term int64
type LogIndex int64

type RaftInstance interface {
	ID() int
	Recv() <-chan RPCMessage
	Send(instanceID int, msg RPCMessage)
	Broadcast(msg RPCMessage)
}

type RPCMessage interface {
	Term() Term
}

type AppendEntriesMessage struct {
	//Term leader’s Term
	term Term
	//leaderId so follower can redirect clients
	leaderId int
	//prevLogIndex index of log entry immediately preceding new ones
	prevLogIndex LogIndex
	//prevLogTerm Term of prevLogIndex entry
	prevLogTerm Term
	//entries[] log entries to store (empty for heartbeat; may send more than one for efficiency)
	entries []Log
	//leaderCommit leader’s commitIndex
	leaderCommit LogIndex
}

func (m *AppendEntriesMessage) Term() Term {
	return m.term
}

type AppendEntriesACKMessage struct {
}

type RequestVoteMessage struct {
	// Term candidate’s Term
	term Term
	// candidateId candidate requesting vote
	candidateId int
	// lastLogIndex index of candidate’s last log entry
	lastLogIndex LogIndex
	//lastLogTerm Term of candidate’s last log entry
	lastLogTerm Term
}

func (m *RequestVoteMessage) Term() Term {
	return m.term
}

type RequestVoteACKMessage struct {
	term        Term
	voteGranted bool
}

func (m *RequestVoteACKMessage) Term() Term {
	return m.term
}

type Log struct {
	term     Term
	logIndex LogIndex
}
