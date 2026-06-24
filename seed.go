package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// Seed format: "<base64(data)>.<4-hex signature>"
// Signature: first 2 bytes of SHA256( encodedData + salt ), hex encoded.
func BuildSignedSeed(rawData string, salt string) string {
	encoded := base64.RawStdEncoding.EncodeToString([]byte(rawData))
	sig := seedSignature(encoded, salt)
	return encoded + "." + sig
}

// VerifySignedSeed validates the signature, then returns the decoded seed data.
func VerifySignedSeed(seed string, salt string) (string, error) {
	encoded, sig, ok := strings.Cut(seed, ".")
	if !ok || encoded == "" || sig == "" {
		return "", fmt.Errorf("invalid seed format")
	}

	expected := seedSignature(encoded, salt)
	if !strings.EqualFold(sig, expected) {
		return "", errors.New(strings.Join([]string{
			"[CRITICAL] SESSION INTEGRITY BREACH DETECTED.",
			"[WARN] IP LOGGED. TRACE INITIALIZED.",
			"[SYSTEM] RELOADING LAST SECURE SNAPSHOT...",
		}, "\n"))
	}

	raw, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("invalid seed data")
	}
	return string(raw), nil
}

func seedSignature(encodedData string, salt string) string {
	sum := sha256.Sum256([]byte(encodedData + salt))
	return hex.EncodeToString(sum[:2])
}

