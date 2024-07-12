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
	censored := false

	textRunes := []rune(text) // Convert the text to a slice of runes for proper Unicode handling
	wordStart := -1           // Use -1 to indicate that we're not currently tracking a word

	// Helper function to process and potentially censor a word
	processWord := func(wordEnd int) {
		if wordStart != -1 { // We have a word to process
			word := string(textRunes[wordStart:wordEnd])
			if c.isBadWord(word, langs...) {
				censored = true
				result.WriteString(strings.Repeat("*", wordEnd-wordStart))
			} else {
				result.WriteString(word)
			}
			wordStart = -1 // Reset wordStart for the next word
		}
	}

	for i, ch := range textRunes {
		if unicode.IsSpace(ch) || !unicode.IsLetter(ch) {
			processWord(i)       // Process the word ending at the current index
			result.WriteRune(ch) // Write the non-word character to the result
		} else if wordStart == -1 {
			wordStart = i // Start a new word
		}
	}

	processWord(len(textRunes)) // Process the last word, if any

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
