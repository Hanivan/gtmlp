package gtmlp

import (
	"strings"
	"sync"
)

// pipeRegistry holds all registered pipes
var (
	pipeRegistry   = make(map[string]PipeFunc)
	registryMutex  sync.RWMutex
)

// RegisterPipe registers a custom pipe function
func RegisterPipe(name string, fn PipeFunc) {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	pipeRegistry[strings.ToLower(name)] = fn
}

// parsePipeDefinition parses a pipe definition like "pipeName:param1:param2"
// Returns pipe name (lowercase) and parameters
func parsePipeDefinition(def string) (string, []string) {
	parts := strings.Split(def, ":")
	if len(parts) == 1 {
		return strings.ToLower(parts[0]), nil
	}
	return strings.ToLower(parts[0]), parts[1:]
}

// getPipe retrieves a pipe from registry (case-insensitive)
func getPipe(name string) PipeFunc {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	return pipeRegistry[strings.ToLower(name)]
}
