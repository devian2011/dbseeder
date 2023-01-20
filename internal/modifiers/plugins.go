package modifiers

import "errors"

var (
	ErrUnknownPlugin = errors.New("unknown plugin")
	ErrWrongType     = errors.New("cannot apply plugin - wrong type")
)

type Modifier interface {
	GetCode() string
	GetFn() ModifierFn
	GetDescription() string
}

type ModifierFn func(v any) (any, error)

type ModifierStore struct {
	store map[string]Modifier
}

func NewModifierStore() *ModifierStore {
	return &ModifierStore{
		store: map[string]Modifier{
			"bcrypt": &bCryptModifier{},
		},
	}
}

func (m ModifierStore) Apply(pluginName string, v any) (any, error) {
	if m, exists := m.store[pluginName]; exists {
		return m.GetFn()(v)
	}
	return v, ErrUnknownPlugin
}

func (m ModifierStore) List() map[string]string {
	out := make(map[string]string, 0)
	for k, md := range m.store {
		out[k] = md.GetDescription()
	}

	return out
}
