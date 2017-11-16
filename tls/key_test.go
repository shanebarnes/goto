package tls

import (
    "testing"
)

func TestTlsGeneratePrivateKey(t *testing.T) {
    algs := []string{"1024", "2048", "3072", "P224", "P256", "P384", "P521"}

    for _, alg := range algs {
        _, err := GeneratePrivateKey(alg)

        if err != nil {
            t.Errorf("Key generation failed: %s\n", err)
        }
    }
}

func TestTlsEncodePrivateKey(t *testing.T) {
    algs := []string{"1024", "2048", "3072", "P224", "P256", "P384", "P521"}

    for _, alg := range algs {
        if priv, err := GeneratePrivateKey(alg); err == nil {
            if _, err := EncodePrivateKey(priv); err != nil {
                t.Errorf("Key encoding failed: %s\n", err)
            }
        } else {
            t.Errorf("Key generation failed: %s\n", err)
        }
    }
}
