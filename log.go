package slacklog

import (
	"fmt"
	"sync"
)

type LogMessage struct {
	tmp *temporaryData
}

func NewLogMessage() *LogMessage {
	return &LogMessage{
		tmp: newTmp(),
	}
}

type temporaryData struct {
	ctx     *contextFields
	message string
}

func newTmp() *temporaryData {
	return &temporaryData{
		ctx: newCtx(),
	}
}

type contextFields struct {
	fields Fields
	mu     *sync.RWMutex
}

func (cf *contextFields) len() int {
	return len(cf.fields)
}

// Log15 lets you add context fields to the message in the style of Log15
func (c *LogMessage) Log15(title string, ctx ...interface{}) *LogMessage {
	if c.tmp == nil {
		c.tmp = newTmp()
	}
	c.tmp.ctx.mu.Lock()
	defer c.tmp.ctx.mu.Unlock()

	for _, field := range keyVals(ctx...) {
		c.tmp.ctx.fields = append(c.tmp.ctx.fields, field)
	}
	c.tmp.message = title
	return c
}

// WithField lets you add a context field to the message
func (c *LogMessage) WithField(key string, value interface{}) *LogMessage {
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
func (c *LogMessage) WithFields(fields Fields) *LogMessage {
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
	fields := make([]*Field, 0)
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
		fields = append(fields, field)
	}
	return fields
}

type Fields []*Field

type Field struct {
	Key string
	Val interface{}
}

func newCtx() *contextFields {
	return &contextFields{
		fields: make(Fields, 0),
		mu:     &sync.RWMutex{},
	}
}

func (f Fields) String() string {
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

// SetTitle sets the main message text
func (c *LogMessage) SetTitle(msg string) *LogMessage {
	if c.tmp == nil {
		c.tmp = newTmp()
	}

	c.tmp.message = msg
	return c
}

//
func (bb *LogMessage) Validate() error {
	return nil
}

//
func (lm *LogMessage) Blocks() []M {

	msg := &Message{}

	if lm.tmp.message != "" {
		header := NewHeader().
			Text(lm.tmp.message).
			Emoji(true)

		msg.AddBlock(header)
	}

	{
		if lm.tmp.ctx.len() > 0 {
			section := NewSection().
				AddField(
					NewField().
						Type(MARKDOWN).
						Text(
							lm.tmp.ctx.fields.String(),
						),
				)
			msg.AddBlock(section)
		}
	}

	if err := msg.Validate(); err != nil {
		panic(err)
	}

	return msg.Blocks()
}
