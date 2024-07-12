package ugucensor

type Censor struct{}

func NewCensor() *Censor {
	return &Censor{}
}

func (c *Censor) AddWord(word string, lang string)     {}
func (c *Censor) AddWords(words []string, lang string) {}

func (c *Censor) CensorText(text string, langs ...string) (string, bool) {
	return text, false
}
