package tls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestCertAndKey(certFile, keyFile string) error {
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
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
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
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	ioutil.WriteFile(certFile, caPEM.Bytes(), 0644)
	ioutil.WriteFile(keyFile, caPrivKeyPEM.Bytes(), 0644)

	return nil
}

func TestNewCertificateReloader_BadCertFile(t *testing.T) {
	certFile, _ := ioutil.TempFile("", "cert.*.pem")
	keyFile, _ := ioutil.TempFile("", "key.*.pem")
	defer certFile.Close()
	defer keyFile.Close()
	createTestCertAndKey(certFile.Name(), keyFile.Name())
	reloader, err := NewCertificateReloader("", keyFile.Name(), time.Minute)
	assert.Nil(t, reloader)
	assert.NotNil(t, err)
}

func TestNewCertificateReloader_BadKeyFile(t *testing.T) {
	certFile, _ := ioutil.TempFile("", "cert.*.pem")
	keyFile, _ := ioutil.TempFile("", "key.*.pem")
	defer certFile.Close()
	defer keyFile.Close()
	createTestCertAndKey(certFile.Name(), keyFile.Name())
	reloader, err := NewCertificateReloader(certFile.Name(), "", time.Minute)
	assert.Nil(t, reloader)
	assert.NotNil(t, err)
}

func TestNewCertificateReloader_BadTtl(t *testing.T) {
	certFile, _ := ioutil.TempFile("", "cert.*.pem")
	keyFile, _ := ioutil.TempFile("", "key.*.pem")
	defer certFile.Close()
	defer keyFile.Close()
	createTestCertAndKey(certFile.Name(), keyFile.Name())
	reloader, err := NewCertificateReloader(certFile.Name(), keyFile.Name(), time.Second)
	assert.Nil(t, reloader)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, syscall.EINVAL))
}

func TestNewCertificateReloader(t *testing.T) {
	certFile, _ := ioutil.TempFile("", "cert.*.pem")
	keyFile, _ := ioutil.TempFile("", "key.*.pem")
	defer certFile.Close()
	defer keyFile.Close()
	createTestCertAndKey(certFile.Name(), keyFile.Name())
	reloader, err := NewCertificateReloader(certFile.Name(), keyFile.Name(), 2*time.Minute)
	assert.NotNil(t, reloader)
	assert.Nil(t, err)
	assert.NotNil(t, reloader.cert)
	assert.Equal(t, certFile.Name(), reloader.certFile)
	assert.Equal(t, keyFile.Name(), reloader.keyFile)
	assert.NotNil(t, reloader.mu)
	assert.Equal(t, int64(0), reloader.reloadCount)
	assert.Equal(t, 2*time.Minute, reloader.ttl)
}

func TestGetCertificateFunc(t *testing.T) {
	certFile, _ := ioutil.TempFile("", "cert.*.pem")
	keyFile, _ := ioutil.TempFile("", "key.*.pem")
	defer certFile.Close()
	defer keyFile.Close()
	createTestCertAndKey(certFile.Name(), keyFile.Name())

	reloader, _ := NewCertificateReloader(certFile.Name(), keyFile.Name(), time.Minute)
	f := reloader.GetCertificateFunc()
	cert, err := f(nil)
	assert.NotNil(t, cert)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), reloader.reloadCount)
	assert.True(t, reloader.reloadDeadline.Sub(time.Now()) > 45*time.Second)

	reloader.ttl = 5 * time.Second
	reloader.reloadDeadline = time.Now()

	getCount := 10
	ch := make(chan int, getCount)
	wg := sync.WaitGroup{}
	wg.Add(getCount)
	for i := 0; i < getCount; i++ {
		go func() {
			defer wg.Done()
			<-ch
			cert, err = f(nil)
			assert.NotNil(t, cert)
			assert.Nil(t, err)
		}()
	}
	for i := 0; i < getCount; i++ {
		ch <- i
	}
	wg.Wait()
	time.Sleep(5 * time.Second)
	assert.Equal(t, int64(1), reloader.reloadCount)

	cert, err = f(nil)
	assert.NotNil(t, cert)
	assert.Nil(t, err)
	time.Sleep(5 * time.Second)
	assert.Equal(t, int64(2), reloader.reloadCount)
}
