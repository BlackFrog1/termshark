// Copyright 2019 Graham Clark. All rights reserved.  Use of this source
// code is governed by the MIT license that can be found in the LICENSE
// file.

package termshark

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"gopkg.in/fsnotify.v1"
)

//======================================================================

var GoRoutineWg *sync.WaitGroup

type ConfigWatcher struct {
	watcher *fsnotify.Watcher
	change  chan struct{}
	close   chan struct{}
}

func NewConfigWatcher() (*ConfigWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	change := make(chan struct{})
	cl := make(chan struct{})

	res := &ConfigWatcher{
		change: change,
		close:  cl,
	}

	TrackedGo(func() {
	Loop:
		for {
			select {
			// watch for events
			case <-watcher.Events:
				res.change <- struct{}{}

			case err := <-watcher.Errors:
				log.Debugf("Error from config watcher: %v", err)

			case <-cl:
				break Loop
			}
		}
	}, GoRoutineWg)

	if err := watcher.Add(ConfFile("termshark.toml")); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	res.watcher = watcher

	return res, nil
}

func (c *ConfigWatcher) Close() error {
	c.close <- struct{}{}
	return c.watcher.Close()
}

func (c *ConfigWatcher) ConfigChanged() <-chan struct{} {
	return c.change
}

//======================================================================
// Local Variables:
// mode: Go
// fill-column: 78
// End:
