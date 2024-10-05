package encryption

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	NewEncryptionController()
}

func TestDecryptionPassword(t *testing.T) {
	setup()
	mockEncryptedPassword := "U2FsdGVkX13Z"
	decryptedPassword, err := EncryptionControllerObj.DecryptPassword(mockEncryptedPassword)

	assert.Nil(t, err)
	assert.Equal(t, decryptedPassword, "password")
}

func TestBase64DecodingError(t *testing.T) {
	setup()
	mockEncryptedPassword := "U2FsdGVkX13Z"
	decryptedPassword, err := EncryptionControllerObj.DecryptPassword(mockEncryptedPassword)

	assert.NotNil(t, err)
	assert.Equal(t, decryptedPassword, "")
}