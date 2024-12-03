package embeded

import lua "github.com/yuin/gopher-lua"

type AdapterMode string

const (
	Client AdapterMode = "Client"
	Server AdapterMode = "Server"
)

type Flags struct {
	Mode       AdapterMode
	Enc        string
	Dec        string
	CryptoAlgo string
	KeyAlgo    string
}

var currentFlags Flags

func ServerModeLua(L *lua.LState) int {
	isServer := L.CheckBool(1)
	if isServer {
		currentFlags.Mode = Server
	} else {
		currentFlags.Mode = Client
	}
	L.Push(lua.LString(string(currentFlags.Mode)))
	return 1
}

func DecryptPortLua(L *lua.LState) int {
	currentFlags.Dec = L.CheckString(1)
	L.Push(lua.LString(currentFlags.Dec))
	return 1
}

func EncryptPortLua(L *lua.LState) int {
	currentFlags.Enc = L.CheckString(1)
	L.Push(lua.LString(currentFlags.Enc))
	return 1
}

func CryptoLua(L *lua.LState) int {
	currentFlags.CryptoAlgo = L.CheckString(1)
	L.Push(lua.LString(currentFlags.CryptoAlgo))
	return 1
}

func KeyExchangeLua(L *lua.LState) int {
	currentFlags.KeyAlgo = L.CheckString(1)
	L.Push(lua.LString(currentFlags.KeyAlgo))
	return 1
}

func ConfigLoader(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"serverMode":  ServerModeLua,
		"decryptPort": DecryptPortLua,
		"encryptPort": EncryptPortLua,
		"crypto":      CryptoLua,
		"keyExchange": KeyExchangeLua,
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}
