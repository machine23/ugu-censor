package trie

import (
	"testing"
)

func TestTrie_Search(t *testing.T) {
	trie := NewTrie()

	wordsToInsert := []string{"apple", "app", "ap", "banana", "band", "bandit"}
	for _, word := range wordsToInsert {
		trie.Insert(word)
	}

	f := func(word string, expected bool) {
		t.Helper()

		if got := trie.Search(word); got != expected {
			t.Errorf("Search(%q) = %v; want %v", word, got, expected)
		}
	}

	f("apple", true)
	f("app", true)
	f("ap", true)
	f("banana", true)
	f("band", true)
	f("bandit", true)
	f("ban", false)
	f("bandi", false)
	f("bandits", false)
	f("applea", false)
	f("apples", false)
	f("applesauce", false)
	f("api", false)
	f("a", false)
	f("", false)
}

func TestTrie_StartsWith(t *testing.T) {
	trie := NewTrie()

	wordsToInsert := []string{"apple", "app", "ap", "banana", "band", "bandit"}
	for _, word := range wordsToInsert {
		trie.Insert(word)
	}

	f := func(prefix string, expected bool, isEnd bool) {
		t.Helper()

		if got, gotEnd := trie.StartsWith(prefix); got != expected || gotEnd != isEnd {
			t.Errorf("StartsWith(%q) = %v, %v; want %v, %v", prefix, got, gotEnd, expected, isEnd)
		}
	}

	f("apple", true, true)
	f("app", true, true)
	f("ap", true, true)
	f("banana", true, true)
	f("band", true, true)
	f("bandit", true, true)
	f("appl", true, false)
	f("ban", true, false)
	f("bandi", true, false)
	f("a", true, false)
	f("b", true, false)
	f("c", false, false)
	f("d", false, false)
	f("", false, false)
	f("bandits", false, false)
	f("applea", false, false)
	f("apples", false, false)
	f("applesauce", false, false)
}

func TestTrie_Remove(t *testing.T) {
	trie := NewTrie()

	wordsToInsert := []string{"apple", "app", "ap", "banana", "band", "bandit"}
	for _, word := range wordsToInsert {
		trie.Insert(word)
	}

	trie.Remove("apple")
	trie.Remove("api") // removing a word that doesn't exist in the trie
	trie.Remove("")    // removing an empty string
	trie.Remove("banana")
	trie.Remove("band")

	f := func(word string, expected bool) {
		t.Helper()

		if got := trie.Search(word); got != expected {
			t.Errorf("Search(%q) = %v; want %v", word, got, expected)
		}
	}

	f("apple", false)
	f("app", true)
	f("ap", true)
	f("banana", false)
	f("bandit", true)
	f("band", false)
}
