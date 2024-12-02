package embeded

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type UIQuestion struct {
	// TODO ADD NAME:Name        string
	Question    string
	Options     []string
	PlaceHolder string
	Answer      string
	RenderFunc  func(question string, options []string, placeholder string) string
}

var UIQuestionList []UIQuestion

func registerUIQuestion(l *lua.LState) int {
	question := l.CheckString(1)
	options := l.CheckTable(2)
	placeholder := l.CheckString(3)
	answer := l.CheckString(4)
	renderFunc := l.CheckFunction(5)

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
		Question:    question,
		Options:     opts,
		PlaceHolder: placeholder,
		Answer:      answer,
		RenderFunc:  render,
	})

	return 0
}

// func (u UIQuestion) Render() string {
// 	return u.RenderFunc(u.Question, u.Options)
// }

// func (u UIQuestion) GetAnswer() string {
// 	return u.Answer
// }

func UIQuestionLoader(l *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"new": registerUIQuestion,
	}

	mod := l.SetFuncs(l.NewTable(), exports)
	l.Push(mod)
	return 1
}
