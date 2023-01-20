package modifiers

type bCryptModifier struct {
}

func (m *bCryptModifier) GetCode() string {
	return "bcrypt"
}

func (m *bCryptModifier) GetDescription() string {
	return "Use BCrypt for description"
}

func (m *bCryptModifier) GetFn() ModifierFn {
	return func(v any) (any, error) {
		switch v.(type) {
		case string:
			return v, nil
		default:
			return nil, ErrWrongType
		}
	}
}
