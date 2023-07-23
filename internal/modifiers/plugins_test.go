package modifiers

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

type mOne struct {
	code string
	fn   ModifierFn
}

func (m *mOne) GetCode() string {
	return m.code
}

func (m *mOne) GetDescription() string {
	return ""
}

func (m *mOne) GetFn() ModifierFn {
	return m.fn
}

func TestModifierStore_Apply(t *testing.T) {
	mStore := &ModifierStore{
		store: map[string]Modifier{
			"one": &mOne{
				code: "one",
				fn: func(v any) (any, error) {
					return "one " + v.(string), nil
				},
			},
		},
	}

	expected := "one hello"
	actual, _ := mStore.Apply("one", "hello")
	assert.Equal(t, expected, actual)
}

func TestModifierStore_ApplyList(t *testing.T) {
	mStore := &ModifierStore{
		store: map[string]Modifier{
			"one": &mOne{
				code: "one",
				fn: func(v any) (any, error) {
					return "one " + v.(string), nil
				},
			},
			"two": &mOne{
				code: "two",
				fn: func(v any) (any, error) {
					return "two " + v.(string), nil
				},
			},
		},
	}

	expected := "two one hello"
	actual, _ := mStore.ApplyList([]string{"one", "two"}, "hello")
	assert.Equal(t, expected, actual)
}
