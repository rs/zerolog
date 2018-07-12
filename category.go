package zerolog

import (
	"math"
	"sync"
)

// Category represents log event categories.
type Category uint8

type categoryDef struct {
	name string
	lvl  Level
}

const maxcategories = math.MaxUint8

var (
	categories [maxcategories]categoryDef
	catmx      sync.RWMutex
)

// SetCategory sets the name and the log level for category cat.
//
// This function is not safe for use by multiple goroutines.
func SetCategory(cat Category, name string, lvl Level) {
	categories[cat].lvl = lvl
	categories[cat].name = name
}

// SetCategoryLevel sets the level for category cat, that is, the minimum
// accepted log level for events of that category to be actually logged. Events
// can be assigned a category by calling Event.Cat().
//
// This function is safe for use by multiple goroutines.
func SetCategoryLevel(cat Category, lvl Level) {
	catmx.Lock()
	SetCategoryLevelUnsafe(cat, lvl)
	catmx.Unlock()
}

// LockCategoryLevels locks the underlying store of category log levels.
//
// When you're done calling SetCategoryLevelUnsafe() you must call UnlockCategoryLevelUnsafe.
func LockCategoryLevels() { catmx.Lock() }

// UnlockCategoryLevels unlocks the underlying store of category log levels.
func UnlockCategoryLevels() { catmx.Unlock() }

// SetCategoryLevelUnsafe sets the level for category cat, that is, the minimum
// accepted log level for events of that category to be actually logged. Events
// can be assigned a category by calling Event.Cat().
//
// NOTICE: useful for batch updating of category log levels when other goroutines are running.
func SetCategoryLevelUnsafe(cat Category, lvl Level) { categories[cat].lvl = lvl }

// shouldCategory returns true if an event with category cat and level lvl
// should be logged. name is set to the category name if that is the case.
func shouldCategory(cat Category, lvl Level) (ok bool, name string) {
	catmx.RLock()
	if lvl >= categories[cat].lvl {
		ok = true
		name = categories[cat].name
	}
	catmx.RUnlock()
	return
}
