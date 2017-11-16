package tls

import (
    "testing"
    "time"
)

func TestTlscertGenerateSerialNumber(t *testing.T) {
    _, err := generateSerialNumber()

    if err != nil {
        t.Errorf("Serial number generation failed: %s\n", err)
    }
}

func TestTlscertCreateCert(t *testing.T) {
    priv, _ := GeneratePrivateKey("2048")
    notBefore, notAfter, _ := createCertValidity("", 365 * 24 * time.Hour)

    if _, err := CreateCert(priv, notBefore, notAfter, "localhost,127.0.0.1", false); err != nil {
        t.Errorf("Certificate generation failed: %s\n", err)
    }
}
