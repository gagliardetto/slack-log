package slacklog

import (
	"errors"
	"fmt"
	"strings"
)

type Block interface {
	Validate() error
	ToMap() M
}

type Type string

const (
	PLAIN_TEXT Type = "plain_text"
	MARKDOWN   Type = "mrkdwn"
)

type Header struct {
	typ   Type
	text  string
	emoji bool
}

type Message struct {
	blocks []Block
}

//
func (msg *Message) AddBlock(bl Block) {
	msg.blocks = append(msg.blocks, bl)
}

//
func (msg *Message) Validate() error {
	for i, bl := range msg.blocks {
		if err := bl.Validate(); err != nil {
			return fmt.Errorf("block #%v is not valid: %s", i, err)
		}
	}
	return nil
}
func NewHeader() *Header {
	return &Header{
		typ: PLAIN_TEXT,
	}
}

//
func (bb *Header) Type(typ Type) *Header {
	bb.typ = typ
	return bb
}

//
func (bb *Header) Text(text string) *Header {
	bb.text = text
	return bb
}

//
func (bb *Header) Emoji(emoji bool) *Header {
	bb.emoji = emoji
	return bb
}

type M map[string]interface{}

//
func (bb *Header) ToMap() M {
	return M{
		"type": "header",
		"text": M{
			"type":  bb.typ,
			"text":  bb.text,
			"emoji": bb.emoji,
		},
	}
}

//
func (bb *Header) Validate() error {
	if bb.typ == "" {
		return errors.New("type not set")
	}
	if bb.text == "" {
		return errors.New("text not set")
	}
	return nil
}

type Section struct {
	fields []*SField
}

func NewSection() *Section {
	return &Section{
		fields: make([]*SField, 0),
	}
}

//
func (bb *Section) AddField(fld *SField) *Section {
	bb.fields = append(bb.fields, fld)
	return bb
}

//
func (bb *Section) Validate() error {
	for i, fld := range bb.fields {
		if err := fld.Validate(); err != nil {
			return fmt.Errorf("field #%v is not valid: %s", i, err)
		}
	}
	return nil
}

//
func (bb *Section) ToMap() M {
	fieldMaps := make([]M, 0)
	for _, fld := range bb.fields {
		fieldMaps = append(fieldMaps, fld.ToMap())
	}
	return M{
		"type":   "section",
		"fields": fieldMaps,
	}
}

type SField struct {
	typ   Type
	lines []string
}

func NewField() *SField {
	return &SField{}
}

//
func (bb *SField) Type(typ Type) *SField {
	bb.typ = typ
	return bb
}

//
func (bb *SField) Text(lines ...string) *SField {
	bb.lines = append(bb.lines, lines...)
	return bb
}

//
func (bb *SField) Validate() error {
	if bb.typ == "" {
		return errors.New("type not set")
	}
	if len(bb.lines) == 0 {
		return errors.New("text not set")
	}
	return nil
}

//
func (bb *SField) ToMap() M {
	return M{
		"type": bb.typ,
		"text": strings.Join(bb.lines, "\n"),
	}
}

func Link(url string, text string) string {
	return fmt.Sprintf("<%s|%s>", url, text)
}
func example_blocks() {
	msg := &Message{}

	header := NewHeader().
		Type(PLAIN_TEXT).
		Text("Hello world").
		Emoji(true)

	msg.AddBlock(header)

	section1 := NewSection().
		AddField(
			NewField().
				Type(PLAIN_TEXT).
				Text("Foo"),
		)

	msg.AddBlock(section1)

	section2 := NewSection().
		AddField(
			NewField().
				Type(MARKDOWN).
				Text(
					"*Hello* **world**",
					"This is another line.",
				),
		)

	msg.AddBlock(section2)

	if err := msg.Validate(); err != nil {
		panic(err)
	}
}
