package services

import (
	"database/sql"
	"testing"

	"github.com/hay-kot/repomgr/app/core/bus"
	"github.com/hay-kot/repomgr/app/core/config"
)

func tAppService(t *testing.T) *AppService {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	b := bus.NewEventBus(10)
	exec := &executeRecorder{}

	conf := &config.Config{
		Concurrency:      8,
		Shell:            "zsh",
		KeyBindings:      map[string]config.KeyCommand{},
		Sources:          []config.Source{},
		Database:         config.Database{},
		Logs:             config.Logs{},
		CloneDirectories: config.CloneDirectories{},
	}

	service, err := NewAppService(db, conf, exec, b)
	if err != nil {
		t.Fatal(err)
	}

	return service
}
