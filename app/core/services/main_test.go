package services

import (
	"database/sql"
	"testing"

	"github.com/hay-kot/repomgr/app/core/bus"
	"github.com/hay-kot/repomgr/app/core/config"
)

type tAppServiceOpts struct {
	recorder *executeRecorder
	cfg      *config.Config
}

func tAppService(t *testing.T, opts ...tAppServiceOpts) *AppService {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	cfg := config.Default()
	exec := &executeRecorder{}

	if len(opts) > 0 {
    cfg = opts[0].cfg 
    exec = opts[0].recorder 
	}

	b := bus.NewEventBus(10)

	service, err := NewAppService(db, cfg, exec, b)
	if err != nil {
		t.Fatal(err)
	}

	return service
}
