package modifiers

import "golang.org/x/crypto/bcrypt"

type bCryptModifier struct {
}

func (m *bCryptModifier) GetCode() string {
	return "bcrypt"
}

func (m *bCryptModifier) GetDescription() string {
	return "Use BCrypt for modify string or []byte fields"
}

func (m *bCryptModifier) GetFn() ModifierFn {
	return func(v any) (any, error) {
		switch v.(type) {
		case string:
			result, err := bcrypt.GenerateFromPassword([]byte(v.(string)), bcrypt.DefaultCost)
			return string(result), err
		case []byte:
			result, err := bcrypt.GenerateFromPassword(v.([]byte), bcrypt.DefaultCost)
			return string(result), err
		default:
			return nil, ErrWrongType
		}
	}
}
