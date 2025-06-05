package lua

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/neo451/ayo/internal/characters"
	"github.com/neo451/ayo/internal/config"
	lua "github.com/yuin/gopher-lua"
)

type LuaModel struct {
	L            *lua.LState
	ui           *LuaUI
	msg          tea.Msg
	components   []Component
	componentMap map[string]Component
}

func exposeUI(L *lua.LState, m *LuaModel) {
	uiTable := L.NewTable()

	// Create a new render function that updates ui.Rendered from a Lua array
	uiTable.RawSetString("render", L.NewFunction(func(L *lua.LState) int {
		arr := L.ToTable(1)
		var sb strings.Builder
		arr.ForEach(func(_, val lua.LValue) {
			sb.WriteString(val.String())
			sb.WriteString("\n")
		})
		m.ui.Rendered = sb.String()
		return 0
	}))

	uiTable.RawSetString("textInput", TI(L, m))

	// Push to global scope
	L.SetGlobal("ui", uiTable)
}

func (m *LuaModel) messageToLuaTable(msg tea.Msg) *lua.LTable {
	tbl := m.L.NewTable()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		tbl.RawSetString("type", lua.LString("key"))
		tbl.RawSetString("key", lua.LString(msg.String()))
		tbl.RawSetString("alt", lua.LBool(msg.Alt))
	case tea.WindowSizeMsg:
		tbl.RawSetString("type", lua.LString("resize"))
		tbl.RawSetString("width", lua.LNumber(msg.Width))
		tbl.RawSetString("height", lua.LNumber(msg.Height))
	default:
		tbl.RawSetString("type", lua.LString("unknown"))
	}
	return tbl
}

func (m *LuaModel) callLuaFunc(funcName string) {
	fn := m.L.GetGlobal(funcName)
	if fn.Type() != lua.LTFunction {
		panic("Lua function not defined: " + funcName)
	}

	msgTable := m.messageToLuaTable(m.msg)
	if err := m.L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    0,
		Protect: true,
	}, msgTable); err != nil {
		fmt.Fprintf(os.Stderr, "Lua error in '%s': %v\n", funcName, err)
	}
}

func (m *LuaModel) callLuaFuncWithReturn(funcName string, msg lua.LValue) lua.LValue {
	fn := m.L.GetGlobal(funcName)
	if fn.Type() != lua.LTFunction {
		return lua.LNil
	}

	if err := m.L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}, msg); err != nil {
		fmt.Fprintf(os.Stderr, "Lua error in '%s': %v\n", funcName, err)
		return lua.LNil
	}

	ret := m.L.Get(-1)
	m.L.Pop(1)
	return ret
}

func (m *LuaModel) Init() tea.Cmd {
	m.callLuaFunc("Init")
	return textinput.Blink
}

func (m *LuaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.msg = msg

	for _, comp := range m.components {
		if ti, ok := comp.(*TextInputComponent); ok {
			ti.Model, _ = ti.Model.Update(msg)
		}
	}

	ret := m.callLuaFuncWithReturn("Update", m.messageToLuaTable(msg))

	switch lv := ret.(type) {
	case lua.LString:
		switch string(lv) {
		case "quit":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *LuaModel) View() string {
	m.callLuaFunc("View")
	if m.ui.Rendered != "" {
		return m.ui.Rendered
	}
	var b strings.Builder
	for _, comp := range m.components {
		b.WriteString(comp.Render())
		b.WriteString("\n")
	}
	return b.String()
}

func Load(cfg config.Config, chars []characters.Character, script string) {
	L := lua.NewState()
	defer L.Close()

	ui := &LuaUI{}
	model := &LuaModel{
		L:            L,
		ui:           ui,
		componentMap: make(map[string]Component),
	}

	// L.SetContext(context.WithValue(context.Background(), "model", model))

	if err := L.DoFile(script); err != nil {
		panic(err)
	}

	exposeUI(L, model)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Println("error running app:", err)
	}
}

// func Load(cfg config.Config, chars []characters.Character, script string) {
// 	L := lua.NewState()
// 	defer L.Close()
//
// 	ui := &LuaUI{}
// 	ti := textinput.New()
// 	ti.Placeholder = "..."
//
// 	model := &LuaModel{
// 		L:         L,
// 		ui:        ui,
// 		textInput: ti,
// 	}
//
// 	// Load the script
// 	if err := L.DoFile(script); err != nil {
// 		panic(err)
// 	}
//
// 	// Now expose UI after loading the script so global `ui` is defined and synced
// 	exposeUI(L, ui)
//
// 	p := tea.NewProgram(model)
// 	if _, err := p.Run(); err != nil {
// 		fmt.Println("error running app:", err)
// 	}
// }
