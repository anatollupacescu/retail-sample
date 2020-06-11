package arbor_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/internal/arbor"
)

func TestSingle(t *testing.T) {
	t.Run("given a single test", func(t *testing.T) {
		var called int

		add := arbor.New("test", func() error {
			called++
			return nil
		})

		add.run()

		t.Run("is called once", func(t *testing.T) {
			assert.Equal(t, 1, called)
		})
	})
}

func TestDep(t *testing.T) {
	t.Run("given a test with a dep", func(t *testing.T) {
		var called bool

		dep := arbor.New("dep", func() error {
			return errors.New("bad result")
		})

		addTest := arbor.New("main", func() error {
			called = true
			return nil
		}, dep)

		addTest.run()

		t.Run("test is not runFunc when dep fails", func(t *testing.T) {
			assert.False(t, called)
			assert.Equal(t, "bad result", dep.FailReason)
			assert.Equal(t, arbor.Pending, addTest.Status)
		})
	})
}

func TestOrder(t *testing.T) {
	t.Run("given a test with a dep", func(t *testing.T) {
		var calls string

		dep := arbor.New("dep", func() error {
			calls += "dep"
			return nil
		})

		addTest := arbor.New("can add two numbers", func() error {
			calls += "main"
			return nil
		}, dep)

		addTest.run()

		t.Run("dep is ran before linking test", func(t *testing.T) {
			assert.Equal(t, "depmain", calls)
		})
	})
}

func TestDiamond(t *testing.T) {
	t.Run("given that a common dependency succeeds", func(t *testing.T) {
		var calls int

		dep := arbor.New("dep", func() error {
			calls = 1
			return nil
		})

		first := arbor.New("first", func() error {
			calls += 10
			return nil
		}, dep)

		second := arbor.New("second", func() error {
			calls += 100
			return nil
		}, dep)

		diamond := arbor.Suite("diamond", first, second)

		diamond.run()

		t.Run("it is ran exactly once", func(t *testing.T) {
			assert.Equal(t, 111, calls)
		})

		t.Run("will compile summary", func(t *testing.T) {
			summary := `[diamond] passed
  ⮡[first] passed
    ⮡[dep] passed
  ⮡[second] passed
    ⮡[dep] passed
`
			assert.Equal(t, summary, diamond.String())
		})
	})
}
