package filewatcher

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func Watch(path string, eventChan chan string, shutdown chan struct{}) {
	// get the absolute path for the watched folder
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	// create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	watcher.Add(absPath)

	log.Println("Watching path:", watcher.WatchList())

	go watchLoop(watcher, eventChan)

	// channel that blocks until program is closed
	<-shutdown
	log.Println("FileWatcher stopped")
}

func watchLoop(watcher *fsnotify.Watcher, eventChan chan string) {
	for {
		select {
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("FileWatcher Error:", err)

		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Create) {
				eventChan <- event.Name
			}
			log.Println("FileWatcher Event:", event)
		}
	}
}
