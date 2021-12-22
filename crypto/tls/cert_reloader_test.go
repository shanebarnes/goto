package tls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestCertAndKey(certFile, keyFile string, notBefore, notAfter time.Time) error {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2021),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}

	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return err
	}

	caPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(certFile, caPEM.Bytes(), 0644); err != nil {
		return err
	}
	return ioutil.WriteFile(keyFile, caPrivKeyPEM.Bytes(), 0644)
}

func destroyTestCertAndKeyFile(certFile, keyFile *os.File) {
	certFile.Close()
	keyFile.Close()
	os.Remove(certFile.Name())
	os.Remove(keyFile.Name())
}

func TestNewCertificateReloader_BadCertFile(t *testing.T) {
	certFile, cErr := ioutil.TempFile("", "cert.*.pem")
	assert.NoError(t, cErr)
	keyFile, kErr := ioutil.TempFile("", "key.*.pem")
	assert.NoError(t, kErr)
	defer destroyTestCertAndKeyFile(certFile, keyFile)

	assert.NoError(t, createTestCertAndKey(certFile.Name(), keyFile.Name(), time.Now(), time.Now().Add(time.Second)))
	reloader, err := NewCertificateReloader("", keyFile.Name())
	assert.Nil(t, reloader)
	assert.Error(t, err)
}

func TestNewCertificateReloader_BadKeyFile(t *testing.T) {
	certFile, cErr := ioutil.TempFile("", "cert.*.pem")
	assert.NoError(t, cErr)
	keyFile, kErr := ioutil.TempFile("", "key.*.pem")
	assert.NoError(t, kErr)
	defer destroyTestCertAndKeyFile(certFile, keyFile)

	assert.NoError(t, createTestCertAndKey(certFile.Name(), keyFile.Name(), time.Now(), time.Now().Add(time.Second)))
	reloader, err := NewCertificateReloader(certFile.Name(), "")
	assert.Nil(t, reloader)
	assert.Error(t, err)
}

func TestNewCertificateReloader(t *testing.T) {
	certFile, cErr := ioutil.TempFile("", "cert.*.pem")
	assert.NoError(t, cErr)
	keyFile, kErr := ioutil.TempFile("", "key.*.pem")
	assert.NoError(t, kErr)
	defer destroyTestCertAndKeyFile(certFile, keyFile)

	assert.NoError(t, createTestCertAndKey(certFile.Name(), keyFile.Name(), time.Now(), time.Now().Add(time.Second)))
	reloader, err := NewCertificateReloader(certFile.Name(), keyFile.Name())
	assert.NotNil(t, reloader)
	assert.NoError(t, err)
	assert.False(t, reloader.certExpiry.IsZero())
	assert.NotNil(t, reloader.certTls)
	assert.Equal(t, certFile.Name(), reloader.certFile)
	assert.NotNil(t, reloader.closeCtx)
	assert.NotNil(t, reloader.closeFunc)
	assert.Equal(t, keyFile.Name(), reloader.keyFile)
	assert.NotNil(t, reloader.mu)

	reloader.Close()
}

func TestGetCertificateFunc(t *testing.T) {
	certFile, cErr := ioutil.TempFile("", "cert.*.pem")
	assert.NoError(t, cErr)
	keyFile, kErr := ioutil.TempFile("", "key.*.pem")
	assert.NoError(t, kErr)
	defer destroyTestCertAndKeyFile(certFile, keyFile)

	before := time.Now()
	after := before.Add(time.Second * 10)
	assert.NoError(t, createTestCertAndKey(certFile.Name(), keyFile.Name(), before, after))
	reloader, _ := NewCertificateReloader(certFile.Name(), keyFile.Name())

	f := reloader.GetCertificateFunc()
	cert, err := f(nil)
	assert.NotNil(t, cert)
	assert.NoError(t, err)

	certTest1, _ := tls.LoadX509KeyPair(certFile.Name(), keyFile.Name())
	assert.Equal(t, certTest1.Certificate, cert.Certificate)
	assert.Equal(t, certTest1.OCSPStaple, cert.OCSPStaple)
	assert.Equal(t, certTest1.PrivateKey, cert.PrivateKey)
	assert.Equal(t, certTest1.SignedCertificateTimestamps, cert.SignedCertificateTimestamps)
	assert.Equal(t, certTest1.SupportedSignatureAlgorithms, cert.SupportedSignatureAlgorithms)

	before = time.Now()
	after = before.Add(time.Hour * 24)
	assert.NoError(t, createTestCertAndKey(certFile.Name(), keyFile.Name(), before, after))
	time.Sleep(time.Second * 15)

	cert, err = f(nil)
	assert.NotNil(t, cert)
	assert.NoError(t, err)

	certTest2, _ := tls.LoadX509KeyPair(certFile.Name(), keyFile.Name())
	assert.Equal(t, certTest2.Certificate, cert.Certificate)
	assert.Equal(t, certTest2.OCSPStaple, cert.OCSPStaple)
	assert.Equal(t, certTest2.PrivateKey, cert.PrivateKey)
	assert.Equal(t, certTest2.SignedCertificateTimestamps, cert.SignedCertificateTimestamps)
	assert.Equal(t, certTest2.SupportedSignatureAlgorithms, cert.SupportedSignatureAlgorithms)
	assert.NotEqual(t, certTest2.Certificate, certTest1.Certificate)
	assert.NotEqual(t, certTest2.PrivateKey, certTest1.PrivateKey)

	reloader.Close()
}

func TestReloader_Close(t *testing.T) {
	certFile, cErr := ioutil.TempFile("", "cert.*.pem")
	assert.NoError(t, cErr)
	keyFile, kErr := ioutil.TempFile("", "key.*.pem")
	assert.NoError(t, kErr)
	defer destroyTestCertAndKeyFile(certFile, keyFile)

	assert.NoError(t, createTestCertAndKey(certFile.Name(), keyFile.Name(), time.Now(), time.Now().Add(time.Second)))
	reloader, _ := NewCertificateReloader(certFile.Name(), keyFile.Name())

	assert.NoError(t, reloader.Close())
	assert.Equal(t, io.EOF, reloader.Close())
}
