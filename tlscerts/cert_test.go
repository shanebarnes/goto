package tlscerts

import (
    "crypto/x509"
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
    notBefore, notAfter, _ := CreateCertValidity("", 365 * 24 * time.Hour)

    if _, err := CreateCert(priv, x509.ExtKeyUsageClientAuth, notBefore, notAfter, "localhost,127.0.0.1", false); err != nil {
        t.Errorf("Client certificate generation failed: %s\n", err)
    }

    if _, err := CreateCert(priv, x509.ExtKeyUsageServerAuth, notBefore, notAfter, "localhost,127.0.0.1", false); err != nil {
        t.Errorf("Server certificate generation failed: %s\n", err)
    }
}
