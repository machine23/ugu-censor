package ugucensor

import "testing"

func TestCensor_CensorText(t *testing.T) {
	c := NewCensor()

	words := []string{"игра", "игрок", "играть", "яблоко"}
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
		f("лучшая игра", "лучшая ****", true)
		f("игра лучшая", "**** лучшая", true)
	})

	t.Run("multiple words 2", func(t *testing.T) {
		f("яблоко и игра", "****** и ****", true)
		f("я яблоко", "я ******", true)
		f("я б яблоко", "я б ******", true)
		f("я я яблоко", "я я ******", true)
		f("я я я яблоко", "я я я ******", true)
		f("я я я я яблоко", "я я я я ******", true)
		f("\nЯблоки и игры с ними прочно вошли в нашу культуру,", "\n****** и **** с ними прочно вошли в нашу культуру,", true)
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

	t.Run("non-standart punctuation", func(t *testing.T) {
		// TODO: full word censoring
		f("и.г.р.а...", "*****.а...", true)
		f("_И_Г_Р_А_", "_*****_А_", true)
		f("Это иг....ра!", "Это ********!", true)
		f("Эта и..гр...а лучшая", "Эта *****...а лучшая", true)
		f("И.Г.Р.А, а я*бл*о*к*о потом.", "*****.А, а *********о потом.", true)
		f("Это та самая и г р а!", "Это та самая ***** а!", true)
	})

	t.Run("mixed punctuation", func(t *testing.T) {
		f("самая и г р а как игр, ат", "самая ***** а как ***, ат", true)
	})

	t.Run("no false positives", func(t *testing.T) {
		f("играция", "играция", false)
		f("и грация", "и грация", false)
		f("подвиг радость", "подвиг радость", false)
	})

	t.Run("multiple forms of the same word", func(t *testing.T) {
		f("Нет игры без правил.", "Нет **** без правил.", true)
		f("За игрой следует игра.", "За ***** следует ****.", true)
		f("Все заваляно играми.", "Все заваляно ******.", true)
		f("Завтра яблоко, а сегодня яблоки.", "Завтра ******, а сегодня ******.", true)
		f("Яблоку быть яблоком.", "****** быть *******.", true)
		f("Торгует яблоками.", "Торгует ********.", true)
		f("в этих играх самое главное", "в этих ***** самое главное", true)
		f("В таких играх игроки готовят блюда из яблок.", "В таких ***** ****** готовят блюда из *****.", true)
		f("\nВ современных настольных играх ****** тоже находят своё место. ", "\nВ современных настольных ***** ****** тоже находят своё место. ", true)
	})

	t.Run("long text", func(t *testing.T) {
		text := `Не забываем и о том, что яблоки часто используются в различных кулинарных играх. В таких играх игроки готовят блюда из яблок, следуя рецептам и выполняя различные кулинарные задания. Это могут быть пироги, салаты, соки и множество других вкусных и полезных блюд. Такие игры помогают игрокам узнать больше о кулинарии, развить свои кулинарные навыки и, возможно, даже вдохновить на создание собственных шедевров на кухне. Яблоки и игры с ними прочно вошли в нашу культуру, становясь символом радости, веселья и творчества. Независимо от того, играете ли вы в старинные игры с яблоками на свежем воздухе, наслаждаетесь настольными играми в кругу семьи или погружаетесь в виртуальные миры компьютерных игр, яблоки всегда добавляют веселья и радости в процесс. Эти фрукты стали неотъемлемой частью нашей жизни, объединяя людей в стремлении к игре, развлечению и новым открытиям.`

		expected := `Не забываем и о том, что ****** часто используются в различных кулинарных *****. В таких ***** ****** готовят блюда из *****, следуя рецептам и выполняя различные кулинарные задания. Это могут быть пироги, салаты, соки и множество других вкусных и полезных блюд. Такие **** помогают ******* узнать больше о кулинарии, развить свои кулинарные навыки и, возможно, даже вдохновить на создание собственных шедевров на кухне. ****** и **** с ними прочно вошли в нашу культуру, становясь символом радости, веселья и творчества. Независимо от того, ******* ли вы в старинные **** с ******** на свежем воздухе, наслаждаетесь настольными ****** в кругу семьи или погружаетесь в виртуальные миры компьютерных ***, ****** всегда добавляют веселья и радости в процесс. Эти фрукты стали неотъемлемой частью нашей жизни, объединяя людей в стремлении к ****, развлечению и новым открытиям.`

		f(text, expected, true)
	})
}

func BenchmarkCensorText(b *testing.B) {
	c := NewCensor()

	words := []string{"игра", "яблоко"}
	c.AddWords(words, "ru")

	text := "Это та самая игра, которую я видел вчера. Яблоко было вкусным. Хотя после игры яблоко уже не казалось таким вкусным."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.CensorText(text, "ru")
	}
}

func TestCensor_findPossibleBadWordStarts(t *testing.T) {
	c := NewCensor()

	words := []string{"игра", "игрок", "играть", "яблоко"}
	c.AddWords(words, "ru")

	f := func(text string, expected []int) {
		t.Helper()

		got := c.findPossibleBadWordStarts([]rune(text), "ru")
		if len(got) != len(expected) {
			t.Errorf("\nfindPossibleBadWordStarts(%q, \"ru\")\n\tgot : %v\n\twant: %v", text, got, expected)
			return
		}

		for i := range got {
			if got[i] != expected[i] {
				t.Errorf("\nfindPossibleBadWordStarts(%q, \"ru\")\n\tgot : %v\n\twant: %v", text, got, expected)
				return
			}
		}
	}

	t.Run("empty text", func(t *testing.T) {
		f("", []int{})
	})

	t.Run("clean text", func(t *testing.T) {
		f("Это чистый текст.", []int{})
	})

	t.Run("single word", func(t *testing.T) {
		f("игра", []int{0})
		f("яблоко", []int{0})
	})

	t.Run("multiple words", func(t *testing.T) {
		f("это игра", []int{4})
		f("игра это", []int{0})
		f("игра игра", []int{0, 5})
		f("игра яблоко", []int{0, 5})
		f("игра яблоко игра", []int{0, 5, 12})
		f("Эта игра хорошая", []int{4})
		f("лучшая игра", []int{7})
		f("игра лучшая", []int{0})
	})

	t.Run("multiple words 2", func(t *testing.T) {
		f("яблоко и игра", []int{0, 9})
		f("я яблоко", []int{2})
		f("я б яблоко", []int{0, 4})
		f("я я яблоко", []int{4})
		f("я я я яблоко", []int{6})
		f("игра я", []int{0})
		f("игра я б", []int{0, 5})
	})

	t.Run("case insensitive", func(t *testing.T) {
		f("Игра", []int{0})
		f("ИГРА", []int{0})
		f("ЯбЛоКо", []int{0})
		f("ЯБЛОКО И ИГРА", []int{0, 9})
		f("яблоко И игра", []int{0, 9})
	})

	t.Run("with symbols", func(t *testing.T) {
		f("игра. яблоко", []int{0, 6})
		f("игра.яблоко", []int{0, 5})
		f("  игра. яблоко", []int{2, 8})
		f("игра. яблоко  ", []int{0, 6})
		f("игра.>яблоко.", []int{0, 6})
		f("***и*г*р*а* яблоко", []int{3, 12})
		f("игра *1*а*я*я******  ...****б*л*о*к*о", []int{0, 12})
	})
}

func TestCensor_findPossibleBadWordBounds(t *testing.T) {
	c := NewCensor()

	words := []string{"игра", "игрок", "играть", "яблоко"}
	c.AddWords(words, "ru")

	f := func(text string, starts []int, expected []PossibleBadWordBounds) {
		t.Helper()

		got := c.findPossibleBadWordBounds([]rune(text), starts, "ru")
		if len(got) != len(expected) {
			t.Errorf("\nfindPossibleBadWordBounds(%q, %v, \"ru\")\n\tgot : %v\n\twant: %v", text, starts, got, expected)
			return
		}

		for i := range got {
			if got[i] != expected[i] {
				t.Errorf("\nfindPossibleBadWordBounds(%q, %v, \"ru\")\n\tgot : %v\n\twant: %v", text, starts, got, expected)
				return
			}
		}
	}

	t.Run("empty text", func(t *testing.T) {
		f("", []int{}, []PossibleBadWordBounds{})
	})

	t.Run("clean text", func(t *testing.T) {
		f("Это чистый текст.", []int{}, []PossibleBadWordBounds{})
	})

	t.Run("single word", func(t *testing.T) {
		f("игра", []int{0}, []PossibleBadWordBounds{{"игр", "игра", 0, 4}})
		f("яблоко", []int{0}, []PossibleBadWordBounds{{"яблок", "яблоко", 0, 6}})
	})

	t.Run("multiple words", func(t *testing.T) {
		f("это игра", []int{4}, []PossibleBadWordBounds{{"игр", "игра", 4, 8}})
		f("игра это", []int{0}, []PossibleBadWordBounds{{"игр", "игра", 0, 4}})
		f("игра игра", []int{0, 5}, []PossibleBadWordBounds{{"игр", "игра", 0, 4}, {"игр", "игра", 5, 9}})
		f("игра яблоко", []int{0, 5}, []PossibleBadWordBounds{{"игр", "игра", 0, 4}, {"яблок", "яблоко", 5, 11}})
		f("яблоко и игра", []int{0, 9}, []PossibleBadWordBounds{{"яблок", "яблоко", 0, 6}, {"игр", "игра", 9, 13}})
		f("игра яблоко игра", []int{0, 5, 12}, []PossibleBadWordBounds{{"игр", "игра", 0, 4}, {"яблок", "яблоко", 5, 11}, {"игр", "игра", 12, 16}})
	})

	t.Run("with symbols", func(t *testing.T) {
		f("игра. яблоко", []int{0, 6}, []PossibleBadWordBounds{{"игр", "игра", 0, 4}, {"яблок", "яблоко", 6, 12}})
		f("игра.яблоко", []int{0, 5}, []PossibleBadWordBounds{{"игр", "игра", 0, 4}, {"яблок", "яблоко", 5, 11}})
		f("  игра. яблоко", []int{2, 8}, []PossibleBadWordBounds{{"игр", "игра", 2, 6}, {"яблок", "яблоко", 8, 14}})
		f("игра. яблоко  ", []int{0, 6}, []PossibleBadWordBounds{{"игр", "игра", 0, 4}, {"яблок", "яблоко", 6, 12}})
		f("игра.>яблоко.", []int{0, 6}, []PossibleBadWordBounds{{"игр", "игра", 0, 4}, {"яблок", "яблоко", 6, 12}})
		f("***и*г*р*а* яблоко", []int{3, 12}, []PossibleBadWordBounds{{"игр", "игр", 3, 8}, {"яблок", "яблоко", 12, 18}})
		f("игра *1*а*я*я******  ...****б*л*о*к*о", []int{0, 12}, []PossibleBadWordBounds{{"игр", "игра", 0, 4}, {"яблок", "яблок", 12, 35}})
	})
}
