package repoautocomplete

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dlvhdr/gh-dash/v4/internal/data"
	"github.com/dlvhdr/gh-dash/v4/internal/tui/components/autocomplete"
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
