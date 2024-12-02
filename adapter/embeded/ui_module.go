package embeded

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type UIQuestion struct {
	Name        string
	Question    string
	Options     []string
	PlaceHolder string
	Answer      string
	RenderFunc  func(question string, options []string, placeholder string) string
}

var UIQuestionList []UIQuestion

func registerUIQuestion(l *lua.LState) int {
	name := l.CheckString(1)
	question := l.CheckString(2)
	options := l.CheckTable(3)
	placeholder := l.CheckString(4)
	answer := l.CheckString(5)
	renderFunc := l.CheckFunction(6)

	var opts []string
	options.ForEach(func(_, value lua.LValue) {
		opts = append(opts, value.String())
	})

	render := func(q string, o []string, pl string) string {
		luaQuestion := lua.LString(q)
		luaOptions := l.NewTable()
		for i, opt := range o {
			l.SetTable(luaOptions, lua.LNumber(i+1), lua.LString(opt))
		}
		luaPlaceHolder := lua.LString(pl)

		err := l.CallByParam(lua.P{
			Fn:      renderFunc,
			NRet:    1,
			Protect: true,
		}, luaQuestion, luaOptions, luaPlaceHolder)

		if err != nil {
			fmt.Println("Error in rendering question:", err)
			return ""
		}

		luaResult := l.Get(-1)
		l.Pop(1)

		return luaResult.String()
	}

	UIQuestionList = append(UIQuestionList, UIQuestion{
		Name:        name,
		Question:    question,
		Options:     opts,
		PlaceHolder: placeholder,
		Answer:      answer,
		RenderFunc:  render,
	})

	return 0
}

func UILoader(l *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"new": registerUIQuestion,
	}

	mod := l.SetFuncs(l.NewTable(), exports)
	l.Push(mod)
	return 1
}
