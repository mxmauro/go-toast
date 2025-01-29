package toast

import (
	"errors"
	"sync"
)

// -----------------------------------------------------------------------------

type Options map[string]interface{}

// -----------------------------------------------------------------------------

var mtx = sync.RWMutex{}
var initialized = false

// -----------------------------------------------------------------------------

func Initialize() error {
	mtx.Lock()
	defer mtx.Unlock()

	err := toastInit()
	if err != nil {
		return err
	}
	initialized = true
	return nil
}

func Finalize() {
	mtx.Lock()
	defer mtx.Unlock()

	if initialized {
		initialized = false
		toastDone()
	}
}

func Show(opts Options) error {
	mtx.RLock()
	defer mtx.RUnlock()

	if !initialized {
		return errors.New("not initialized")
	}

	err := toastShow(opts)

	// Done
	return err
}
