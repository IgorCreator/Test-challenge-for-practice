package crypto

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	key := "6aQqE17SgkXypLNtAsfbntSLpl7kMP/qdRQThhCtdwE="
	cipher, err := NewCipherFromBase64(key)
	if err != nil {
		t.Fatalf("cipher: %v", err)
	}
	secret := "breeder@example.com"
	enc, nonce, err := cipher.Encrypt(secret)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	out, err := cipher.Decrypt(enc, nonce)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if out != secret {
		t.Fatalf("expected %q, got %q", secret, out)
	}
}

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("strong_pass_123")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	ok, err := VerifyPassword("strong_pass_123", hash)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !ok {
		t.Fatalf("expected password to verify")
	}
}
