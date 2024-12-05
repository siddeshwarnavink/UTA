package embeded

import lua "github.com/yuin/gopher-lua"

func HandleLua(l *lua.LState, configPath string) {
	l.PreloadModule("config", ConfigLoader)
	l.PreloadModule("crypto", CryptoLoader)
	l.PreloadModule("algo.aes", AlogAesLoader)
	l.PreloadModule("keyExchange", KeyExchangeLoader)
	l.PreloadModule("keyalgo.dh", DiffieHellmanLoader)
	l.PreloadModule("keyalgo.rsa", RSAKeyExchangeLoader)

	if err := l.DoFile(configPath); err != nil {
		panic(err)
	}
}
