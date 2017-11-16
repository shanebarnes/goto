package tls

import (
    "bytes"
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "errors"
)

func EncodePrivateKey(priv interface{}) (*bytes.Buffer, error) {
    var err error = nil
    var block *pem.Block = nil
    var keyOut *bytes.Buffer = new(bytes.Buffer)

    if block, err = getPemBlock(priv); err == nil {
        err = pem.Encode(keyOut, block)
    }

    return keyOut, err
}

func GeneratePrivateKey(algorithm string) (interface{}, error) {
    var err error = nil
    var key interface{}

    switch algorithm {
    case "1024":
        key, err = rsa.GenerateKey(rand.Reader, 1024)
    case "2048":
        key, err = rsa.GenerateKey(rand.Reader, 2048)
    case "3072":
        key, err = rsa.GenerateKey(rand.Reader, 3072)
    case "P224":
        key, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
    case "P256":
        key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    case "P384":
        key, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
    case "P521":
        key, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
    default:
        err = errors.New("Unsupported algorithm: " + algorithm)
    }

    return key, err
}

func getPemBlock(priv interface{}) (*pem.Block, error) {
    var err error = nil
    var block *pem.Block = nil
    var buffer []byte

    switch k := priv.(type) {
    case *rsa.PrivateKey:
        block =  &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
    case *ecdsa.PrivateKey:
        buffer, err = x509.MarshalECPrivateKey(k)
        if err == nil {
            block = &pem.Block{Type: "EC PRIVATE KEY", Bytes: buffer}
        }
    default:
    }

    return block, err
}

func getPublicKey(priv interface{}) interface{} {
    var pub interface{}

    switch k := priv.(type) {
    case *rsa.PrivateKey:
        pub = &k.PublicKey
    case *ecdsa.PrivateKey:
        pub = &k.PublicKey
    default:
    }

    return pub
}
