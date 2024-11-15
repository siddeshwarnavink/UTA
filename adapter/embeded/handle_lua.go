package embeded

import lua "github.com/yuin/gopher-lua"

func HandleLua(l *lua.LState) {
	// define 'crypto' module in lua
	l.PreloadModule("crypto", CryptoLoader)
	l.PreloadModule("algo.aes", AlogAesLoader)

	if err := l.DoFile("adapter/config/init.lua"); err != nil {
		panic(err)
	}
}
