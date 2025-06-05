-- Store past entries
local entries = {}

local state = {}

function inspect(value, indent)
	indent = indent or 0
	local padding = string.rep("  ", indent)

	if type(value) == "table" then
		local parts = { "{" }
		for k, v in pairs(value) do
			local keyStr = tostring(k)
			local valStr = inspect(v, indent + 1)
			table.insert(parts, padding .. "  " .. keyStr .. " = " .. valStr)
		end
		table.insert(parts, padding .. "}")
		return table.concat(parts, "\n")
	elseif type(value) == "string" then
		return string.format("%q", value)
	else
		return tostring(value)
	end
end

function Init()
	local ti = ui.textInput.new()
	ti.placeholder = "n"
	state.textInput = ti
end

function Update(msg)
	if msg.type == "key" then
		if msg.key == "ctrl+c" then
			return "quit"
		elseif msg.key == "enter" then
			table.insert(entries, state.textInput.value or "")
			state.textInput.value = ""
		end
	end
	return "none"
end

function View()
	local lines = {
		"Type something and press Enter. Press Ctrl+C to quit.",
		"",
		"Input: " .. state.textInput.View(),
		"",
		"Previous entries:",
	}

	for _, line in ipairs(entries) do
		table.insert(lines, "- " .. line)
	end

	ui.render(lines)
end
