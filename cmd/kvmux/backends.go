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

type SimBackend struct {
	FailPowerOn bool
	FailPowerOff bool
	FailMediaAttach bool
}

func (b SimBackend) On(n *Node) error { if b.FailPowerOn { return errors.New("injected power-on failure") }; n.Power="on"; n.Health="unknown"; return nil }
func (b SimBackend) Off(n *Node) error { if b.FailPowerOff { return errors.New("injected power-off failure") }; n.Power="off"; n.Health="unknown"; return nil }
func (SimBackend) Cycle(n *Node,d time.Duration) error {
	n.Power="cycling"; time.Sleep(d); n.Power="on"; n.Health="unknown"; return nil
}
func (SimBackend) Reset(n *Node) error {
	if n.Power!="on" { return errors.New("cannot reset powered-off node") }
	n.Health="unknown"; return nil
}
func (b SimBackend) Attach(n *Node,im Image) error { if b.FailMediaAttach { return errors.New("injected media-attach failure") }; n.AttachedImage=im.ID; return nil }
func (SimBackend) Detach(n *Node) error { n.AttachedImage=""; return nil }
