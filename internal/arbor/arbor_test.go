package arbor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingle(t *testing.T) {
	t.Run("given a single test", func(t *testing.T) {
		var called int

		add := New("can add two numbers", func() error {
			called++
			return nil
		})

		add.Run()

		t.Run("is called once", func(t *testing.T) {
			assert.Equal(t, 1, called)
		})
	})
}

func TestDep(t *testing.T) {
	t.Run("given a test with a dep", func(t *testing.T) {
		var called bool

		dep := New("dep", func() error {
			return errors.New("bad result")
		})

		addTest := New("can add two numbers", func() error {
			called = true
			return nil
		}, dep)

		addTest.Run()

		t.Run("test is not run when dep fails", func(t *testing.T) {
			assert.False(t, called)
			assert.Equal(t, "bad result", dep.failReason)
			assert.Equal(t, pending, addTest.status)
		})
	})
}

func TestOrder(t *testing.T) {
	t.Run("given a test with a dep", func(t *testing.T) {
		var calls string

		dep := New("dep", func() error {
			calls += "dep"
			return nil
		})

		addTest := New("can add two numbers", func() error {
			calls += "main"
			return nil
		}, dep)

		addTest.Run()

		t.Run("dep is ran before linking test", func(t *testing.T) {
			assert.Equal(t, "depmain", calls)
		})
	})
}

func TestDiamond(t *testing.T) {
	t.Run("given that a common dependency succeeds", func(t *testing.T) {
		var calls int

		dep := New("dep", func() error {
			calls = 1
			return nil
		})

		first := New("first", func() error {
			calls += 10
			return nil
		}, dep)

		second := New("second", func() error {
			calls += 100
			return nil
		}, dep)

		diamond := New("diamond", func() error { return nil }, first, second)

		diamond.Run()

		t.Run("it is ran exactly once", func(t *testing.T) {
			assert.Equal(t, 111, calls)
		})

		t.Run("will compile summary", func(t *testing.T) {
			summary := `[diamond] ok
	⮡[first] ok
		⮡[dep] ok
	⮡[second] ok
		⮡[dep] ok
`
			assert.Equal(t, summary, diamond.String())
		})
	})
}
