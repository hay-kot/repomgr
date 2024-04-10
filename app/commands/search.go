package commands

import (
	"context"
	"fmt"
)

func (ctrl *Controller) Search(ctx context.Context) error {
	r, err := ctrl.repos.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, repo := range r {
		fmt.Println(repo.Name)
	}

	return nil
}
