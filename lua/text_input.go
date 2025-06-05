package lua

import (
	"github.com/charmbracelet/bubbles/textinput"
	lua "github.com/yuin/gopher-lua"
)

func TI(L *lua.LState, m *LuaModel) lua.LValue {
	ti := L.NewTable()

	ti.RawSetString("new", L.NewFunction(func(L *lua.LState) int {
		// Create Go model
		tiModel := textinput.New()
		tiModel.Placeholder = "..."
		tiModel.Focus()

		comp := &TextInputComponent{Model: tiModel, Id: "textInput"}
		m.componentMap["textInput"] = comp
		m.components = []Component{comp}

		// Create Lua instance table
		inst := L.NewTable()
		mt := L.NewTypeMetatable("textInputMeta")

		L.SetField(mt, "__index", L.NewFunction(func(L *lua.LState) int {
			key := L.CheckString(2)

			switch key {
			case "value":
				if comp, ok := m.componentMap["textInput"].(*TextInputComponent); ok {
					L.Push(lua.LString(comp.Model.Value()))
					return 1
				}
			case "View":
				L.Push(L.NewFunction(func(L *lua.LState) int {
					if comp, ok := m.componentMap["textInput"].(*TextInputComponent); ok {
						L.Push(lua.LString(comp.Model.View()))
						return 1
					}
					return 0
				}))
				return 1
			case "placeholder":
				if comp, ok := m.componentMap["textInput"].(*TextInputComponent); ok {
					L.Push(lua.LString(comp.Model.Placeholder))
					return 1
				}
			}

			return 0
		}))

		L.SetField(mt, "__newindex", L.NewFunction(func(L *lua.LState) int {
			key := L.CheckString(2)
			val := L.Get(3)

			if comp, ok := m.componentMap["textInput"].(*TextInputComponent); ok {
				switch key {
				case "value":
					comp.Model.SetValue(val.String())
				case "placeholder":
					comp.Model.Placeholder = val.String()
				}
			}
			return 0
		}))

		L.SetMetatable(inst, mt)
		L.Push(inst)
		return 1
	}))

	return ti
}
