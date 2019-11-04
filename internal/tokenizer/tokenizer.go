package tokenizer

import kagomeTok "github.com/ikawaha/kagome/tokenizer"

type Tokenizer struct {
	kagome kagomeTok.Tokenizer
}

func New() Tokenizer {
	dic := kagomeTok.SysDicIPA()
	udic, err := kagomeTok.NewUserDic("/home/osamu/data/kuromoji-dict/mydict.csv")
	if err != nil {
		panic(err)
	}
	kagome := kagomeTok.NewWithDic(dic)
	kagome.SetUserDic(udic)
	return Tokenizer{
		kagome: kagome,
	}
}

func (t Tokenizer) Tokenize(text string) []string {
	tokens := t.kagome.Analyze(text, kagomeTok.Normal)
	result := make([]string, len(tokens))
	for i := 0; i < len(tokens); i++ {
		result[i] = tokens[i].Surface
	}
	return result
}