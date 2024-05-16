package prik

import "testing"


// MockDBFactory is an example factory for a "db" key.
func MockDBFactory() (AnyFactory, DisposeFn) {
	// Mocked value for testing purposes
	mockDB := "MockedDB"

	// Mocked DisposeFn (no-op for this example)
	dispose := func() {}

	return mockDB, dispose
}

func TestFunctionResolveOrDependency(t *testing.T) {
	factories := Factories{
		"db": MockDBFactory,
	}

	ctx := CreateContext(factories)

	// Success case
	value, err := ResolveOr[string](ctx, "db")

	if err != nil {
		t.Fatalf("Unexpected error while resolving key 'db': %v", err)
	}

	// Check that the resolved value matches the expected mocked value
	expectedValue := "MockedDB"
	if *value != expectedValue {
		t.Errorf("Expected resolved value to be '%s', got '%v'", expectedValue, value)
	}

	// Failure case (non-existent key)
	value, err = ResolveOr[string](ctx, "non-existent-key")
	if err == nil {
		t.Fatalf("Expected error while resolving non-existent key, got nil")
	}
	if err.Error() != "No factory found: non-existent-key" {
		t.Fatalf("Expected error message to be 'No factory found: non-existent-key', got '%v'", err)
	}
	if value != nil {
		t.Fatalf("Expected resolved value to be nil, got '%v'", value)
	}
	
	// Failure case (type assertion error)
	failvalue, failerr := ResolveOr[bool](ctx, "db")
	if failerr == nil {
		t.Fatalf("Expected error while resolving key 'db' with wrong type, got nil")
	}
	if failerr.Error() != "Failed to cast factory MockedDB to type bool" {
		t.Fatalf("Expected error message to be 'Failed to cast factory MockedDB to type bool', got '%v'", failerr)
	}
	if failvalue != nil {
		t.Fatalf("Expected resolved value to be nil, got '%v'", failvalue)
	}
}

func TestFunctionResolveDependency(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	factories := Factories{
		"db": MockDBFactory,
	}

	ctx := CreateContext(factories)

	// Success case
	value := Resolve[string](ctx, "db")

	// Check that the resolved value matches the expected mocked value
	expectedValue := "MockedDB"
	if *value != expectedValue {
		t.Errorf("Expected resolved value to be '%s', got '%v'", expectedValue, value)
	}

	// Failure case (non-existent key)\
	Resolve[string](ctx, "non-existent-key")

	// Failure case (type assertion error)
	Resolve[bool](ctx, "db")
}

func TestWithoutShared(t *testing.T) {
	var	callCount int = 0

	mockFactory := func() (AnyFactory, DisposeFn) {
		callCount++
		return "MockedValue", func() {}
	}

	factories := Factories{
		"db": mockFactory,
	}

	ctx := CreateContext(factories)

	for i := 0; i < 5; i++ {
		value, err := ResolveOr[string](ctx, "db")
		
		if err != nil {
			t.Fatalf("Unexpected error while resolving key 'db': %v", err)
		}

		// Check that the resolved value matches the expected mocked value
		expectedValue := "MockedValue"
		if *value != expectedValue {
			t.Errorf("Expected resolved value to be '%s', got '%v'", expectedValue, value)
		}
	}

	if callCount != 5 {
		t.Errorf("Expected the mock factory to be called five times, got %d times", callCount)
	}
}

func TestShared(t *testing.T) {
	var	callCount int = 0

	mockFactory := func() (AnyFactory, DisposeFn) {
		callCount++
		return "MockedValue", func() {}
	}

	factories := Factories{
		"db": Shared(mockFactory),
	}

	ctx := CreateContext(factories)

	for i := 0; i < 5; i++ {
		value, err := ResolveOr[string](ctx, "db")
		
		if err != nil {
			t.Fatalf("Unexpected error while resolving key 'db': %v", err)
		}

		// Check that the resolved value matches the expected mocked value
		expectedValue := "MockedValue"
		if *value != expectedValue {
			t.Errorf("Expected resolved value to be '%s', got '%v'", expectedValue, value)
		}
	}

	if callCount != 1 {
		t.Errorf("Expected the mock factory to be called once, got %d times", callCount)
	}
}
