// Copyright 2023 jim.zoumo@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"fmt"
	"sync"
)

// Registry provides a place binding name and interface{}
type Registry interface {
	// Register registers a interface by name.
	// It will panic if name corresponds to an already registered interface
	// and the registry does not allow user to override the interface.
	Register(name string, v interface{}) error

	// Get returns an interface registered with the given name
	Get(name string) (interface{}, bool)

	// Range calls f sequentially for each key and value present in the registry.
	// If f returns false, range stops the iteration.
	Range(func(key string, value interface{}) bool)

	// Keys returns the name of all registered interfaces
	Keys() []string

	// Values returns all registered interfaces
	Values() []interface{}
}

// registry is a struct binding name and interface such as Constructor
type registry struct {
	data            sync.Map
	overrideAllowed bool
}

// Config is a struct containing all config for registry
type Config struct {

	// OverrideAllowed allows the registry to override
	// an already registered interface by name if it is true,
	// otherwise registry will panic.
	OverrideAllowed bool
}

var (
	defaultConfig = &Config{
		OverrideAllowed: false,
	}
)

// New returns a new registry
func New(config *Config) Registry {
	if config == nil {
		config = defaultConfig
	}

	return &registry{
		data:            sync.Map{},
		overrideAllowed: config.OverrideAllowed,
	}
}

// Register registers a interface by name.
// It will panic if name corresponds to an already registered interface
// and the registry does not allow user to override the interface.
func (r *registry) Register(name string, v interface{}) error {
	if r.overrideAllowed {
		r.data.Store(name, v)
	} else {
		_, ok := r.data.LoadOrStore(name, v)
		if ok {
			return fmt.Errorf("[registry] Repeated registration key: %v", name)
		}
	}
	return nil
}

// Get returns an interface registered with the given name
func (r *registry) Get(name string) (interface{}, bool) {
	return r.data.Load(name)
}

// Range calls f sequentially for each key and value present in the registry.
// If f returns false, range stops the iteration.
func (r *registry) Range(f func(key string, value interface{}) bool) {
	r.data.Range(func(k, v interface{}) bool {
		return f(k.(string), v)
	})
}

// Keys returns the name of all registered interfaces
func (r *registry) Keys() []string {
	names := []string{}
	r.data.Range(func(k, v interface{}) bool {
		names = append(names, k.(string))
		return true
	})
	return names
}

// Values returns all registered interfaces
func (r *registry) Values() []interface{} {
	ret := []interface{}{}
	r.data.Range(func(k, v interface{}) bool {
		ret = append(ret, v)
		return true
	})
	return ret
}
