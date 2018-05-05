package door

import (
	"sync"
)

var defaultDirectory Directory
var once sync.Once

func DefaultDirectory() Directory {
	once.Do(func() {
		defaultDirectory = NewDirectory()
	})
	return defaultDirectory
}
