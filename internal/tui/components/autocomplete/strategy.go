package autocomplete

// Context represents the extractable text region around the cursor.
// Start and End are rune indices (not byte indices).
type Context struct {
	Start   int
	End     int
	Content string
}

// Autocompleter defines a strategy for context extraction, suggestion insertion,
// and exclusion filtering.
type Autocompleter interface {
	ExtractContext(input string, cursorPos int) Context
	InsertSuggestion(input string, suggestion string, contextStart int, contextEnd int) (newInput string, newCursorPos int)
	ItemsToExclude(input string, cursorPos int) []string
}
