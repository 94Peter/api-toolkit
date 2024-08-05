package auth_test

import (
	"testing"

	"github.com/94peter/api-toolkit/auth"
)

func TestJwtConf_IsRsaKeysExist(t *testing.T) {
	// Test case 1: Both privateKeyFile and publicKeyFile exist
	j := &auth.JwtConf{
		PrivateKeyFile: "./privateKeyFile",
		PublicKeyFile:  "./publicKeyFile",
	}
	if !j.IsRsaKeysExist() {
		err := j.GenerateRsaKeys(2048)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}

	if !j.IsRsaKeysExist() {
		t.Errorf("Expected true, got false")
	}

	// Test case 2: PrivateKeyFile does not exist
	j = &auth.JwtConf{
		PrivateKeyFile: "/path/to/nonExistentFile",
		PublicKeyFile:  "/path/to/publicKeyFile",
	}
	if j.IsRsaKeysExist() {
		t.Errorf("Expected false, got true")
	}

	// Test case 3: PublicKeyFile does not exist
	j = &auth.JwtConf{
		PrivateKeyFile: "/path/to/privateKeyFile",
		PublicKeyFile:  "/path/to/nonExistentFile",
	}
	if j.IsRsaKeysExist() {
		t.Errorf("Expected false, got true")
	}
}
