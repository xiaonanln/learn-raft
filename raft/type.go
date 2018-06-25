package raft

import (
	"fmt"
	"strings"
)

type Term int64
type LogIndex int64
type LogData []byte

func (d LogData) String() string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i, b := range d {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%02x", b))
	}
	sb.WriteByte(']')
	return sb.String()
}

type RaftInstance interface {
	ID() int
	Recv() <-chan RPCMessage
	Send(instanceID int, msg RPCMessage)
	Broadcast(msg RPCMessage)
	InputLog() <-chan LogData
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
	entries []*Log
	//leaderCommit leader’s commitIndex
	leaderCommit LogIndex
}

func (m *AppendEntriesMessage) Term() Term {
	return m.term
}

type AppendEntriesACKMessage struct {
	term    Term
	success bool
}

func (m *AppendEntriesACKMessage) Term() Term {
	return m.term
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
	term  Term
	index LogIndex
	data  LogData
}
