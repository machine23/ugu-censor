package ugucensor

import "testing"

func TestCensor_CensorText(t *testing.T) {
	c := NewCensor()

	words := []string{"игра", "яблоко"}
	c.AddWords(words, "ru")

	f := func(text string, expected string, expectedCensored bool) {
		t.Helper()

		got, gotCensored := c.CensorText(text, "ru")
		if got != expected {
			t.Errorf("\nCensorText(%q, \"ru\")\n\tgot : %s\n\twant: %s", text, got, expected)
		}

		if gotCensored != expectedCensored {
			t.Errorf("\nCensorText(%q, \"ru\")\n\tisCensored\n\tgot : %v\n\twant: %v", text, gotCensored, expectedCensored)
		}
	}

	t.Run("empty text", func(t *testing.T) {
		f("", "", false)
	})

	t.Run("clean text", func(t *testing.T) {
		f("Это чистый текст.", "Это чистый текст.", false)
	})

	t.Run("single word", func(t *testing.T) {
		f("игра", "****", true)
		f("яблоко", "******", true)
	})

	t.Run("case insensitive one word", func(t *testing.T) {
		f("Игра", "****", true)
		f("ИГРА", "****", true)
		f("ЯбЛоКо", "******", true)
	})

	t.Run("multiple words", func(t *testing.T) {
		f("это игра", "это ****", true)
		f("игра это", "**** это", true)
		f("игра игра", "**** ****", true)
		f("игра яблоко", "**** ******", true)
		f("игра яблоко игра", "**** ****** ****", true)
		f("Эта игра хорошая", "Эта **** хорошая", true)
	})

	t.Run("standart punctuation", func(t *testing.T) {
		f("игра.", "****.", true)
		f("игра!", "****!", true)
		f("игра?", "****?", true)
		f("игра...", "****...", true)
		f("игра!?", "****!?", true)
		f("игра, игра", "****, ****", true)
		f("игра,игра", "****,****", true)
		f("игра, игра.", "****, ****.", true)
		f("Это та самая игра!", "Это та самая ****!", true)
		f("Это та самая игра?", "Это та самая ****?", true)
		f("Игра, та самая игра.", "****, та самая ****.", true)
		f("Игра, а потом яблоко.", "****, а потом ******.", true)
	})
}
