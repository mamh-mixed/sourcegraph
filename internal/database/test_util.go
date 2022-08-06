package database

import (
	"context"

	"github.com/sourcegraph/sourcegraph/internal/encryption"
)

func MockEmailExistsErr() error {
	return errCannotCreateUser{errorCodeEmailExists}
}

func MockUsernameExistsErr() error {
	return errCannotCreateUser{errorCodeEmailExists}
}

func strptr(s string) *string {
	return &s
}

func boolptr(b bool) *bool {
	return &b
}

func testEncryptionKeyID(key encryption.Key) string {
	v, _ := key.Version(context.Background())
	return v.JSON()
}
