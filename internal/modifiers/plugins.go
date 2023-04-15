package modifiers

import "errors"

var (
	// ErrUnknownPlugin unknown plugin error
	ErrUnknownPlugin = errors.New("unknown plugin")
	// ErrWrongType wrong plugin type for apply
	ErrWrongType = errors.New("cannot apply plugin - wrong type")
)

// Modifier interface
type Modifier interface {
	GetCode() string
	GetFn() ModifierFn
	GetDescription() string
}

type ModifierFn func(v any) (any, error)

// ModifierStore instace store plugins
type ModifierStore struct {
	store map[string]Modifier
}

// NewModifierStore create modifier store instance
func NewModifierStore() *ModifierStore {
	return &ModifierStore{
		store: map[string]Modifier{
			"bcrypt": &bCryptModifier{},
		},
	}
}

// Apply apply plugin to value
func (m ModifierStore) Apply(pluginName string, v any) (any, error) {
	if m, exists := m.store[pluginName]; exists {
		return m.GetFn()(v)
	}
	return v, ErrUnknownPlugin
}

// ApplyList apply list of plugins to value
func (m *ModifierStore) ApplyList(pluginNameList []string, v any) (any, error) {
	var err error
	for _, pluginName := range pluginNameList {
		v, err = m.Apply(pluginName, v)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

// List show plugins list
func (m ModifierStore) List() map[string]string {
	out := make(map[string]string, 0)
	for k, md := range m.store {
		out[k] = md.GetDescription()
	}

	return out
}
