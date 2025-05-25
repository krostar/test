package compare

import (
	"os"
	"slices"
	"testing"

	"github.com/matryer/is"
	"github.com/stretchr/testify/assert"
	gotestassert "gotest.tools/v3/assert"

	"github.com/krostar/test"
	"github.com/krostar/test/double"
)

func Benchmark_SimpleAssertion_Success(b *testing.B) {
	fakeT := double.NewFake()

	b.Run("krostar/test", func(b *testing.B) {
		for b.Loop() {
			test.Assert(fakeT, true)
		}
	})

	b.Run("testify", func(b *testing.B) {
		adapter := &testifyAdapter{fakeT}

		for b.Loop() {
			assert.True(adapter, true)
		}
	})

	b.Run("matryer/is", func(b *testing.B) {
		is := is.New(fakeT)

		for b.Loop() {
			is.True(true)
		}
	})

	b.Run("gotesttools", func(b *testing.B) {
		adapter := &gotestAdapter{fakeT}

		for b.Loop() {
			gotestassert.Check(adapter, true)
		}
	})
}

func Benchmark_SimpleAssertion_Failing(b *testing.B) {
	fakeT := double.NewFake()

	b.Run("krostar/test", func(b *testing.B) {
		for b.Loop() {
			test.Assert(fakeT, false)
		}
	})

	b.Run("testify", func(b *testing.B) {
		adapter := &testifyAdapter{fakeT}

		for b.Loop() {
			assert.True(adapter, false)
		}
	})

	b.Run("matryer/is", func(b *testing.B) {
		stdout := os.Stdout
		os.Stdout, _ = os.Open(os.DevNull)
		b.Cleanup(func() { os.Stdout = stdout })

		is := is.New(fakeT)

		for b.Loop() {
			is.True(false)
		}
	})

	b.Run("gotesttools", func(b *testing.B) {
		adapter := &gotestAdapter{fakeT}

		for b.Loop() {
			gotestassert.Check(adapter, false)
		}
	})
}

func Benchmark_SimpleEquality_Success(b *testing.B) {
	fakeT := double.NewFake()
	user1 := struct {
		ID    int
		Name  string
		Email string
		Tags  []string
	}{
		ID:    123,
		Name:  "John Doe",
		Email: "john@example.com",
		Tags:  []string{"admin", "user"},
	}
	user2 := user1

	b.Run("krostar/test", func(b *testing.B) {
		for b.Loop() {
			test.Assert(fakeT, user1.ID == user2.ID)
			test.Assert(fakeT, user1.Name == user2.Name)
			test.Assert(fakeT, user1.Email == user2.Email)
			test.Assert(fakeT, slices.Equal(user1.Tags, user2.Tags))
		}
	})

	b.Run("testify", func(b *testing.B) {
		adapter := &testifyAdapter{fakeT}

		for b.Loop() {
			assert.Equal(adapter, user1.ID, user2.ID)
			assert.Equal(adapter, user1.Name, user2.Name)
			assert.Equal(adapter, user1.Email, user2.Email)
			assert.Equal(adapter, user1.Tags, user2.Tags)
		}
	})

	b.Run("matryer/is", func(b *testing.B) {
		is := is.New(fakeT)

		for b.Loop() {
			is.Equal(user1.ID, user2.ID)
			is.Equal(user1.Name, user2.Name)
			is.Equal(user1.Email, user2.Email)
			is.Equal(user1.Tags, user2.Tags)
		}
	})

	b.Run("gotesttools", func(b *testing.B) {
		adapter := &gotestAdapter{fakeT}

		for b.Loop() {
			gotestassert.Equal(adapter, user1.ID, user2.ID)
			gotestassert.Equal(adapter, user1.Name, user2.Name)
			gotestassert.Equal(adapter, user1.Email, user2.Email)
			gotestassert.Check(adapter, slices.Equal(user1.Tags, user2.Tags))
		}
	})
}

func Benchmark_SimpleEquality_Failure(b *testing.B) {
	fakeT := double.NewFake()
	user1, user2 := struct {
		ID    int
		Name  string
		Email string
		Tags  []string
	}{
		ID:    123,
		Name:  "John",
		Email: "john@example.com",
		Tags:  []string{"regular", "user"},
	}, struct {
		ID    int
		Name  string
		Email string
		Tags  []string
	}{
		ID:    321,
		Name:  "Alice",
		Email: "alice@example.com",
		Tags:  []string{"admin", "user"},
	}

	b.Run("krostar/test", func(b *testing.B) {
		for b.Loop() {
			test.Assert(fakeT, user1.ID == user2.ID)
			test.Assert(fakeT, user1.Name == user2.Name)
			test.Assert(fakeT, user1.Email == user2.Email)
			test.Assert(fakeT, slices.Equal(user1.Tags, user2.Tags))
		}
	})

	b.Run("testify", func(b *testing.B) {
		adapter := &testifyAdapter{fakeT}

		for b.Loop() {
			assert.Equal(adapter, user1.ID, user2.ID)
			assert.Equal(adapter, user1.Name, user2.Name)
			assert.Equal(adapter, user1.Email, user2.Email)
			assert.Equal(adapter, user1.Tags, user2.Tags)
		}
	})

	b.Run("matryer/is", func(b *testing.B) {
		stdout := os.Stdout
		os.Stdout, _ = os.Open(os.DevNull)
		b.Cleanup(func() { os.Stdout = stdout })

		is := is.New(fakeT)

		for b.Loop() {
			is.Equal(user1.ID, user2.ID)
			is.Equal(user1.Name, user2.Name)
			is.Equal(user1.Email, user2.Email)
			is.Equal(user1.Tags, user2.Tags)
		}
	})

	b.Run("gotesttools", func(b *testing.B) {
		adapter := &gotestAdapter{fakeT}

		for b.Loop() {
			gotestassert.Equal(adapter, user1.ID, user2.ID)
			gotestassert.Equal(adapter, user1.Name, user2.Name)
			gotestassert.Equal(adapter, user1.Email, user2.Email)
			gotestassert.Check(adapter, slices.Equal(user1.Tags, user2.Tags))
		}
	})
}

// testifyAdapter adapts double.TestingT to work with testify's assert.TestingT interface
type testifyAdapter struct{ double.TestingT }

func (t *testifyAdapter) Errorf(format string, args ...interface{}) {
	t.TestingT.Logf(format, args...)
	t.TestingT.Fail()
}

// gotestAdapter adapts double.TestingT to work with gotest.tools' assert.TestingT interface
type gotestAdapter struct{ double.TestingT }

func (t *gotestAdapter) Log(args ...interface{}) { t.TestingT.Logf("%v", args) }
