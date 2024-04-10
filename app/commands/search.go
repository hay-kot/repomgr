package commands

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/repomgr/app/commands/ui"
	"github.com/hay-kot/repomgr/app/repos"
)

func (ctrl *Controller) Search(ctx context.Context) (repos.Repository, error) {
	r, err := ctrl.repos.GetAll(ctx)
	if err != nil {
		return repos.Repository{}, err
	}

	var (
		searchCtrl = ui.NewSearchCtrl(ctrl.conf.KeyBindings, r)
		search     = ui.NewSearchView(searchCtrl)
		layout     = ui.NewLayout(search)
	)

	p := tea.NewProgram(layout, tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		return repos.Repository{}, err
	}

	selected := searchCtrl.Selected()

	return selected, nil
}
