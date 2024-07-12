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
	var cleanWord strings.Builder
	var postPrefix strings.Builder
	var hasPrefix, isBad bool
	censored := false

	runes := []rune(text)
	for i, ch := range runes {
		isLetter := unicode.IsLetter(ch)

		if isLetter {
			if postPrefix.Len() > 0 && hasPrefix {
				word.WriteString(postPrefix.String())
				postPrefix.Reset()
			} else if postPrefix.Len() > 0 {
				result.WriteString(postPrefix.String())
				postPrefix.Reset()
			}
			word.WriteRune(ch)
			cleanWord.WriteRune(unicode.ToLower(ch))
		} else {
			postPrefix.WriteRune(ch)
		}

		if !isLetter || i == len(runes)-1 {
			cleanWordStr := cleanWord.String()
			hasPrefix, isBad = c.isBadWord(cleanWordStr, langs...)
			if hasPrefix && !isBad {
				continue
			}

			if isBad {
				censored = true
				result.WriteString(strings.Repeat("*", len([]rune(word.String()))))
				result.WriteString(postPrefix.String())
				postPrefix.Reset()
				word.Reset()
				cleanWord.Reset()
				continue
			}

			if !hasPrefix || i == len(runes)-1 {

				result.WriteString(word.String())
				result.WriteString(postPrefix.String())
				postPrefix.Reset()
				word.Reset()
				cleanWord.Reset()
			}
		}
	}

	if postPrefix.Len() > 0 {
		result.WriteString(postPrefix.String())
	}

	return result.String(), censored
}

func (c *Censor) isBadWord(word string, langs ...string) (bool, bool) {
	word = strings.ToLower(word)
	var hasPrefix, isCompleteWord bool
	for _, lang := range langs {
		if dict, ok := c.dicts[lang]; ok {
			hasPrefix, isCompleteWord = dict.StartsWith(word)
			if isCompleteWord {
				return true, true
			}
		}
	}

	return hasPrefix, false
}

func (c Censor) charFilter(ch rune) rune {
	if unicode.IsLetter(ch) {
		return ch
	}
	return -1
}
