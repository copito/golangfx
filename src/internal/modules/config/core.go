package config

import (
	"sync"

	"github.com/copito/runner/src/internal/entities"
)

// ConfigChangeCallback is called when configuration changes. It receives the old and new configuration as parameters.
type ConfigChangeCallback func(oldConfig, newConfig *entities.Config)

type ConfigProvider interface {
	Get() *entities.Config
	Set(newConfig entities.Config)
	OnChange(callback ConfigChangeCallback)
}

// Ensure compliance with ConfigProvider interface
var _ ConfigProvider = (*configProvider)(nil)

type configProvider struct {
	config    *entities.Config
	Callbacks []ConfigChangeCallback
	mu        sync.RWMutex
}

func NewConfigProvider(initialConfig entities.Config) ConfigProvider {
	return &configProvider{
		config:    &initialConfig,
		Callbacks: make([]ConfigChangeCallback, 0),
	}
}

func (cp *configProvider) Get() *entities.Config {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return cp.config
}

func (cp *configProvider) Set(newConfig entities.Config) {
	cp.mu.Lock()
	oldConfig := cp.config
	cp.config = &newConfig
	cp.mu.Unlock()

	// Notify all registered callbacks about the configuration change
	for _, callback := range cp.Callbacks {
		callback(oldConfig, &newConfig)
	}
}

func (cp *configProvider) OnChange(callback ConfigChangeCallback) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.Callbacks = append(cp.Callbacks, callback)
}
