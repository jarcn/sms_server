package utils

import (
	"reflect"
	"testing"
)

type User struct {
	Name string
	Age  int
}

func NoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func Error(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func Equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("not equal:"+"expected: %v actual: %+v ", expected, actual)
	}
}

func getCache() *RedisCache {
	cache, err := New(Options{
		Prefix: "chenjia_",
	})
	if err != nil {
		panic(err)
	}
	return cache
}

func TestGetSet(t *testing.T) {
	var err error
	cache := getCache()
	err = cache.Set("test", "23", 0)
	NoError(t, err)
	getInt, err := cache.GetInt("age")
	NoError(t, err)
	Equal(t, 23, getInt)
}
