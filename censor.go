package ugucensor

import (
	"strings"

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
	if c.isBadWord(text, langs...) {
		return strings.Repeat("*", len([]rune(text))), true
	}

	return text, false
}

func (c *Censor) isBadWord(word string, langs ...string) bool {
	for _, lang := range langs {
		if dict, ok := c.dicts[lang]; ok {
			if dict.Search(word) {
				return true
			}
		}
	}
	return false
}
