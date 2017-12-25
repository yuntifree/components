package strutil

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

//GenSalt generate 32byte random salt
func GenSalt() string {
	u := uuid.NewV4()
	return strings.Join(strings.Split(u.String(), "-"), "")
}

//GenUUID generate 36byte uuid
func GenUUID() string {
	u := uuid.NewV4()
	return u.String()
}
