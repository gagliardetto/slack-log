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

type Msg interface {
	Blocks() []M
}

type Type string

const (
	PLAIN_TEXT Type = "plain_text"
	MARKDOWN   Type = "mrkdwn"
)

type Message struct {
	blocks []Block
}
type M map[string]interface{}

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

type Header struct {
	text  string
	emoji bool
}

func NewHeader() *Header {
	return &Header{}
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

//
func (bb *Header) ToMap() M {
	return M{
		"type": "header",
		"text": M{
			"type":  PLAIN_TEXT,
			"text":  bb.text,
			"emoji": bb.emoji,
		},
	}
}

//
func (bb *Header) Validate() error {
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

type Divider struct{}

func NewDivider() *Divider {
	return &Divider{}
}

//
func (bb *Divider) ToMap() M {
	return M{
		"type": "divider",
	}
}

//
func (bb *Divider) Validate() error {
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

//
func (bb *Message) Blocks() []M {
	blocks := make([]M, 0)

	for _, block := range bb.blocks {
		blocks = append(blocks, block.ToMap())
	}
	return blocks
}
func Link(url string, text string) string {
	return fmt.Sprintf("<%s|%s>", url, text)
}

type Image struct {
	title *Header
	alt   string
	url   string
}

func NewImage() *Image {
	return &Image{
		alt: "image",
	}
}

//
func (bb *Image) Title(title *Header) *Image {
	bb.title = title
	return bb
}

//
func (bb *Image) URL(url string) *Image {
	bb.url = url
	return bb
}

//
func (bb *Image) Alt(altText string) *Image {
	bb.alt = altText
	return bb
}

//
func (bb *Image) ToMap() M {
	m := M{
		"type":      "image",
		"image_url": bb.url,
		"alt_text":  bb.alt,
	}
	if bb.title != nil {
		m["title"] = M{
			"type":  PLAIN_TEXT,
			"text":  bb.title.text,
			"emoji": bb.title.emoji,
		}
	}
	return m
}

//
func (bb *Image) Validate() error {
	if bb.url == "" {
		return errors.New("image_url not set")
	}
	if bb.alt == "" {
		return errors.New("alt_text not set")
	}
	if bb.title != nil {
		if err := bb.title.Validate(); err != nil {
			return fmt.Errorf("error while validating title: %s", err)
		}
	}
	return nil
}
func example_blocks() {
	msg := &Message{}

	header := NewHeader().
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

	img := NewImage().
		URL("https://i1.wp.com/thetempest.co/wp-content/uploads/2017/08/The-wise-words-of-Michael-Scott-Imgur-2.jpg?w=1024&ssl=1").
		Alt("This is a quote").
		Title(
			NewHeader().
				Emoji(true).
				Text("An inspiring quote for *you*"),
		)
	msg.AddBlock(img)

	if err := msg.Validate(); err != nil {
		panic(err)
	}
}
