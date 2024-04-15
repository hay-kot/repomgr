package commands

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/repomgr/app/commands/ui"
)

func (ctrl *Controller) Search(ctx context.Context) (string, error) {
	r, err := ctrl.app.GetAll(ctx)
	if err != nil {
		return "", err
	}

	var (
		searchCtrl = ui.NewSearchCtrl(ctrl.conf.KeyBindings, r)
		search     = ui.NewSearchView(searchCtrl, ctrl.app)
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
