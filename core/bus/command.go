package bus

import (
	"errors"
	"time"
)

type Command struct {
	Message
	ttl         time.Duration
	expectReply bool
}

type CommandOption func(f *Command) error

func WithTtl(ttl time.Duration) CommandOption {
	return func(c *Command) error {
		if ttl <= time.Duration(0) {
			return errors.New("ttl must be greater than 0")
		}
		c.ttl = ttl
		return nil
	}
}

func WithExpectReply(expectReply bool) CommandOption {
	return func(c *Command) error {
		c.expectReply = expectReply
		return nil
	}
}

func NewCommand(name string, opts ...CommandOption) (*Command, error) {
	message, err := NewMessage(name)
	if err != nil {
		return nil, err
	}

	c := &Command{ttl: time.Minute, expectReply: true, Message: message}

	for _, applyOpt := range opts {
		err = applyOpt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Command) Ttl() time.Duration {
	return c.ttl
}

func (c *Command) ExpectReply() bool {
	return c.expectReply
}
