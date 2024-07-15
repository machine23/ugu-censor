package ugucensor

import (
	"strings"
	"unicode"

	"github.com/machine23/ugu-censor/trie"
	ugustemmer "github.com/machine23/ugu-stemmer"
)

type Stemmer interface {
	Stem(word string) string
}

type Censor struct {
	dicts    map[string]*trie.Trie
	stemmers map[string]Stemmer
}

func NewCensor() *Censor {
	return &Censor{
		dicts:    make(map[string]*trie.Trie),
		stemmers: make(map[string]Stemmer),
	}
}

func (c *Censor) AddWord(word string, lang string) {
	stemmer, ok := c.stemmers[lang]
	if !ok {
		stemmer = ugustemmer.NewSnowballStemmer(lang)
		c.stemmers[lang] = stemmer
	}
	if _, ok := c.dicts[lang]; !ok {
		c.dicts[lang] = trie.NewTrie()
	}
	if stemmer != nil {
		word = stemmer.Stem(word)
	}
	c.dicts[lang].Insert(word)
}

func (c *Censor) AddWords(words []string, lang string) {
	for _, word := range words {
		c.AddWord(word, lang)
	}
}

func (c *Censor) CensorText(text string, lang string) (string, bool) {
	var result strings.Builder
	var censored bool
	var possibleBadPart strings.Builder
	var word strings.Builder
	var rawWord strings.Builder
	var postWord strings.Builder
	var endOfWord bool
	var inWord bool

	var hasBadPrefix bool
	var hasBadPart bool
	var hasBadWord bool
	var badWord string

	cursor := c.dicts[lang].Cursor()

	runes := []rune(text)
	lenRunes := len(runes)
	newWord := true
	for i, ch := range runes {
		if unicode.IsLetter(ch) {
			inWord = inWord || true
		}
		if inWord {
			if unicode.IsLetter(ch) {
				if postWord.Len() > 0 {
					rawWord.WriteString(postWord.String())
					postWord.Reset()
				}
				rawWord.WriteRune(ch)
				word.WriteRune(ch)
			} else {
				postWord.WriteRune(ch)
			}
		} else {
			result.WriteRune(ch)
			continue
		}

		if unicode.IsLetter(ch) {
			if newWord || hasBadPrefix {
				newWord = false

				hasBadPrefix, hasBadPart = cursor.Advance(unicode.ToLower(ch))
				if hasBadPrefix {
					inWord = true
					possibleBadPart.WriteRune(unicode.ToLower(ch))
					if hasBadPart {
						hasBadWord = true
						badWord = possibleBadPart.String()
					}
					if i < lenRunes-1 {
						continue
					}
				}
			}
		} else {
			endOfWord = !newWord && (!hasBadPrefix || hasBadWord)
		}

		endOfWord = endOfWord || i == lenRunes-1

		// write result

		if endOfWord {
			if hasBadWord {
				stemmer := c.stemmers[lang]
				isBadWord := badWord == word.String()
				if stemmer != nil && !isBadWord {
					isBadWord = badWord == stemmer.Stem(word.String())
				}

				if isBadWord {
					result.WriteString(strings.Repeat("*", len([]rune(rawWord.String()))))
					censored = true
				} else {
					result.WriteString(rawWord.String())
				}
			} else {
				result.WriteString(rawWord.String())
			}
			if postWord.Len() > 0 {
				result.WriteString(postWord.String())
			}
			postWord.Reset()
			word.Reset()
			endOfWord = false
			possibleBadPart.Reset()
			cursor.Reset()
			inWord = false
			newWord = true
			hasBadWord = false
			rawWord.Reset()
		}
	}
	return result.String(), censored
}
