package prik

import (
	"fmt"
	"sync"
)

type DisposeFn func()
type Dispose func(fn DisposeFn)

type AnyFactory interface{}
type Factory func() (AnyFactory, DisposeFn)
type Factories map[string]Factory

type Context struct {
	factories   Factories
	disposables []DisposeFn
}

func Exec[T any](factory *Factory) (*T, DisposeFn, error) {
	value, dispose := (*factory)()
	tf, ok := value.(T); if !ok {
		return nil, nil, fmt.Errorf("Failed to cast factory %s to type %T", value, *new(T))
	}

	return &tf, dispose, nil
}

func Resolve[T any](ctx *Context, name string) *T {
	_factory, ok := ctx.factories[name]; if !ok {
		panic(fmt.Errorf("No factory found: %s", name))
	}

	factory, dispose, err := Exec[T](&_factory); if err != nil {
		panic(err)
	}

	ctx.disposables = append(ctx.disposables, dispose)
	return factory
}

func ResolveOr[T any](ctx *Context, name string) (*T, error) {
	_factory, ok := ctx.factories[name]; if !ok {
		return nil, fmt.Errorf("No factory found: %s", name)
	}

	factory, dispose, err := Exec[T](&_factory); if err != nil {
		return nil, err
	}

	ctx.disposables = append(ctx.disposables, dispose)
	return factory, nil
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
		instance AnyFactory
		dispose  DisposeFn
		once sync.Once
		mutex sync.Mutex
	)

	return func() (AnyFactory, DisposeFn) {
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
