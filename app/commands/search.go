package commands

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/repomgr/app/commands/ui"
)

func (ctrl *Controller) Search(ctx context.Context) (string, error) {
	r, err := ctrl.store.GetAll(ctx)
	if err != nil {
		return "", err
	}

	var (
		searchCtrl = ui.NewSearchCtrl(r, ctrl.rfs, ctrl.commander)
		search     = ui.NewSearchView(searchCtrl)
		layout     = ui.NewLayout(search)
	)

	p := tea.NewProgram(layout, tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		return "", err
	}

	msg := searchCtrl.ExitMessage()
	return msg, nil
}
