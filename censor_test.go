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

	f("", "", false)
	f("Это чистый текст.", "Это чистый текст.", false)
}
