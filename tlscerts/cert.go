package tlscerts

import (
    "bytes"
    "crypto/rand"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "math/big"
    "net"
    "strings"
    "time"
)

func CreateCert(priv interface{}, extKeyUsage x509.ExtKeyUsage, notBefore, notAfter time.Time, host string, isCertAuth bool) (*bytes.Buffer, error) {
    serialNumber, err := generateSerialNumber()

    template := x509.Certificate{
        SerialNumber: serialNumber,
        Subject: pkix.Name{
            Organization: []string{"Acme Co"},
        },
        NotBefore: notBefore,
        NotAfter: notAfter,

        KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage: []x509.ExtKeyUsage{extKeyUsage},
        BasicConstraintsValid: true,
    }

    hosts := strings.Split(host, ",")
    for _, h := range hosts {
        if ip := net.ParseIP(h); ip != nil {
            template.IPAddresses = append(template.IPAddresses, ip)
        } else {
            template.DNSNames = append(template.DNSNames, h)
        }
    }

    if isCertAuth {
        template.IsCA = true
        template.KeyUsage |= x509.KeyUsageCertSign
    }

    var certOut *bytes.Buffer = new(bytes.Buffer)
    var derBytes []byte
    if derBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, getPublicKey(priv), priv); err == nil {
        pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
    }

    return certOut, err
}

func CreateCertValidity(validFrom string, validFor time.Duration) (time.Time, time.Time, error) {
    var err error = nil
    var notBefore, notAfter time.Time

    if len(validFrom) == 0 {
        notBefore = time.Now()
    } else {
        notBefore, err = time.Parse("Jan 2 15:04:05 2006", validFrom)
    }

    if err == nil {
        notAfter = notBefore.Add(validFor)
    }

    return notBefore, notAfter, err
}

func generateSerialNumber() (*big.Int, error) {
    limit := new(big.Int).Lsh(big.NewInt(1), 128)
    return rand.Int(rand.Reader, limit)
}
