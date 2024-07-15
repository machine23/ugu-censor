package ugucensor

import (
	"strings"
	"unicode"
	"unicode/utf8"

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
	var (
		result          strings.Builder
		possibleBadPart strings.Builder
		word            strings.Builder
		rawWord         strings.Builder
		postWord        strings.Builder

		endOfWord    bool
		inWord       bool
		censored     bool
		hasBadPrefix bool
		hasBadPart   bool
		hasBadWord   bool
		isLetter     bool
		badWord      string
	)
	result.Grow(len(text))

	cursor := c.dicts[lang].Cursor()

	runes := []rune(text)
	lenRunes := len(runes)
	newWord := true

	stemmer := c.stemmers[lang]
	dict := c.dicts[lang]

	for i, ch := range runes {
		isLetter = unicode.IsLetter(ch)
		inWord = inWord || isLetter
		if !inWord {
			result.WriteRune(ch)
			continue
		}

		if isLetter {
			lowerCh := unicode.ToLower(ch)
			if postWord.Len() > 0 {
				rawWord.WriteString(postWord.String())
				postWord.Reset()
			}
			rawWord.WriteRune(ch)
			word.WriteRune(lowerCh)

			if newWord || hasBadPrefix {
				newWord = false

				hasBadPrefix, hasBadPart = cursor.Advance(lowerCh)
				if hasBadPrefix {
					inWord = true
					possibleBadPart.WriteRune(lowerCh)
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
			postWord.WriteRune(ch)
			endOfWord = !newWord && (!hasBadPrefix || hasBadWord)
		}

		endOfWord = endOfWord || i == lenRunes-1

		// write result

		if endOfWord {
			rawWordStr := rawWord.String()
			if hasBadWord {
				wordStr := word.String()
				isBadWord := badWord == wordStr
				if stemmer != nil && !isBadWord {
					isBadWord = dict.Search(stemmer.Stem(wordStr))
				}

				if isBadWord {
					result.WriteString(strings.Repeat("*", utf8.RuneCountInString(rawWordStr)))
					censored = true
				} else {
					result.WriteString(rawWordStr)
				}
			} else {
				result.WriteString(rawWordStr)
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
			badWord = ""
		}
	}
	return result.String(), censored
}
