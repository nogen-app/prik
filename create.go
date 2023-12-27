package nogencontext

import (
	"fmt"
	"sync"
)

type DisposeFn func()
type Dispose func(fn DisposeFn)

type Factory func() (interface{}, DisposeFn)
type Factories map[string]Factory

type ContextInterface interface {
	Resolve(name string) (interface{}, error)
	Dispose()
} 

type Context struct {
	//	ContextInterface

	factories Factories
	disposables []DisposeFn
}

func (c *Context) Resolve(name string) (interface{}, error) {
	factory, ok := c.factories[name]; if !ok {
		return nil, fmt.Errorf("No factory found: %s", name)
	}
	value, dispose := factory()
	c.disposables = append(c.disposables, dispose)
	return value, nil
}

func (c *Context) Dispose() {
	for _, d := range c.disposables { d() }
}

func CreateContext(factories Factories) *Context {
	return &Context{
		factories: factories,
		disposables: make([]DisposeFn, 0),
	}
}

func Shared(factory Factory) Factory {
	var (
		instance interface{}
		dispose DisposeFn
		once sync.Once
		mutex sync.Mutex
	)

	return func() (interface{}, DisposeFn) {
		once.Do(func() {
			value, disposefn := factory()
			instance = value
			dispose = disposefn
		})
		mutex.Lock()
		defer mutex.Unlock()

		return instance, dispose
	}
}

// func Shared(factory Factory) Factory {
// 	var (
// 		instance interface{}
// 		once sync.Once
// 		mutex sync.Mutex
// 	)

// 	return func() (interface{}, error) {
// 		once.Do(func() {
// 			value, disposefn := factory()
// 			instance = value
// 		})
// 		mutex.Lock()
// 		defer mutex.Unlock()

// 		return instance, nil
// 	}
// }

// func make(factories Factories) Context {
// 	disposables := map[string]DisposeFn{}

// 	ctx := Context{
// 		resolve: func(name string) (interface{}, error) {
// 			f := factories[name]
// 			return f(dispose DisposeFn) {
// 				dispodisposables[name] = dispose
// 			})
// 			//f(
// 		},
// 		dispose: func() {
// 			for x, d := range disposables { d() }
// 		},
// 	}
// }
