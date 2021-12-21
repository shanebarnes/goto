package tls

import (
	"crypto/tls"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

const minTtl = time.Minute

type Reloader struct {
	cert           *tls.Certificate
	certFile       string
	keyFile        string
	mu             *sync.RWMutex
	reloadCount    int64
	reloadDeadline time.Time
	ttl            time.Duration
}

func NewCertificateReloader(certFile, keyFile string, ttl time.Duration) (*Reloader, error) {
	if ttl < minTtl {
		return nil, errors.WithMessage(syscall.EINVAL, "Certificate TTL "+ttl.String()+" < "+minTtl.String())
	} else if cert, err := tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		return nil, err
	} else {
		return &Reloader{
			cert:           &cert,
			certFile:       certFile,
			keyFile:        keyFile,
			mu:             &sync.RWMutex{},
			reloadCount:    int64(0),
			reloadDeadline: time.Now().Add(ttl),
			ttl:            ttl,
		}, nil
	}
}

func (r *Reloader) GetCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		return r.getCertificate()
	}
}

func (r *Reloader) getCertificate() (*tls.Certificate, error) {
	r.mu.RLock()
	reloadDeadline := r.reloadDeadline

	defer func() {
		now := time.Now()
		if now.After(reloadDeadline) {
			reload := false

			r.mu.Lock()
			// Only allow one reload per deadline expiry
			if reloadDeadline == r.reloadDeadline {
				r.reloadDeadline = now.Add(r.ttl)
				reload = true
			}
			r.mu.Unlock()

			if reload {
				go r.reloadCertificate()
			}
		}
	}()
	defer r.mu.RUnlock()

	return r.cert, nil
}

func (r *Reloader) reloadCertificate() {
	if cert, err := tls.LoadX509KeyPair(r.certFile, r.keyFile); err == nil {
		r.mu.Lock()
		r.reloadCount++
		r.cert = &cert
		r.mu.Unlock()
	}
}
