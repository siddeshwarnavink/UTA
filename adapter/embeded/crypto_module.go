package embeded

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type CryptoAlgo struct {
	Name    string
	Encrypt func([]byte, []byte) []byte // <- key, data -> encrypted data
	Decrypt func([]byte, []byte) []byte // <- key, data -> decrypted data
}

// list of all registerd cryto algo
var CryptoList []CryptoAlgo

func registerCrypto(l *lua.LState) int {
	algoName := l.CheckString(1)      // get first arg from function call in lua
	encryptFunc := l.CheckFunction(2) // similarly second
	decryptFunc := l.CheckFunction(3) // and third

	encrypt := func(key, data []byte) []byte {
		luaKey := lua.LString(string(key))
		luaData := lua.LString(string(data))

		err := l.CallByParam(lua.P{
			Fn:      encryptFunc,
			NRet:    1,
			Protect: true,
		}, luaKey, luaData)

		if err != nil {
			fmt.Println("Error in encryption:", err)
			return nil
		}

		luaResult := l.Get(-1)
		l.Pop(1)

		return []byte(luaResult.String())
	}

	decrypt := func(key, data []byte) []byte {
		luaKey := lua.LString(string(key))
		luaData := lua.LString(string(data))

		err := l.CallByParam(lua.P{
			Fn:      decryptFunc,
			NRet:    1,
			Protect: true,
		}, luaKey, luaData)

		if err != nil {
			fmt.Println("Error in decryption:", err)
			return nil
		}

		luaResult := l.Get(-1)
		l.Pop(1)

		return []byte(luaResult.String())
	}

	CryptoList = append(CryptoList, CryptoAlgo{
		Name:    algoName,
		Encrypt: encrypt,
		Decrypt: decrypt,
	})

	return 0
}

// YET TO BE IMPLEMENTED
func getCryptoNames(l *lua.LState) int {
	tbl := l.NewTable()
	for i, algo := range CryptoList {
		tbl.RawSetInt(i+1, lua.LString(algo.Name))
	}
	l.Push(tbl)
	return 1
}

func CryptoLoader(l *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"register": registerCrypto,
		"list":     getCryptoNames,
	}

	mod := l.SetFuncs(l.NewTable(), exports)
	l.Push(mod)
	return 1
}
