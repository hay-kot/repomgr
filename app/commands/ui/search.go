package ui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/repomgr/app/core/config"
	"github.com/hay-kot/repomgr/app/core/services"
	"github.com/hay-kot/repomgr/app/repos"
	"github.com/hay-kot/repomgr/internal/icons"
	"github.com/hay-kot/repomgr/internal/styles"
	"github.com/rs/zerolog/log"
	"github.com/sahilm/fuzzy"
)

// SearchCtrl is the controller/model for the fuzzer finding UI component
type SearchCtrl struct {
	rs           *services.RepositoryService
	index        []string
	repos        []repos.Repository
	searchLength int
	selected     int
	// key = filtered index, value = original index
	indexmap map[int]int
	limit    int

	keybinds config.KeyBindings
}

func NewSearchCtrl(
	rs *services.RepositoryService,
	bindings config.KeyBindings,
	r []repos.Repository,
) *SearchCtrl {
	return &SearchCtrl{
		rs:       rs,
		repos:    r,
		keybinds: bindings,
	}
}

// Selected returns the active selection by the user, or any empty object
// if no selection has been made OR the active index is out of range.
func (c *SearchCtrl) Selected() repos.Repository {
	if c.indexmap == nil {
		if c.selected < 0 || c.selected >= len(c.repos) {
			return repos.Repository{}
		}
		return c.repos[c.selected]
	}

	idx, ok := c.indexmap[c.selected]
	if !ok {
		return repos.Repository{}
	}

	if idx < 0 || idx >= len(c.repos) {
		return repos.Repository{}
	}

	return c.repos[idx]
}

// search returns a sorted list of matches uses a fuzzy search algorithm
func (c *SearchCtrl) search(str string) []repos.Repository {
	if str == "" {
		c.searchLength = len(c.repos)
		return c.repos
	}

	c.indexmap = make(map[int]int)

	if c.index == nil {
		c.index = make([]string, len(c.repos))
		for i, repo := range c.repos {
			c.index[i] = repo.DisplayName()
		}
	}

	matches := fuzzy.Find(str, c.index)
	results := make([]repos.Repository, len(matches))
	for i, match := range matches {
		results[i] = c.repos[match.Index]
		c.indexmap[i] = match.Index
	}

	c.searchLength = len(results)
	return results
}

type SearchView struct {
	ctrl   *SearchCtrl
	cmd    *services.CommandService
	search textinput.Model
	height int
	shift  int
}

func NewSearchView(ctrl *SearchCtrl, service *services.CommandService) SearchView {
	ti := textinput.New()
	ti.Focus()
	ti.Prompt = styles.AccentBlue("> ")
	ti.CharLimit = 256
	ti.Width = 80

	return SearchView{
		ctrl:   ctrl,
		cmd:    service,
		search: ti,
	}
}

func (m SearchView) Init() tea.Cmd {
	return textinput.Blink
}

func (m SearchView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

		switch msg.String() {
		case "up":
			if m.ctrl.selected > 0 {
				m.ctrl.selected--

				if m.shift > 0 {
					m.shift--
				}
			}
		case "down":
			if m.ctrl.selected < m.ctrl.limit-1 {
				m.ctrl.selected++

				if m.ctrl.selected >= len(m.ctrl.repos) {
					m.shift++
				}
			}
		default:
			ok, cmd := m.cmd.GetBoundCommand(msg.Type)
			if !ok {
				break
			}

			err := m.cmd.Run(m.ctrl.Selected(), cmd)
			if err != nil {
				log.Err(err).Msg("failed to run command")
				return m, tea.Quit
			}
		}
	}

	old := m.search

	m.search, cmd = m.search.Update(msg)
	if old.Value() != m.search.Value() {
		m.ctrl.selected = 0
	}

	return m, cmd
}

func (m SearchView) View() string {
	results := m.ctrl.search(m.search.Value())
	str := strings.Builder{}

	// Calculate the number of allowed_rows we can display
	m.ctrl.limit = m.height - 8

	var determinedMax int
	if m.ctrl.limit < 0 {
		determinedMax = len(results)
	} else if len(results) > m.ctrl.limit {
		determinedMax = m.ctrl.limit
	} else {
		determinedMax = len(results)
	}

	m.ctrl.limit = determinedMax

	if m.ctrl.selected > m.ctrl.limit {
		m.ctrl.selected = 0
	}

	str.WriteString(m.search.View())
	str.WriteString(styles.Subtle(fmt.Sprintf("\n  %d/%d", len(results), len(m.ctrl.repos))) + "\n")
	str.WriteString(m.fmtMatches(results[:determinedMax]))

	// fill remaining height - 1
	for i := 0; i < m.height-determinedMax-3; i++ {
		str.WriteString("\n")
	}

	str.WriteString(m.keyHelp())
	return str.String()
}

func (m SearchView) keyHelp() string {
	keys := make([]string, 0, len(m.ctrl.keybinds))
	for key, cmd := range m.ctrl.keybinds {
		keys = append(keys, fmt.Sprintf("%s: %s", key, cmd.Desc))
	}
	slices.Sort(keys)

	bldr := strings.Builder{}
	for i, key := range keys {
		bldr.WriteString(styles.Subtle(" " + key))
		if i < len(keys)-1 {
			bldr.WriteString(" ")
			bldr.WriteString(styles.Subtle(icons.Dot))
		}
	}

	return bldr.String()
}

func (m SearchView) fmtMatches(repos []repos.Repository) string {
	longest := 0

	for _, repo := range repos {
		if len(repo.Name) > longest {
			longest = len(repo.DisplayName())
		}
	}

	search := m.search.Value()

	str := strings.Builder{}
	for i, repo := range repos {
		spaces := (longest + 5) - len(repo.Name)

		var (
			prefix     = " "
			iconPrefix = " "
			iconSpace  = "    "
		)

		if repo.IsFork {
			iconPrefix += styles.Subtle(icons.Fork) + " "
		} else {
			iconPrefix += "  " // double width icon
		}

		if m.ctrl.rs.IsCloned(repo) {
			iconPrefix += styles.Subtle(icons.Folder) + " "
		} else {
			iconPrefix += " "
		}

		text := "github.com/" + repo.DisplayName() + strings.Repeat(" ", spaces)

		if m.ctrl.selected == i {
			prefix = styles.HighlightRow(styles.AccentRed(">"))
			text = styles.HighlightRow(styles.Bold.Render(text))
			iconPrefix = styles.HighlightRow(styles.Bold.Render(iconPrefix))
			iconSpace = styles.HighlightRow(styles.Bold.Render(iconSpace))
		} else {
			if search != "" && strings.Contains(repo.Name, search) {
				// Highlight the search term
				text = strings.ReplaceAll(text, search, styles.Bold.Render(search))
			}
		}

		str.WriteString(prefix + iconPrefix + iconSpace + text + "\n")
	}

	return str.String()
}
