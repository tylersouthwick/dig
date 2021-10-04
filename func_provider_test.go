package dig

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type MyInterface interface {
	Get() string
}

type FooBar struct {
	Hello string
}

type FooBarWithInterface struct {
	MyInterface
	Hello string
}

func (i FooBarWithInterface) Get() string {
	return i.Hello
}

type FooBar2 struct {
	Hello  string
	FooBar *FooBar
}

func TestFuncProvider_SingleField(t *testing.T) {
	d := New()
	require.NoError(t, d.Provide(func() string {
		return "hello world"
	}), "provide failed")
	require.NoError(t, d.Provide(funcProvider(reflect.TypeOf(FooBar{}))), "provide failed")

	a := assert.New(t)
	require.NoError(t, d.Invoke(func(foobar *FooBar) {
		a.Equal("hello world", foobar.Hello)
	}), "invoke failed")
}

func TestFuncProvider_MultipleFields(t *testing.T) {
	d := New()
	require.NoError(t, d.Provide(func() string {
		return "hello world"
	}), "provide failed")
	require.NoError(t, d.Provide(funcProvider(reflect.TypeOf(FooBar{}))), "provide failed")
	require.NoError(t, d.Provide(funcProvider(reflect.TypeOf(FooBar2{}))), "provide failed")

	a := assert.New(t)
	require.NoError(t, d.Invoke(func(foobar *FooBar2) {
		a.Equal("hello world", foobar.Hello)
		a.Equal("hello world", foobar.FooBar.Hello)
	}), "Invoke failed")
}

func TestFuncProvider_Interface(t *testing.T) {
	d := New()
	require.NoError(t, d.Provide(func() string {
		return "hello world"
	}), "provide failed")
	require.NoError(t, d.Provide(funcProvider(reflect.TypeOf(FooBarWithInterface{}))), "provide failed")

	a := assert.New(t)
	require.NoError(t, d.Invoke(func(foobar *FooBarWithInterface) {
		a.Equal("hello world", foobar.Get())
	}), "Invoke failed")
}
