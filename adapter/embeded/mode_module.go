package embeded

import lua "github.com/yuin/gopher-lua"

func ModeLoader(l *lua.LState) int {
	// Create a new Lua table to hold the modes
	modes := l.NewTable()

	// Define the modes as constants (this is similar to an enum)
	l.SetField(modes, "SERVER", lua.LString("server"))
	l.SetField(modes, "RSERVER", lua.LString("ResourceServer"))
	l.SetField(modes, "CLIENT", lua.LString("client"))
	l.SetField(modes, "BCLIENT", lua.LString("BroadcastClient"))
	l.SetField(modes, "BSERVER", lua.LString("BroadcastServer"))
	l.SetField(modes, "ICLIENT", lua.LString("IntermediateClient"))
	l.SetField(modes, "ISERVER", lua.LString("IntermediateServer"))

	// Push the modes table onto the Lua stack and return it as a module
	l.Push(modes)
	return 1
}
