package autocomplete

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNextSkipsPaddedSuggestions(t *testing.T) {
	m := Model{
		filtered: []Suggestion{
			{Value: "first"},
			{},
			{Value: "second"},
			{},
		},
		selected: 0,
	}

	m.Next()
	require.Equal(t, 2, m.selected)

	m.Next()
	require.Equal(t, 0, m.selected)
}

func TestPrevSkipsPaddedSuggestions(t *testing.T) {
	m := Model{
		filtered: []Suggestion{
			{Value: "first"},
			{},
			{Value: "second"},
			{},
		},
		selected: 0,
	}

	m.Prev()
	require.Equal(t, 2, m.selected)

	m.Prev()
	require.Equal(t, 0, m.selected)
}

func TestSelectedReturnsEmptyForPaddedSuggestion(t *testing.T) {
	m := Model{
		filtered: []Suggestion{
			{Value: "first"},
			{},
		},
		selected: 1,
	}

	require.Equal(t, Suggestion{}, m.Selected())
}

func TestNavigationNoopsWhenOnlyPaddedSuggestions(t *testing.T) {
	m := Model{
		filtered: []Suggestion{{}, {}, {}},
		selected: 1,
	}

	m.Next()
	require.Equal(t, 1, m.selected)

	m.Prev()
	require.Equal(t, 1, m.selected)

	require.Equal(t, Suggestion{}, m.Selected())
}
