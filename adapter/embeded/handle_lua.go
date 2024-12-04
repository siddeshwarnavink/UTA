package embeded

import lua "github.com/yuin/gopher-lua"

func HandleLua(l *lua.LState, configPath string) {
	l.PreloadModule("mode", ModeLoader)
	l.PreloadModule("config", ConfigLoader)

	// define 'crypto' module in lua
	l.PreloadModule("crypto", CryptoLoader)
	// the now-standard algorithms will be provided by us.
	l.PreloadModule("algo.aes", AlogAesLoader)

	// Key exchange
	l.PreloadModule("keyExchange", KeyExchangeLoader)
	l.PreloadModule("keyalgo.dh", DiffieHellmanLoader)
	l.PreloadModule("keyalgo.rsa", RSAKeyExchangeLoader)

	l.PreloadModule("ui", UILoader)
	l.PreloadModule("ui.form", UIFormLoader)
	l.PreloadModule("ui.mcq", UIMCQLoader)

	if err := l.DoFile(configPath); err != nil {
		panic(err)
	}
}
