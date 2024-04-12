package bus

import (
	"context"
	"sync"

	"github.com/hay-kot/repomgr/app/repos"
	"github.com/rs/zerolog/log"
)

// EventBus is a simple in-process event bus for the application to handle event pub/sub for. This
// is intended for small, short lived event handlers like cache busting for a single (or few) items.
//
// For longer running background tasks look at the tasks package in the sys folder. That provides
// a more robust solution for background tasks.
type EventBus struct {
	started bool
	ch      chan eventData

	mu          sync.RWMutex
	subscribers map[Topic][]func(any)
}

func NewEventBus(size int) *EventBus {
	return &EventBus{
		ch: make(chan eventData, size),
		subscribers: map[Topic][]func(any){
			TopicRepoCloned: {},
		},
	}
}

func (e *EventBus) Start(ctx context.Context) error {
	if e.started {
		panic("event bus already started")
	}

	if len(e.subscribers) == 0 {
		panic("no subscribers, you must have at least one")
	}

	e.started = true

	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-e.ch:
			e.mu.RLock()
			arr, ok := e.subscribers[event.topic]
			e.mu.RUnlock()

			if !ok {
				log.Warn().Str("topic", string(event.topic)).Msg("no subscribers")
				continue
			}

			log.Debug().Int("listeners", len(arr)).
				Str("topic", string(event.topic)).Msg("event received")
			for _, fn := range arr {
				fn(event.data)
			}
		}
	}
}

func (e *EventBus) publish(topic Topic, data any) {
	log.Debug().Str("topic", string(topic)).Msg("publishing event")
	e.ch <- eventData{
		topic: topic,
		data:  data,
	}
}

func (e *EventBus) subscribe(event Topic, fn func(any)) {
	e.mu.Lock()
	defer e.mu.Unlock()

	arr, ok := e.subscribers[event]
	if !ok {
		panic("event not found")
	}

	e.subscribers[event] = append(arr, fn)
}

func (e *EventBus) PubCloneEvent(repo repos.Repository, cloneDir string) {
	e.publish(TopicRepoCloned, RepoClonedEvent{
		Repo:     repo,
		CloneDir: cloneDir,
	})
}

func (e *EventBus) SubCloneEvent(fn func(RepoClonedEvent)) {
	e.subscribe(TopicRepoCloned, func(a any) {
		fn(a.(RepoClonedEvent))
	})
}
