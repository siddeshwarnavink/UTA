package embeded

import lua "github.com/yuin/gopher-lua"

func HandleLua(l *lua.LState) {
	// define 'crypto' module in lua
	l.PreloadModule("crypto", CryptoLoader)

	if err := l.DoFile("config/init.lua"); err != nil {
		panic(err)
	}
}
