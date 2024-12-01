package embeded

// enum QuestionType {
// 	MCQ,
// 	From
// }

// type UIModule struct {
// 	Question string
// 	options  []string
// 	Questiontype QuestionType
// 	Answer string
// 	RenderFunc func() string
// }

// var UIModuleList []UIModule

// func registerUIModule(l *lua.LState) int {
// 	question := l.CheckString(1)
// 	questionType := l.CheckInt(2)
// 	options := l.CheckTable(3)
// 	answer := l.CheckString(4)

// 	var optionsList []string
// 	options.ForEach(func(_, value lua.LValue) {
// 		optionsList = append(optionsList, value.String())
// 	})

// 	UIModuleList = append(UIModuleList, UIModule{
// 		Question: question,
// 		options:  optionsList,
// 		Questiontype: QuestionType(questionType),
// 		Answer: answer,
// 	})

// 	return 0
// }

// func UIModuleLoader(l *lua.LState) int {
// 	var exports = map[string]lua.LGFunction{
// 		"add": registerUIModule,
// 	}

// 	mod := l.SetFuncs(l.NewTable(), exports)
// 	l.Push(mod)
// 	return 1
// }
