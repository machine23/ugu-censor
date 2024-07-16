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
	return c.twoPassCensorText(text, lang)
}

func (c *Censor) onePassCensorText(text string, lang string) (string, bool) {
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

type PossibleBadWordBounds struct {
	BadPart string
	Word    string
	Start   int
	End     int
}

func (c *Censor) twoPassCensorText(text string, lang string) (string, bool) {
	var (
		result   strings.Builder
		censored bool

		possibleBadWordStarts []int
	)

	result.Grow(len(text))

	stemmer := c.stemmers[lang]

	runes := []rune(text)

	// first pass
	// find all possible bad word starts

	possibleBadWordStarts = c.findPossibleBadWordStarts(runes, lang)

	// second pass
	// check all possible bad word starts and censor them

	wordBounds := c.findPossibleBadWordBounds(runes, possibleBadWordStarts, lang)

	if len(wordBounds) == 0 {
		return text, false
	}

	// check all possible bad words, censor them and write result

	for i, wb := range wordBounds {
		// Determine the previous end index or start from 0
		prevEnd := 0
		if i > 0 {
			prevEnd = wordBounds[i-1].End
		}
		// Append the text before the current word
		result.WriteString(string(runes[prevEnd:wb.Start]))

		// Check if the word is a bad word, either directly or via the stemmer
		isBadWord := wb.Word == wb.BadPart ||
			(stemmer != nil && c.dicts[lang].Search(stemmer.Stem(wb.Word)))

		if isBadWord {
			// Replace the bad word with asterisks
			result.WriteString(strings.Repeat("*", wb.End-wb.Start))
			censored = true
			continue
		}

		// Append the current word as it is not censored
		result.WriteString(string(runes[wb.Start:wb.End]))
	}

	// write the rest of the text
	result.WriteString(string(runes[wordBounds[len(wordBounds)-1].End:]))

	return result.String(), censored
}

func (c *Censor) findPossibleBadWordStarts(runes []rune, lang string) []int {
	var (
		possibleBadWordStarts []int
		newWord               bool

		cursor   = c.dicts[lang].Cursor()
		lenRunes = len(runes)
	)

	for i, ch := range runes {
		isLetter := unicode.IsLetter(ch)
		if isLetter {
			newWord = i == 0 || !unicode.IsLetter(runes[i-1])
		}

		if newWord {
			cursor.Reset()

			// check first letter
			bp, _ := cursor.Advance(unicode.ToLower(ch))
			if !bp {
				continue
			}

			// find and check second letter, skip all non-letters
			// if there is no second letter, then it's not a bad word
			// if second letter is prefix of bad word, then add i index
			// to possibleBadWordStarts
			for j := i + 1; j < lenRunes; j++ {
				nextCh := runes[j]
				if unicode.IsLetter(nextCh) {
					possibleBadWordStart, _ := cursor.Advance(unicode.ToLower(nextCh))
					if possibleBadWordStart {
						possibleBadWordStarts = append(possibleBadWordStarts, i)
					}
					break
				}
			}
		}
	}
	return possibleBadWordStarts
}

func (c *Censor) findPossibleBadWordBounds(runes []rune, starts []int, lang string) []PossibleBadWordBounds {
	var (
		badWords []PossibleBadWordBounds
		cursor   = c.dicts[lang].Cursor()
		lenRunes = len(runes)
	)

	var badPart strings.Builder
	var badWord strings.Builder
	for _, bwStart := range starts {
		cursor.Reset()
		var hasBadWord bool
		var hasBadPrefix, hasBadPart bool

		badPart.Reset()
		badWord.Reset()

		badWordBounds := PossibleBadWordBounds{}
		for i := bwStart; i < lenRunes; i++ {
			ch := unicode.ToLower(runes[i])
			if unicode.IsLetter(ch) {
				hasBadPrefix, hasBadPart = cursor.Advance(ch)
				if !hasBadPrefix && !hasBadWord {
					break
				}

				if hasBadPrefix {
					badPart.WriteRune(ch)
				}

				badWord.WriteRune(ch)

				if !hasBadWord && hasBadPart {
					hasBadWord = true
					badWordBounds.BadPart = badPart.String()
				}

				if i == lenRunes-1 && hasBadWord {
					badWordBounds.Start = bwStart
					badWordBounds.End = i + 1
					badWordBounds.Word = badWord.String()
					badWords = append(badWords, badWordBounds)
				}
			} else {
				if hasBadWord {
					badWordBounds.Start = bwStart
					badWordBounds.End = i
					badWordBounds.Word = badWord.String()
					badWords = append(badWords, badWordBounds)
					break
				}
			}

		}

	}
	return badWords
}
