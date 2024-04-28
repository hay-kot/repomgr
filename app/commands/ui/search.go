package ui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/repomgr/app/core/commander"
	"github.com/hay-kot/repomgr/app/core/repofs"
	"github.com/hay-kot/repomgr/app/repos"
	"github.com/hay-kot/repomgr/internal/icons"
	"github.com/hay-kot/repomgr/internal/styles"
	"github.com/rs/zerolog/log"
	"github.com/sahilm/fuzzy"
)

// SearchCtrl is the controller/model for the fuzzer finding UI component
type SearchCtrl struct {
	*state

	index        []string
	repos        []repos.Repository
	searchLength int
	selected     int

	// key = filtered index, value = original index
	indexmap map[int]int
	limit    int

	rfs       *repofs.RepoFS
	commander *commander.Commander
	keybinds  commander.KeyBindings
}

func NewSearchCtrl(r []repos.Repository, rfs *repofs.RepoFS, cmd *commander.Commander) *SearchCtrl {
	return &SearchCtrl{
		state:     &state{},
		repos:     r,
		rfs:       rfs,
		commander: cmd,
		keybinds:  cmd.Bindings(),
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
		c.indexmap = nil
		return c.repos
	}

	if c.index == nil {
		c.index = make([]string, len(c.repos))
		for i, repo := range c.repos {
			c.index[i] = repo.DisplayName()
		}
	}

	matches := fuzzy.Find(str, c.index)
	c.indexmap = make(map[int]int, len(matches))
	results := make([]repos.Repository, len(matches))
	for i, match := range matches {
		results[i] = c.repos[match.Index]
		c.indexmap[i] = match.Index
	}

	c.searchLength = len(results)
	return results
}

type SearchView struct {
	results []repos.Repository
	ctrl    *SearchCtrl
	search  textinput.Model
	height  int
	shift   int
}

func NewSearchView(ctrl *SearchCtrl) *SearchView {
	ti := textinput.New()
	ti.Focus()
	ti.Prompt = styles.AccentBlue("> ")
	ti.CharLimit = 256
	ti.Width = 80

	return &SearchView{
		ctrl:   ctrl,
		search: ti,
	}
}

func (m *SearchView) Init() tea.Cmd {
	return textinput.Blink
}

func (m *SearchView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.ctrl.isExit() {
		return m, tea.Quit
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

		msgstr := msg.String()
		switch {
		case "up" == msgstr:
			if m.ctrl.selected > 0 {
				m.ctrl.selected--

				if m.shift > 0 {
					m.shift--
				}
			}
		case "down" == msgstr:
			if m.ctrl.selected < m.ctrl.limit-1 {
				m.ctrl.selected++

				if m.ctrl.selected >= len(m.ctrl.repos) {
					m.shift++
				}
			}
		case msgstr == "enter", strings.HasPrefix(msgstr, "ctrl"):
			action, ok := m.ctrl.commander.GetAction(msg.Type.String(), m.ctrl.Selected())
			if !ok {
				break
			}

			if action.IsExit() {
				m.ctrl.signalExit(action.ExitMessage())
				return m, tea.Quit
			}

			switch action.Mode {
			case commander.ModeReadOnly:
				cmdModel := NewCommandView(action, m)
				return cmdModel, cmdModel.Init()
			case commander.ModeInteractive:
				cmd, ok := action.IsExec()
				if !ok {
					log.Error().Msg("action defined as interactive but no exec command found")
					break
				}

				teacmd := tea.ExecProcess(cmd, func(err error) tea.Msg {
					if err != nil {
						log.Error().Err(err).Msg("failed to execute command")
						return tea.Quit
					}

					return nil
				})

				return m, teacmd
			case commander.ModeBackground:
				ch := action.GoRun()
				m.ctrl.recieveError(ch)
			default:
				panic("unhandled mode " + action.Mode)
			}
		}
	}

	old := m.search.Value()
	m.search, cmd = m.search.Update(msg)
	if old != m.search.Value() {
		m.ctrl.selected = 0
	}

	m.results = m.ctrl.search(m.search.Value())
	return m, cmd
}

func (m *SearchView) View() string {
	results := m.results
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

func (m *SearchView) keyHelp() string {
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

		if m.ctrl.rfs.IsCloned(repo) {
			iconPrefix += styles.Subtle(icons.Folder) + " "
		} else {
			iconPrefix += "  "
		}

		text := "github.com/" + repo.DisplayName() + strings.Repeat(" ", spaces)

		if m.ctrl.selected == i {
			prefix = styles.HighlightRow(styles.AccentRed(">"))
			text = styles.HighlightRow(styles.Bold.Render(text))
			iconPrefix = styles.HighlightRow(iconPrefix)
			iconSpace = styles.HighlightRow(iconSpace)
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
