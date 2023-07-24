package modifiers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestBCryptModifier_GetFn(t *testing.T) {
	o := &bCryptModifier{}
	_, err := o.GetFn()(123)
	assert.Equal(t, ErrWrongType, err)
	password := "pa$$w0rd"
	actual, err := o.GetFn()(password)
	if err != nil {
		t.Errorf("BCrypt must support string values")
	}
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(actual.(string)), []byte(password)))
}
