package embeded

import lua "github.com/yuin/gopher-lua"

func HandleLua(l *lua.LState) {
	// define 'crypto' module in lua
	l.PreloadModule("crypto", CryptoLoader)
	// the now-standard algorithms will be provided by us.
	l.PreloadModule("algo.aes", AlogAesLoader)

	// Key exchange
	l.PreloadModule("keyExchange", KeyExchangeLoader)
	l.PreloadModule("keyalgo.dh", DiffieHellmanLoader)
	l.PreloadModule("keyalgo.rsa", RSAKeyExchangeLoader)

	if err := l.DoFile("config/init.lua"); err != nil {
		panic(err)
	}
}
