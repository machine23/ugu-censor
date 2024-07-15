package trie

type trieNode struct {
	children map[rune]*trieNode
	isEnd    bool
}

// Trie represents a trie (prefix tree) data structure for efficient word insertion and search.
type Trie struct {
	root *trieNode
}

// NewTrie creates and returns a new instance of a Trie.
func NewTrie() *Trie {
	return &Trie{
		root: &trieNode{
			children: make(map[rune]*trieNode),
		},
	}
}

// Insert adds a word to the trie.
func (t *Trie) Insert(word string) {
	node := t.root
	for _, c := range word {
		if _, ok := node.children[c]; !ok {
			node.children[c] = &trieNode{
				children: make(map[rune]*trieNode),
			}
		}
		node = node.children[c]
	}
	node.isEnd = true
}

// Search checks if the given word is present in the trie.
//
// Example:
//
//	trie := NewTrie()
//	trie.Insert("hello")
//	found := trie.Search("hello") // Returns true
//	notFound := trie.Search("hell") // Returns false
func (t *Trie) Search(word string) bool {
	node := t.root
	for _, c := range word {
		child, ok := node.children[c]
		if !ok {
			return false
		}
		node = child
	}
	return node.isEnd
}

// StartsWith checks if there is any word in the trie that starts with the given prefix.
// It returns two booleans: the first indicates if a word with the prefix exists,
// and the second indicates if the prefix itself is a complete word in the trie.
//
// Example:
//
//	trie := NewTrie()
//	trie.Insert("apple")
//	hasPrefix, isCompleteWord := trie.StartsWith("app")  // true, false
//	hasPrefix, isCompleteWord = trie.StartsWith("apple") // true, true
//	hasPrefix, isCompleteWord = trie.StartsWith("apz")   // false, false
func (t *Trie) StartsWith(prefix string) (bool, bool) {
	if prefix == "" {
		return false, false
	}

	node := t.root
	for _, c := range prefix {
		child, ok := node.children[c]
		if !ok {
			return false, false
		}
		node = child
	}
	return true, node.isEnd
}

// Remove deletes a word from the trie.
func (t *Trie) Remove(word string) {
	t.remove(t.root, word, 0)
}

func (t *Trie) remove(node *trieNode, word string, index int) bool {
	if index == len(word) {
		if !node.isEnd {
			return false // Word does not exist
		}
		node.isEnd = false             // Mark the end of the word as false
		return len(node.children) == 0 // If no children, node can be deleted
	}

	char := rune(word[index])
	child, ok := node.children[char]
	if !ok {
		return false // Character not found, word does not exist
	}

	// Recursive call to remove the word from the child node
	shouldDeleteChild := t.remove(child, word, index+1)

	if shouldDeleteChild {
		delete(node.children, char) // Delete the child if it has no children
		// Return true if current node is not an end and has no children
		return !node.isEnd && len(node.children) == 0
	}

	return false
}

func (t *Trie) Cursor() *TrieCursor {
	return NewTrieCursor(t.root)
}

// TrieCursor manages traversal within a trie data structure, allowing for character-by-character navigation.
type TrieCursor struct {
	root    *trieNode
	current *trieNode
}

// NewTrieCursor initializes a TrieCursor with a given root node of the trie.
// This function sets both the root and current pointers to the provided rootNode, preparing the cursor for navigation.
// Parameters:
//   - rootNode: A pointer to the root node of the trie from which navigation will start.
//
// Returns:
//   - A pointer to the newly created TrieCursor instance.
func NewTrieCursor(rootNode *trieNode) *TrieCursor {
	return &TrieCursor{
		root:    rootNode,
		current: rootNode,
	}
}

// Advance attempts to move the cursor forward in the trie based on the provided character.
// This method checks if the current node has a child corresponding to the character and, if so, moves the cursor to that child node.
// Parameters:
//   - ch: The rune character to advance the cursor by.
//
// Returns:
//   - A boolean indicating whether the cursor successfully moved to a child node.
//   - A boolean indicating whether the new current node marks the end of a word in the trie.
func (cursor *TrieCursor) Advance(ch rune) (bool, bool) {
	if nextNode, exists := cursor.current.children[ch]; exists {
		cursor.current = nextNode
		return true, nextNode.isEnd
	}
	return false, false
}

// Reset repositions the cursor back to the root of the trie.
// This method is useful for restarting a traversal from the beginning of the trie
// without creating a new TrieCursor instance.
func (cursor *TrieCursor) Reset() {
	cursor.current = cursor.root
}
