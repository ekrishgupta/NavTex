package core

import (
	"log"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher listens for file changes in the workspace and emits events.
type Watcher struct {
	watcher *fsnotify.Watcher
	Events  chan string
	done    chan bool
}

// NewWatcher creates a new file system watcher for the given directory.
func NewWatcher(dir string) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		w.Close()
		return nil, err
	}

	err = w.Add(absDir)
	if err != nil {
		w.Close()
		return nil, err
	}

	sub := &Watcher{
		watcher: w,
		Events:  make(chan string, 10),
		done:    make(chan bool),
	}

	go sub.listen()

	return sub, nil
}

// listen runs in a goroutine and forwards relevant events with debouncing.
func (w *Watcher) listen() {
	var debounceTimer <-chan time.Time
	var pendingEvent string

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			// We only care about writes, creates, or removes
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
				pendingEvent = event.Name
				// Reset the debounce timer on every new event
				debounceTimer = time.After(200 * time.Millisecond)
			}
		case <-debounceTimer:
			// Fire the debounced event to the UI
			w.Events <- pendingEvent
			debounceTimer = nil // Stop firing until another event
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		case <-w.done:
			return
		}
	}
}

// Close stops the watcher.
func (w *Watcher) Close() {
	w.done <- true
	w.watcher.Close()
	close(w.Events)
}
