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

var CurrentFlags Flags

func ServerModeLua(L *lua.LState) int {
	isServer := L.CheckBool(1)
	if isServer {
		CurrentFlags.Mode = Server
	} else {
		CurrentFlags.Mode = Client
	}
	L.Push(lua.LString(string(CurrentFlags.Mode)))
	return 1
}

func DecryptPortLua(L *lua.LState) int {
	CurrentFlags.Dec = L.CheckString(1)
	L.Push(lua.LString(CurrentFlags.Dec))
	return 1
}

func EncryptPortLua(L *lua.LState) int {
	CurrentFlags.Enc = L.CheckString(1)
	L.Push(lua.LString(CurrentFlags.Enc))
	return 1
}

func CryptoLua(L *lua.LState) int {
	CurrentFlags.CryptoAlgo = L.CheckString(1)
	L.Push(lua.LString(CurrentFlags.CryptoAlgo))
	return 1
}

func KeyExchangeLua(L *lua.LState) int {
	CurrentFlags.KeyAlgo = L.CheckString(1)
	L.Push(lua.LString(CurrentFlags.KeyAlgo))
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
