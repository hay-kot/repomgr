package commands

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/repomgr/app/commands/ui"
)

func (ctrl *Controller) Search(ctx context.Context) error {
	r, err := ctrl.repos.GetAll(ctx)
	if err != nil {
		return err
	}

	search := ui.NewSearchView(ui.NewSearch(ctrl.conf.KeyBindings, r))
	layout := ui.NewLayout(search)

	p := tea.NewProgram(layout, tea.WithAltScreen())

	_, err = p.Run()
	if err != nil {
		return err
	}

	return nil
}
