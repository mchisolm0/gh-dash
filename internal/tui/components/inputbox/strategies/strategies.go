package strategies

import (
	"strings"

	"github.com/dlvhdr/gh-dash/v4/internal/tui/components/autocomplete"
)

type LabelInfo struct {
	Label    string
	StartIdx int
	EndIdx   int
	IsFirst  bool
	IsLast   bool
}

type WordInfo struct {
	Word     string
	StartIdx int
	EndIdx   int
	IsFirst  bool
	IsLast   bool
}

type labelAutocompleter struct{}
type userMentionAutocompleter struct{}
type whitespaceAutocompleter struct{}

var (
	LabelCompleter          autocomplete.Autocompleter = labelAutocompleter{}
	UserMentionCompleter    autocomplete.Autocompleter = userMentionAutocompleter{}
	WhitespaceWordCompleter autocomplete.Autocompleter = whitespaceAutocompleter{}
)

func (labelAutocompleter) ExtractContext(input string, cursorPos int) autocomplete.Context {
	content, start, end := LabelContextExtractor(input, cursorPos)
	return autocomplete.Context{
		Start:   start,
		End:     end,
		Content: content,
	}
}

func (labelAutocompleter) InsertSuggestion(input string, suggestion string, contextStart int, contextEnd int) (newInput string, newCursorPos int) {
	return LabelSuggestionInserter(input, suggestion, contextStart, contextEnd)
}

func (labelAutocompleter) ItemsToExclude(input string, cursorPos int) []string {
	return LabelItemsToExclude(input, cursorPos)
}

func (userMentionAutocompleter) ExtractContext(input string, cursorPos int) autocomplete.Context {
	content, start, end := UserMentionContextExtractor(input, cursorPos)
	return autocomplete.Context{
		Start:   start,
		End:     end,
		Content: content,
	}
}

func (userMentionAutocompleter) InsertSuggestion(input string, suggestion string, contextStart int, contextEnd int) (newInput string, newCursorPos int) {
	return UserMentionSuggestionInserter(input, suggestion, contextStart, contextEnd)
}

func (userMentionAutocompleter) ItemsToExclude(input string, cursorPos int) []string {
	return UserMentionItemsToExclude(input, cursorPos)
}

func (whitespaceAutocompleter) ExtractContext(input string, cursorPos int) autocomplete.Context {
	content, start, end := WhitespaceContextExtractor(input, cursorPos)
	return autocomplete.Context{
		Start:   start,
		End:     end,
		Content: content,
	}
}

func (whitespaceAutocompleter) InsertSuggestion(input string, suggestion string, contextStart int, contextEnd int) (newInput string, newCursorPos int) {
	return WhitespaceSuggestionInserter(input, suggestion, contextStart, contextEnd)
}

func (whitespaceAutocompleter) ItemsToExclude(input string, cursorPos int) []string {
	return WhitespaceItemsToExclude(input, cursorPos)
}

func LabelContextExtractor(input string, cursorPos int) (context string, start int, end int) {
	info := ExtractLabelAtCursor(input, cursorPos)
	return info.Label, info.StartIdx, info.EndIdx
}

func ExtractLabelAtCursor(input string, cursorPos int) LabelInfo {
	if input == "" {
		return LabelInfo{
			Label:    "",
			StartIdx: 0,
			EndIdx:   0,
			IsFirst:  true,
			IsLast:   true,
		}
	}

	runes := []rune(input)
	if cursorPos < 0 {
		cursorPos = 0
	}
	if cursorPos > len(runes) {
		cursorPos = len(runes)
	}

	startIdx := 0
	for i := cursorPos - 1; i >= 0; i-- {
		if runes[i] == ',' {
			startIdx = i + 1
			break
		}
	}

	endIdx := len(runes)
	for i := cursorPos; i < len(runes); i++ {
		if runes[i] == ',' {
			endIdx = i
			break
		}
	}

	label := strings.TrimSpace(string(runes[startIdx:endIdx]))
	isFirst := startIdx == 0
	isLast := endIdx == len(runes)

	return LabelInfo{
		Label:    label,
		StartIdx: startIdx,
		EndIdx:   endIdx,
		IsFirst:  isFirst,
		IsLast:   isLast,
	}
}

func CurrentLabels(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	labels := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			labels = append(labels, trimmed)
		}
	}
	return labels
}

func LabelItemsToExclude(input string, cursorPos int) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}
	currentLabels := CurrentLabels(input)
	if currentLabels == nil {
		return nil
	}
	excluded := make([]string, 0, len(currentLabels))
	for _, label := range currentLabels {
		if label != "" {
			excluded = append(excluded, label)
		}
	}
	return excluded
}

func LabelSuggestionInserter(input string, suggestion string, contextStart int, contextEnd int) (newInput string, newCursorPos int) {
	labelInfo := ExtractLabelAtCursor(input, contextStart)
	runes := []rune(input)

	var replacement string
	if labelInfo.IsFirst {
		replacement = suggestion + ", "
	} else {
		replacement = " " + suggestion + ", "
	}

	remainingInput := string(runes[labelInfo.EndIdx:])
	remainingInput = strings.TrimPrefix(remainingInput, ",")
	remainingInput = strings.TrimLeft(remainingInput, " \t")

	newValue := string(runes[:labelInfo.StartIdx]) + replacement + remainingInput
	newCursorPos = labelInfo.StartIdx + len([]rune(replacement))

	return newValue, newCursorPos
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isWordBoundary(r rune) bool {
	return isWhitespace(r) || r == ',' || r == '.' || r == '!' || r == '?' || r == ';' || r == ':' || r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}' || r == '<' || r == '>' || r == '"' || r == '\'' || r == '`'
}

func UserMentionContextExtractor(input string, cursorPos int) (context string, start int, end int) {
	if input == "" {
		return "", -1, -1
	}

	runes := []rune(input)
	if cursorPos < 0 {
		cursorPos = 0
	}
	if cursorPos > len(runes) {
		cursorPos = len(runes)
	}

	wordStart := 0
	for i := cursorPos - 1; i >= 0; i-- {
		if isWordBoundary(runes[i]) {
			wordStart = i + 1
			break
		}
		wordStart = i
	}

	if wordStart >= len(runes) {
		return "", -1, -1
	}
	if runes[wordStart] != '@' {
		return "", -1, -1
	}

	wordEnd := len(runes)
	for i := cursorPos; i < len(runes); i++ {
		if isWordBoundary(runes[i]) {
			wordEnd = i
			break
		}
	}

	mentionStart := wordStart + 1
	mentionText := string(runes[mentionStart:wordEnd])
	return mentionText, wordStart, wordEnd
}

func UserMentionSuggestionInserter(input string, suggestion string, contextStart int, contextEnd int) (newInput string, newCursorPos int) {
	runes := []rune(input)
	replacement := "@" + suggestion + " "
	newValue := string(runes[:contextStart]) + replacement + string(runes[contextEnd:])
	newCursorPos = contextStart + len([]rune(replacement))
	return newValue, newCursorPos
}

func UserMentionItemsToExclude(input string, cursorPos int) []string {
	return nil
}

func WhitespaceContextExtractor(input string, cursorPos int) (context string, start int, end int) {
	info := ExtractWordAtCursor(input, cursorPos)
	return info.Word, info.StartIdx, info.EndIdx
}

func ExtractWordAtCursor(input string, cursorPos int) WordInfo {
	if input == "" {
		return WordInfo{
			Word:     "",
			StartIdx: 0,
			EndIdx:   0,
			IsFirst:  true,
			IsLast:   true,
		}
	}

	runes := []rune(input)
	if cursorPos < 0 {
		cursorPos = 0
	}
	if cursorPos > len(runes) {
		cursorPos = len(runes)
	}

	wordStart := 0
	for i := cursorPos - 1; i >= 0; i-- {
		if isWhitespace(runes[i]) {
			wordStart = i + 1
			break
		}
		wordStart = i
	}

	wordEnd := len(runes)
	for i := cursorPos; i < len(runes); i++ {
		if isWhitespace(runes[i]) {
			wordEnd = i
			break
		}
	}

	wordText := strings.TrimSpace(string((runes[wordStart:wordEnd])))
	isFirst := wordStart == 0
	isLast := wordEnd == len(runes)

	return WordInfo{
		Word:     wordText,
		StartIdx: wordStart,
		EndIdx:   wordEnd,
		IsFirst:  isFirst,
		IsLast:   isLast,
	}
}

func WhitespaceSuggestionInserter(input string, suggestion string, contextStart int, contextEnd int) (newInput string, newCursorPos int) {
	runes := []rune(input)
	replacement := suggestion + " "
	newValue := string(runes[:contextStart]) + replacement + string(runes[contextEnd:])
	newCursorPos = contextStart + len([]rune(replacement))
	return newValue, newCursorPos
}

func WhitespaceItemsToExclude(input string, cursorPos int) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}

	wordInfo := ExtractWordAtCursor(input, cursorPos)
	allWords := AllWords(input)
	if allWords == nil {
		return nil
	}

	excluded := make([]string, 0, len(allWords))
	for _, word := range allWords {
		if word != wordInfo.Word {
			excluded = append(excluded, word)
		}
	}

	return excluded
}

func AllWords(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Fields(value)
	words := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			words = append(words, trimmed)
		}
	}
	return words
}
