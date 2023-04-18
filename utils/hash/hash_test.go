package hash

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test001(t *testing.T) {
	salt := "abc"
	passwd := "1212123"
	hashPasswd, _ := HashPassword(passwd + salt)
	fmt.Println(passwd, hashPasswd)

	ret := CheckPasswordHash(passwd+salt, hashPasswd)
	assert.True(t, ret)
}
