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

	for _, field := range keyVals(ctx...) {
		c.tmp.ctx.fields = append(c.tmp.ctx.fields, field)
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

	// TODO: deduplicate fields with same name?
	field := &Field{
		Key: key,
		Val: value,
	}
	c.tmp.ctx.fields = append(c.tmp.ctx.fields, field)

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
		c.tmp.ctx.fields = append(c.tmp.ctx.fields, fields[k])
	}
	return c
}

func keyVals(keyvals ...interface{}) []*Field {
	if len(keyvals) == 0 {
		return nil
	}
	meta := make([]*Field, (len(keyvals)+1)/2)
	for i := 0; i < len(keyvals); i += 2 {
		k := keyvals[i]
		var v interface{} = "MISSING"
		if i+1 < len(keyvals) {
			v = keyvals[i+1]
		}
		field := &Field{
			Key: fmt.Sprint(k),
			Val: v,
		}
		meta = append(meta, field)
	}
	return meta
}

// ContextFields is just map[string]interface{}
type ContextFields []*Field
type Field struct {
	Key string
	Val interface{}
}

func newCtx() *contextFields {
	return &contextFields{
		fields: make(ContextFields, 0),
		mu:     &sync.RWMutex{},
	}
}

func (f ContextFields) String() string {
	s := ""
	for _, v := range f {
		if s != "" {
			s = s + "\n"
		}
		if v == nil {
			continue
		}
		s = s + fmt.Sprintf("*%v* = %v", v.Key, v.Val)
	}
	return s
}
