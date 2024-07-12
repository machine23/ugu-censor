package ugucensor

import (
	"strings"
	"unicode"

	"github.com/machine23/ugu-censor/trie"
)

type Censor struct {
	dicts map[string]*trie.Trie
}

func NewCensor() *Censor {
	return &Censor{
		dicts: make(map[string]*trie.Trie),
	}
}

func (c *Censor) AddWord(word string, lang string) {
	if _, ok := c.dicts[lang]; !ok {
		c.dicts[lang] = trie.NewTrie()
	}
	c.dicts[lang].Insert(word)
}

func (c *Censor) AddWords(words []string, lang string) {
	for _, word := range words {
		c.AddWord(word, lang)
	}
}

func (c *Censor) CensorText(text string, langs ...string) (string, bool) {
	var result strings.Builder
	var word strings.Builder
	var censored bool

	processWord := func() {
		if word.Len() > 0 {
			w := word.String()
			if c.isBadWord(w, langs...) {
				censored = true
				result.WriteString(strings.Repeat("*", len([]rune(w))))
			} else {
				result.WriteString(w)
			}
			word.Reset()
		}
	}

	for _, ch := range text {
		if unicode.IsSpace(ch) || !unicode.IsLetter(ch) {
			processWord()
			result.WriteRune(ch)
		} else {
			word.WriteRune(ch)
		}
	}

	processWord() // Process the last word if any

	return result.String(), censored
}

func (c *Censor) isBadWord(word string, langs ...string) bool {
	w := strings.ToLower(word)
	for _, lang := range langs {
		if dict, ok := c.dicts[lang]; ok {
			if dict.Search(w) {
				return true
			}
		}
	}
	return false
}
