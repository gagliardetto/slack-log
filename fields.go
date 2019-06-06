package slacklog

import (
	"fmt"
	"sync"
)

type contextFields struct {
	fields ContextFields
	mu     *sync.RWMutex
}

func (cf *contextFields) len() int {
	return len(cf.fields)
}

// Log15 lets you add context fields to the message in the style of Log15
func (c *Client) Log15(msg string, ctx ...interface{}) ClientInterface {
	if c.tmp == nil {
		c.tmp = newTmp()
	}
	c.tmp.ctx.mu.Lock()
	defer c.tmp.ctx.mu.Unlock()

	for k, v := range keyVals(ctx...) {
		c.tmp.ctx.fields[k] = v
	}
	c.tmp.message = msg
	return c
}

// WithField lets you add a context field to the message
func (c *Client) WithField(key string, value interface{}) ClientInterface {
	if c.tmp == nil {
		c.tmp = newTmp()
	}
	c.tmp.ctx.mu.Lock()
	defer c.tmp.ctx.mu.Unlock()

	c.tmp.ctx.fields[key] = value
	return c
}

// WithFields lets you add context fields to the message
func (c *Client) WithFields(fields ContextFields) ClientInterface {
	if c.tmp == nil {
		c.tmp = newTmp()
	}
	c.tmp.ctx.mu.Lock()
	defer c.tmp.ctx.mu.Unlock()

	for k := range fields {
		c.tmp.ctx.fields[k] = fields[k]
	}
	return c
}

func keyVals(keyvals ...interface{}) map[string]interface{} {
	if len(keyvals) == 0 {
		return nil
	}
	meta := make(map[string]interface{}, (len(keyvals)+1)/2)
	for i := 0; i < len(keyvals); i += 2 {
		k := keyvals[i]
		var v interface{} = "MISSING"
		if i+1 < len(keyvals) {
			v = keyvals[i+1]
		}
		meta[fmt.Sprint(k)] = v
	}
	return meta
}

// ContextFields is just map[string]interface{}
type ContextFields map[string]interface{}

func newCtx() *contextFields {
	return &contextFields{
		fields: make(ContextFields),
		mu:     &sync.RWMutex{},
	}
}

func (f ContextFields) String() string {
	s := ""
	for k, v := range f {
		if s != "" {
			s = s + "\n"
		}
		s = s + fmt.Sprintf("*%v* = %v", k, v)
	}
	return s
}
