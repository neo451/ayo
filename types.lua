---@meta

---@class TextInput
---@field value string
---@field placeholder string
---@field focused boolean
---@field new fun(): TextInput
---@field View fun(): string

---@class LuaUI
---@field textInput TextInput
---@field show_message string
---@field render fun(lines: string[]): nil

---@type LuaUI
ui = {}

---@class tea.Msg
---@field type "key" | "win" TODO:
---@field key string
---@field height integer
---@field width integer

---@param msg table
Init = function(msg) end
