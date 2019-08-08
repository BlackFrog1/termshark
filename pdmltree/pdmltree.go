// Copyright 2019 Graham Clark. All rights reserved.  Use of this source
// code is governed by the MIT license that can be found in the LICENSE
// file.

// A demonstration of gowid's tree widget.
package pdmltree

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/tree"
	"github.com/gcla/termshark/widgets/hexdumper"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//======================================================================

type EmptyIterator struct{}

var _ tree.IIterator = EmptyIterator{}

func (e EmptyIterator) Next() bool {
	return false
}

func (e EmptyIterator) Value() tree.IModel {
	panic(errors.New("Should not call"))
}

// pos points one head, so logically is -1 on init, but I use zero so the
// go default init makes sense.
type Iterator struct {
	tree *Model
	pos  int
}

var _ tree.IIterator = (*Iterator)(nil)

func (p *Iterator) Next() bool {
	p.pos += 1
	return (p.pos - 1) < len(p.tree.Children_)
}

func (p *Iterator) Value() tree.IModel {
	return p.tree.Children_[p.pos-1]
}

type Model struct {
	UiName    string            `xml:"-"`
	Name      string            `xml:"-"` // needed for stripping general info from UI
	Expanded  bool              `xml:"-"`
	Pos       int               `xml:"-"`
	Size      int               `xml:"-"`
	Hide      bool              `xml:"-"`
	Children_ []*Model          `xml:",any"`
	Content   []byte            `xml:",innerxml"` // needed for copying Packet Description Markup Language to clipboard
	NodeName  string            `xml:"-"`
	Attrs     map[string]string `xml:"-"`
}

var _ tree.IModel = (*Model)(nil)

// This ignores the first child, "Frame 15", because its range covers the whole packet
// which results in me always including that in the layers for any position.
func (m *Model) HexLayers(pos int, includeFirst bool) []hexdumper.LayerStyler {
	res := make([]hexdumper.LayerStyler, 0)
	sidX := 1
	if includeFirst {
		sidX = 0
	}
	for _, c := range m.Children_[sidX:] {
		if c.Pos <= pos && pos < c.Pos+c.Size {
			res = append(res, hexdumper.LayerStyler{
				Start:         c.Pos,
				End:           c.Pos + c.Size,
				ColUnselected: "hex-bottom-unselected",
				ColSelected:   "hex-bottom-selected",
			})
			for _, c2 := range c.Children_ {
				if c2.Pos <= pos && pos < c2.Pos+c2.Size {
					res = append(res, hexdumper.LayerStyler{
						Start:         c2.Pos,
						End:           c2.Pos + c2.Size,
						ColUnselected: "hex-top-unselected",
						ColSelected:   "hex-top-selected",
					})
				}
			}
		}
	}
	return res
}

// Implement xml.Unmarshal handler
func (m *Model) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error
	m.Attrs = map[string]string{}
	for _, a := range start.Attr {
		m.Attrs[a.Name.Local] = a.Value
		switch a.Name.Local {
		case "pos":
			m.Pos, err = strconv.Atoi(a.Value)
			if err != nil {
				return errors.WithStack(err)
			}
		case "size":
			m.Size, err = strconv.Atoi(a.Value)
			if err != nil {
				return errors.WithStack(err)
			}
		case "showname":
			m.UiName = a.Value
		case "show":
			if m.UiName == "" {
				m.UiName = a.Value
			}
		case "hide":
			m.Hide = a.Value == "yes"
		case "name":
			m.Name = a.Value
		}
	}

	m.NodeName = start.Name.Local

	type pt Model
	res := d.DecodeElement((*pt)(m), &start)
	return res
}

func DecodePacket(data []byte) *Model { // nil if failure
	d := xml.NewDecoder(bytes.NewReader(data))

	var n Model
	err := d.Decode(&n)
	if err != nil {
		log.Error(err)
		return nil
	}

	tr := n.removeUnneeded()
	return tr
}

func (m *Model) removeUnneeded() *Model {
	if m.Hide {
		return nil
	}
	if m.Name == "geninfo" {
		return nil
	}
	if m.Name == "fake-field-wrapper" { // for now...
		return nil
	}
	ch := make([]*Model, 0, len(m.Children_))
	for _, c := range m.Children_ {
		nc := c.removeUnneeded()
		if nc != nil {
			ch = append(ch, nc)
		}
	}
	m.Children_ = ch
	return m
}

func (m *Model) Children() tree.IIterator {
	if m.Expanded {
		return &Iterator{
			tree: m,
		}
	} else {
		return EmptyIterator{}
	}
}

func (m *Model) HasChildren() bool {
	return len(m.Children_) > 0
}

func (m *Model) Leaf() string {
	return m.UiName
}

func (m *Model) String() string {
	return m.stringAt(1)
}

func (m *Model) stringAt(level int) string {
	x := make([]string, len(m.Children_))
	for i, t := range m.Children_ {
		//x[i] = t.(*ModelExt).String2(level + 1)
		x[i] = t.stringAt(level + 1)
	}
	for i := range x {
		x[i] = strings.Repeat(" ", 2*level) + x[i]
	}
	if len(x) == 0 {
		return fmt.Sprintf("[%s]", m.UiName)
	} else {
		return fmt.Sprintf("[%s]\n%s", m.UiName, strings.Join(x, "\n"))
	}
}

//func (p *Model) Children() tree.IIterator {
//}

func (m *Model) IsCollapsed() bool {
	//return false
	return !m.Expanded
	// fp := d.FullPath()
	// if v, res := (*d.cache)[fp]; res {
	// 	return (v == collapsed)
	// } else {
	// 	return true
	// }
}

func (m *Model) SetCollapsed(app gowid.IApp, isCollapsed bool) {
	// fp := d.FullPath()
	if isCollapsed {
		m.Expanded = false
	} else {
		m.Expanded = true
	}
}

//======================================================================
// Local Variables:
// mode: Go
// fill-column: 110
// End:
