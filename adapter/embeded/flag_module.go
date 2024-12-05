package embeded

import lua "github.com/yuin/gopher-lua"

// Define the AdapterMode type and constants
type AdapterMode string

const (
	Client          AdapterMode = "Client"
	Server          AdapterMode = "Server"
	ResourceServer  AdapterMode = "ResourceServer"
	BroadCastClient AdapterMode = "BroadCastClient"
	BroadCastServer AdapterMode = "BroadCastServer"
	InterClient     AdapterMode = "InterClient"
	InterServer     AdapterMode = "InterServer"
)

type Flags struct {
	Mode       AdapterMode
	Enc        string
	Dec        string
	CryptoAlgo string
	KeyAlgo    string
}

var CurrentFlags Flags

func ModeLua(l *lua.LState) int {
	mode := l.CheckString(1)

	switch mode {
	case "Client":
		CurrentFlags.Mode = Client
	case "Server":
		CurrentFlags.Mode = Server
	case "ResourceServer":
		CurrentFlags.Mode = ResourceServer
	case "BroadCastClient":
		CurrentFlags.Mode = BroadCastClient
	case "BroadCastServer":
		CurrentFlags.Mode = BroadCastServer
	case "InterClient":
		CurrentFlags.Mode = InterClient
	case "InterServer":
		CurrentFlags.Mode = InterServer
	}
	l.Push(lua.LString(string(CurrentFlags.Mode)))
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
	exports := map[string]lua.LGFunction{
		"Mode":        ModeLua,
		"decryptPort": DecryptPortLua,
		"encryptPort": EncryptPortLua,
		"crypto":      CryptoLua,
		"keyExchange": KeyExchangeLua,
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}
