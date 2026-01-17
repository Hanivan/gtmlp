package pipe

import (
	"fmt"
	"sync"
)

// PipeFactory creates a new Pipe instance.
type PipeFactory func() Pipe

// Registry holds registered pipe factories.
var pipeRegistry = struct {
	sync.RWMutex
	factories map[string]PipeFactory
}{
	factories: make(map[string]PipeFactory),
}

// RegisterPipe registers a new pipe factory with the given name.
func RegisterPipe(name string, factory PipeFactory) {
	pipeRegistry.Lock()
	defer pipeRegistry.Unlock()

	pipeRegistry.factories[name] = factory
}

// CreatePipe creates a pipe by name. Returns error if pipe not found.
func CreatePipe(name string) (Pipe, error) {
	pipeRegistry.RLock()
	factory, exists := pipeRegistry.factories[name]
	pipeRegistry.RUnlock()

	if !exists {
		return nil, fmt.Errorf("pipe not found: %s", name)
	}

	return factory(), nil
}

// ListPipes returns all registered pipe names.
func ListPipes() []string {
	pipeRegistry.RLock()
	defer pipeRegistry.RUnlock()

	names := make([]string, 0, len(pipeRegistry.factories))
	for name := range pipeRegistry.factories {
		names = append(names, name)
	}
	return names
}

// UnregisterPipe removes a pipe from the registry.
func UnregisterPipe(name string) {
	pipeRegistry.Lock()
	defer pipeRegistry.Unlock()

	delete(pipeRegistry.factories, name)
}

// init registers all built-in pipes.
func init() {
	// Basic pipes
	RegisterPipe("trim", func() Pipe { return NewTrimPipe() })
	RegisterPipe("lowerCase", func() Pipe { return NewLowerCasePipe() })
	RegisterPipe("upperCase", func() Pipe { return NewUpperCasePipe() })
	RegisterPipe("decode", func() Pipe { return NewDecodePipe() })
	RegisterPipe("trimLeft", func() Pipe { return NewTrimLeftPipe() })
	RegisterPipe("trimRight", func() Pipe { return NewTrimRightPipe() })
	RegisterPipe("stripHTML", func() Pipe { return NewStripHTMLPipe() })

	// Advanced pipes
	RegisterPipe("numNormalize", func() Pipe { return NewNumberNormalizePipe() })
	RegisterPipe("extractEmail", func() Pipe { return NewExtractEmailPipe() })
	RegisterPipe("validateEmail", func() Pipe { return NewValidateEmailPipe() })
	RegisterPipe("validateURL", func() Pipe { return NewValidateURLPipe() })
}
