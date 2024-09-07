package embeded

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type CryptoAlgo struct {
	Name    string
	Encrypt func([]byte, []byte) // <- key, data
	Decrypt func([]byte, []byte) // <- key, data
}

// list of all registerd cryto algo
var CryptoList []CryptoAlgo

func registerCrypto(l *lua.LState) int {
	algoName := l.CheckString(1)      // get first arg from function call in lua
	encryptFunc := l.CheckFunction(2) // similarly second
	decryptFunc := l.CheckFunction(3) // and third

	encrypt := func(key, data []byte) {
		luaKey := lua.LString(string(key))
		luaData := lua.LString(string(data))

		err := l.CallByParam(lua.P{
			Fn:      encryptFunc,
			NRet:    0,
			Protect: true,
		}, luaKey, luaData)

		if err != nil {
			fmt.Println("Error in encryption:", err)
		}
	}

	decrypt := func(key, data []byte) {
		luaKey := lua.LString(string(key))
		luaData := lua.LString(string(data))

		err := l.CallByParam(lua.P{
			Fn:      decryptFunc,
			NRet:    0,
			Protect: true,
		}, luaKey, luaData)

		if err != nil {
			fmt.Println("Error in decryption:", err)
		}
	}

	CryptoList = append(CryptoList, CryptoAlgo{
		Name:    algoName,
		Encrypt: encrypt,
		Decrypt: decrypt,
	})

	return 0
}

func CryptoLoader(l *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"register": registerCrypto,
	}

	mod := l.SetFuncs(l.NewTable(), exports)
	l.Push(mod)
	return 1
}
