package lua

import (
	"github.com/charmbracelet/bubbles/textinput"
	lua "github.com/yuin/gopher-lua"
)

type Component interface {
	ID() string
	UpdateFromLua(*lua.LTable)
	Render() string
}

type TextInputComponent struct {
	Model textinput.Model
	Id    string
}

func (t *TextInputComponent) ID() string {
	return t.Id
}

func (t *TextInputComponent) UpdateFromLua(tbl *lua.LTable) {
	if v := tbl.RawGetString("placeholder"); v.Type() == lua.LTString {
		t.Model.Placeholder = v.String()
	}
	if v := tbl.RawGetString("value"); v.Type() == lua.LTString {
		t.Model.SetValue(v.String())
	}
}

func (t *TextInputComponent) Render() string {
	return t.Model.View()
}

type Label struct {
	Text string
}

func (l *Label) ID() string { return "" }

func (l *Label) UpdateFromLua(tbl *lua.LTable) {
	if val, ok := tbl.RawGetString("text").(lua.LString); ok {
		l.Text = string(val)
	}
}

func (l *Label) Render() string {
	return l.Text
}
