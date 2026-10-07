package main

import (
	"errors"
	"time"
)

type PowerBackend interface {
	On(*Node) error
	Off(*Node) error
	Cycle(*Node, time.Duration) error
}

type ResetBackend interface { Reset(*Node) error }

type MediaBackend interface {
	Attach(*Node, Image) error
	Detach(*Node) error
}

type SimBackend struct{}

func (SimBackend) On(n *Node) error { n.Power="on"; n.Health="unknown"; return nil }
func (SimBackend) Off(n *Node) error { n.Power="off"; n.Health="unknown"; return nil }
func (SimBackend) Cycle(n *Node,d time.Duration) error {
	n.Power="cycling"; time.Sleep(d); n.Power="on"; n.Health="unknown"; return nil
}
func (SimBackend) Reset(n *Node) error {
	if n.Power!="on" { return errors.New("cannot reset powered-off node") }
	n.Health="unknown"; return nil
}
func (SimBackend) Attach(n *Node,im Image) error { n.AttachedImage=im.ID; return nil }
func (SimBackend) Detach(n *Node) error { n.AttachedImage=""; return nil }
