package embeded

import (
	"fmt"
	"net"

	lua "github.com/yuin/gopher-lua"
)

// KeyExchangeAlgo is a struct that holds the name of the key exchange algorithm and the function to generate a key
type KeyExchangeAlgo struct {
	Name     string
	Generate func(net.Conn) []byte // <- key -> generated key
	Key      []byte
}

// list of all registerd key exchange algo
var KeyExchangeList []KeyExchangeAlgo

// registerKeyExchange is a function that registers a key exchange algorithm in the lua state
func registerKeyExchange(l *lua.LState) int {
	algoName := l.CheckString(1)       // get first arg from function call in lua
	generateFunc := l.CheckFunction(2) // similarly second

	generate := func(conn net.Conn) []byte {
		luaConn := l.NewUserData()
		luaConn.Value = conn

		err := l.CallByParam(lua.P{
			Fn:      generateFunc,
			NRet:    1,
			Protect: true,
		}, luaConn)

		if err != nil {
			fmt.Println("Error in key generation:", err)
			return nil
		}

		luaResult := l.Get(-1)
		l.Pop(1)

		return []byte(luaResult.String())
	}

	KeyExchangeList = append(KeyExchangeList, KeyExchangeAlgo{
		Name:     algoName,
		Generate: generate,
		Key:      nil,
	})

	return 0
}

// YET TO BE IMPLEMENTED
func getKeyExchangeNames(l *lua.LState) int {
	tbl := l.NewTable()
	for i, algo := range KeyExchangeList {
		tbl.RawSetInt(i+1, lua.LString(algo.Name))
	}
	l.Push(tbl)
	return 1
}

// KeyExchangeLoader is a function that loads the key exchange module in the lua state
func KeyExchangeLoader(l *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"register": registerKeyExchange,
		"list":     getKeyExchangeNames,
	}

	mod := l.SetFuncs(l.NewTable(), exports)
	l.Push(mod)
	return 1
}
