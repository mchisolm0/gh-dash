package repoautocomplete

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dlvhdr/gh-dash/v4/internal/data"
	"github.com/dlvhdr/gh-dash/v4/internal/tui/components/autocomplete"
	"github.com/dlvhdr/gh-dash/v4/internal/tui/components/inputbox"
	"github.com/dlvhdr/gh-dash/v4/internal/tui/components/inputbox/strategies"
)

type RepoLabelsFetchedMsg struct {
	Labels []data.Label
}

type RepoLabelsFetchFailedMsg struct {
	Err error
}

type RepoUsersFetchedMsg struct {
	Users []data.User
}

type RepoUsersFetchFailedMsg struct {
	Err error
}

func LabelSuggestions(labels []data.Label) []autocomplete.Suggestion {
	suggestions := make([]autocomplete.Suggestion, len(labels))
	for i, label := range labels {
		suggestions[i] = autocomplete.Suggestion{Value: label.Name}
	}
	return suggestions
}

func UserSuggestions(users []data.User) []autocomplete.Suggestion {
	suggestions := make([]autocomplete.Suggestion, len(users))
	for i, user := range users {
		suggestions[i] = autocomplete.Suggestion{
			Value:  user.Login,
			Detail: user.Name,
		}
	}
	return suggestions
}

func FetchLabels(repoNameWithOwner string, ac *autocomplete.Model) tea.Cmd {
	spinnerTickCmd := ac.SetFetchLoading()

	fetchCmd := func() tea.Msg {
		labels, err := data.FetchRepoLabels(repoNameWithOwner)
		if err != nil {
			return RepoLabelsFetchFailedMsg{Err: err}
		}
		return RepoLabelsFetchedMsg{Labels: labels}
	}

	return tea.Batch(spinnerTickCmd, fetchCmd)
}

func FetchUsers(repoNameWithOwner string, ac *autocomplete.Model, withSpinner bool) tea.Cmd {
	fetchCmd := func() tea.Msg {
		repoOwner, repoName, ok := strings.Cut(repoNameWithOwner, "/")
		if !ok {
			return RepoUsersFetchFailedMsg{Err: fmt.Errorf("invalid repo name with owner: %q", repoNameWithOwner)}
		}
		users, err := data.FetchRepoUsers(repoName, repoOwner, repoNameWithOwner)
		if err != nil {
			return RepoUsersFetchFailedMsg{Err: err}
		}
		return RepoUsersFetchedMsg{Users: users}
	}

	if !withSpinner {
		return fetchCmd
	}

	spinnerTickCmd := ac.SetFetchLoading()
	return tea.Batch(spinnerTickCmd, fetchCmd)
}

func SetupCommentEntry(inputBox *inputbox.Model, ac *autocomplete.Model) {
	inputBox.Reset()
	ac.Reset()
	inputBox.Autocompleter = strategies.UserMentionCompleter
}

func SetupWhitespaceEntry(inputBox *inputbox.Model, ac *autocomplete.Model) {
	inputBox.Reset()
	ac.Reset()
	inputBox.Autocompleter = strategies.WhitespaceWordCompleter
}

func SetupLabelEntry(inputBox *inputbox.Model) {
	inputBox.Reset()
	inputBox.Autocompleter = strategies.LabelCompleter
}

func SetupUnassignEntry(inputBox *inputbox.Model, ac *autocomplete.Model, resetAutocomplete bool) {
	inputBox.Reset()
	if resetAutocomplete {
		ac.Reset()
	}
}

func ResetSuggestions(ac *autocomplete.Model) {
	ac.Hide()
	ac.SetSuggestions(nil)
}

func SeedUserMentionSuggestions(inputBox inputbox.Model, ac *autocomplete.Model, users []data.User) {
	ac.SetSuggestions(UserSuggestions(users))

	cursorPos := inputBox.CursorPosition()
	mention, _, _ := strategies.UserMentionContextExtractor(inputBox.Value(), cursorPos)
	if mention != "" {
		ac.Show(mention, nil)
	}
}

func SeedWhitespaceSuggestions(inputBox inputbox.Model, ac *autocomplete.Model, users []data.User) {
	ac.SetSuggestions(UserSuggestions(users))

	cursorPos := inputBox.CursorPosition()
	currentWord, _, _ := strategies.WhitespaceContextExtractor(inputBox.Value(), cursorPos)
	existingWords := strategies.WhitespaceItemsToExclude(inputBox.Value(), cursorPos)
	ac.Show(currentWord, existingWords)
}

func SeedLabelSuggestions(inputBox inputbox.Model, ac *autocomplete.Model, labels []data.Label) {
	ac.SetSuggestions(LabelSuggestions(labels))

	cursorPos := inputBox.CursorPosition()
	currentLabel, _, _ := strategies.LabelContextExtractor(inputBox.Value(), cursorPos)
	existingLabels := strategies.LabelItemsToExclude(inputBox.Value(), cursorPos)
	ac.Show(currentLabel, existingLabels)
}

func JoinedListWithTrailingEmpty(items []string, sep string) string {
	values := append([]string{}, items...)
	values = append(values, "")
	return strings.Join(values, sep)
}
