package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"sync"
	"time"
)

type Reloader struct {
	certExpiry time.Time
	certFile   string
	certTls    *tls.Certificate
	closeCtx   context.Context
	closeFunc  context.CancelFunc
	closeOnce  sync.Once
	keyFile    string
	mu         *sync.RWMutex
}

func NewCertificateReloader(certFile, keyFile string) (*Reloader, error) {
	r := &Reloader{
		certFile: certFile,
		keyFile:  keyFile,
		mu:       &sync.RWMutex{},
	}

	err := r.loadCertificate()
	if err == nil {
		r.closeCtx, r.closeFunc = context.WithCancel(context.Background())
		go r.reloadCertificate()
	} else {
		r = nil
	}

	return r, err
}

func (r *Reloader) Close() error {
	err := io.EOF

	r.closeOnce.Do(func() {
		err = nil
		r.closeFunc()
	})

	return err
}

func (r *Reloader) GetCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		return r.getCertificate()
	}
}

func (r *Reloader) getCertificate() (*tls.Certificate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.certTls, nil
}

func (r *Reloader) getCertificateExpiry() time.Duration {
	// Reduce the expiry time by one hour to prevent waiting until the
	// "last minute" to attempt reloading of the certificate
	return r.certExpiry.Sub(time.Now().Add(time.Hour))
}

func (r *Reloader) loadCertificate() error {
	if certTls, err := tls.LoadX509KeyPair(r.certFile, r.keyFile); err != nil {
		return err
	} else if certX509, err := x509.ParseCertificate(certTls.Certificate[0]); err != nil {
		return err
	} else {
		r.certExpiry = certX509.NotAfter

		r.mu.Lock()
		r.certTls = &certTls
		r.mu.Unlock()

	}
	return nil
}

func (r *Reloader) reloadCertificate() {
	minDuration := time.Second * 10
	expiry := r.getCertificateExpiry()

	for {
		select {
		case <-time.After(expiry):
			if r.loadCertificate() == nil {
				if expiry = r.getCertificateExpiry(); expiry < minDuration {
					expiry = minDuration
				}
			} else {
				expiry = minDuration
			}
		case <-r.closeCtx.Done():
			return
		}
	}
}
