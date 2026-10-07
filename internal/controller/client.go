package controller

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	mu sync.Mutex
	rw io.ReadWriter
	seq uint64
}

func New(rw io.ReadWriter) *Client { return &Client{rw: rw} }

func (c *Client) Command(name string, args ...string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := atomic.AddUint64(&c.seq, 1)
	fields := append([]string{strconv.FormatUint(id, 10), name}, args...)
	if _, err := fmt.Fprintln(c.rw, strings.Join(fields, " ")); err != nil { return nil, err }
	line, err := bufio.NewReader(c.rw).ReadString('\n')
	if err != nil { return nil, err }
	resp := strings.Fields(line)
	if len(resp) < 3 { return nil, fmt.Errorf("malformed controller response") }
	if resp[0] != strconv.FormatUint(id, 10) { return nil, fmt.Errorf("response id mismatch: got %s want %d", resp[0], id) }
	switch resp[1] {
	case "OK":
		return resp[2:], nil
	case "ERR":
		return nil, fmt.Errorf("controller: %s", strings.Join(resp[2:], " "))
	default:
		return nil, fmt.Errorf("unknown controller response %q", resp[1])
	}
}

func (c *Client) Ping() error { _, err := c.Command("PING"); return err }
func (c *Client) PowerOn(port int) error { _, err := c.Command("POWER_ON", strconv.Itoa(port)); return err }
func (c *Client) PowerOff(port int) error { _, err := c.Command("POWER_OFF", strconv.Itoa(port)); return err }
func (c *Client) PowerCycle(port int, off time.Duration) error {
	_, err := c.Command("POWER_CYCLE", strconv.Itoa(port), strconv.FormatInt(off.Milliseconds(), 10))
	return err
}
func (c *Client) Reset(port int, pulse time.Duration) error {
	_, err := c.Command("RESET", strconv.Itoa(port), strconv.FormatInt(pulse.Milliseconds(), 10))
	return err
}
func (c *Client) ServiceSelect(port int) error { _, err := c.Command("SERVICE_SELECT", strconv.Itoa(port)); return err }
func (c *Client) AllSafe() error { _, err := c.Command("ALL_SAFE"); return err }
