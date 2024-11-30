package utils

import (
	"sync"
)

// Injector struct to hold the singleton objects
type Injector struct {
	singletonObj map[string]interface{}
	mu           sync.Mutex // To handle concurrency
}

// NewInjector creates a new Injector instance
func NewInjector() *Injector {
	return &Injector{
		singletonObj: make(map[string]interface{}),
	}
}

// Get requests the singleton with the given name
func (i *Injector) Get(name string) interface{} {
	i.mu.Lock()
	defer i.mu.Unlock()

	if obj, exists := i.singletonObj[name]; exists {
		return obj
	}
	return nil
}

// Bind puts a singleton into memory
func (i *Injector) Bind(singleton interface{}, name string) bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exists := i.singletonObj[name]; exists {
		return false
	}
	i.singletonObj[name] = singleton
	return true
}

// Update updates a singleton in memory
func (i *Injector) Update(singleton interface{}, name string) bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exists := i.singletonObj[name]; exists {
		i.singletonObj[name] = singleton
		return true
	}
	i.singletonObj[name] = singleton
	return false
}

// Delete removes a singleton from memory
func (i *Injector) Delete(name string) bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exists := i.singletonObj[name]; exists {
		delete(i.singletonObj, name)
		return true
	}
	return false
}

// Clear removes all singletons from memory
func (i *Injector) Clear() {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.singletonObj = make(map[string]interface{})
}

// SingletonInjector is the global instance of Injector
var SingletonInjector = NewInjector()
