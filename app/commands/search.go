package commands

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/repomgr/app/commands/ui"
	"github.com/hay-kot/repomgr/app/core/services"
	"github.com/hay-kot/repomgr/app/repos"
)

func (ctrl *Controller) Search(ctx context.Context) (repos.Repository, error) {
	r, err := ctrl.repos.GetAll(ctx)
	if err != nil {
		return repos.Repository{}, err
	}

	var (
		exec       = services.NewShellExecutor(ctrl.conf.Shell)
		cmdService = services.NewCommandService(ctrl.conf.CloneDirectories, ctrl.conf.KeyBindings, exec, ctrl.bus)
		searchCtrl = ui.NewSearchCtrl(ctrl.repos, ctrl.conf.KeyBindings, r)
		search     = ui.NewSearchView(searchCtrl, cmdService)
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
