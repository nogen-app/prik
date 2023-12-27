package nogencontext

import "testing"


// MockDBFactory is an example factory for a "db" key.
func MockDBFactory() (interface{}, DisposeFn) {
	// Mocked value for testing purposes
	mockDB := "MockedDB"

	// Mocked DisposeFn (no-op for this example)
	dispose := func() {}

	return mockDB, dispose
}

func TestResolveDependency(t *testing.T) {
	factories := Factories{
		"db": MockDBFactory,
	}

	ctx := CreateContext(factories)

	value, err := ctx.Resolve("db")

	if err != nil {
		t.Fatalf("Unexpected error while resolving key 'db': %v", err)
	}

	// Check that the resolved value matches the expected mocked value
	expectedValue := "MockedDB"
	if value != expectedValue {
		t.Errorf("Expected resolved value to be '%s', got '%v'", expectedValue, value)
	}
}

func TestWithoutShared(t *testing.T) {
	var	callCount int = 0

	mockFactory := func() (interface{}, DisposeFn) {
		callCount++
		return "MockedValue", func() {}
	}

	factories := Factories{
		"db": mockFactory,
	}

	ctx := CreateContext(factories)

	for i := 0; i < 5; i++ {
		value, err := ctx.Resolve("db")
		
		if err != nil {
			t.Fatalf("Unexpected error while resolving key 'db': %v", err)
		}

		// Check that the resolved value matches the expected mocked value
		expectedValue := "MockedValue"
		if value != expectedValue {
			t.Errorf("Expected resolved value to be '%s', got '%v'", expectedValue, value)
		}
	}

	if callCount != 5 {
		t.Errorf("Expected the mock factory to be called five times, got %d times", callCount)
	}
}

func TestShared(t *testing.T) {
	var	callCount int = 0

	mockFactory := func() (interface{}, DisposeFn) {
		callCount++
		return "MockedValue", func() {}
	}

	factories := Factories{
		"db": Shared(mockFactory),
	}

	ctx := CreateContext(factories)

	for i := 0; i < 5; i++ {
		value, err := ctx.Resolve("db")
		
		if err != nil {
			t.Fatalf("Unexpected error while resolving key 'db': %v", err)
		}

		// Check that the resolved value matches the expected mocked value
		expectedValue := "MockedValue"
		if value != expectedValue {
			t.Errorf("Expected resolved value to be '%s', got '%v'", expectedValue, value)
		}
	}

	if callCount != 1 {
		t.Errorf("Expected the mock factory to be called once, got %d times", callCount)
	}
}
