package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Raft roles
const (
	Follower  = "Follower"
	Candidate = "Candidate"
	Leader    = "Leader"
)

// Node represents a Raft node
type Node struct {
	ID         int
	state      string
	voteCount  int
	term       int
	peers      []*Node
	mu         sync.Mutex
	voteCh     chan bool
	heartbeat  chan bool
	stopCh     chan any
	electionTO time.Duration
}

// NewNode create a new Raft node
func NewNode(id int) *Node {
	return &Node{
		ID:         id,
		state:      Follower,
		term:       0,
		peers:      make([]*Node, 0),
		voteCh:     make(chan bool, 1),
		heartbeat:  make(chan bool, 1),
		stopCh:     make(chan any, 1),
		electionTO: time.Duration(rand.Intn(150)+150) * time.Millisecond, // 150~300ms election timeout
	}
}

func (n *Node) AddPeers(peers ...*Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.peers = append(n.peers, peers...)
}

func (n *Node) RemovePeer(peer *Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	newPeers := make([]*Node, 0, len(n.peers)-1)
	for _, p := range n.peers {
		if p != nil && p.ID != peer.ID {
			newPeers = append(newPeers, p)
		}
	}
	n.peers = newPeers
}

// Start the node
func (n *Node) Start() {
	go n.run()
}

// Stop shutdown the node
func (n *Node) Stop() {
	n.stopCh <- struct{}{}
}

// run a role as state
func (n *Node) run() {
	fmt.Printf("> Member %d: Hi\n", n.ID)
	for {
		select {
		case <-n.stopCh:
			fmt.Printf("> Member %d: Manually killed by user\n", n.ID)
			return
		default:
			switch n.state {
			case Follower:
				n.runFollower()
			case Candidate:
				n.runCandidate()
			case Leader:
				n.runLeader()
			}
		}
	}
}

// runFollower
func (n *Node) runFollower() {
	select {
	case <-n.heartbeat:
		// Received leader heartbeat, do nothing
	case <-time.After(n.electionTO):
		// Timeout, transition to Candidate and start election
		n.mu.Lock()
		n.state = Candidate
		n.mu.Unlock()
	}
}

// runCandidate
func (n *Node) runCandidate() {
	n.mu.Lock()
	n.term++ // incr term
	n.voteCount = 1
	fmt.Printf("> Member %d: I want to be leader for term %d\n", n.ID, n.term)
	n.mu.Unlock()

	// send RequestVote to other nodes
	for _, peer := range n.peers {
		if peer != n {
			go func(p *Node) {
				vote := p.RequestVote(n.term, n.ID)
				if vote {
					n.voteCh <- true
				}
			}(peer)
		}
	}

	// wait for election result
	votesReceived := 1
	for {
		select {
		case <-n.voteCh:
			votesReceived++
			if votesReceived > (len(n.peers)+1)/2 {
				n.mu.Lock()
				n.state = Leader
				n.mu.Unlock()
				fmt.Printf("> Member %d voted to be leader: (%d > %d/2)\n", n.ID, votesReceived, len(n.peers)+1)
				return
			}
		case <-time.After(n.electionTO):
			// election failed, revert to follower
			n.mu.Lock()
			n.state = Follower
			n.mu.Unlock()
			return
		}
	}
}

// runLeader
func (n *Node) runLeader() {
	for _, peer := range n.peers {
		if peer != nil && peer != n {
			go peer.ReceiveHeartbeat(n.term)
		}
	}
	time.Sleep(100 * time.Millisecond)
}

// RequestVote candidate request vote from peers
func (n *Node) RequestVote(term int, candidateID int) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	if term > n.term {
		n.term = term
		n.state = Follower
		fmt.Printf("> Member %d: Accept member %d to be leader, Term: %d\n", n.ID, candidateID, term)
		return true
	}
	return false
}

// ReceiveHeartbeat follower received a heartbeat from leader
func (n *Node) ReceiveHeartbeat(term int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if term >= n.term {
		n.term = term
		n.state = Follower
		n.heartbeat <- true
	}
}
