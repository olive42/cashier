package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type key interface{}
type keyfunc func(int) (key, ssh.PublicKey, error)

var (
	keytypes = map[string]keyfunc{
		"rsa":   generateRSAKey,
		"ecdsa": generateECDSAKey,
	}
)

func generateRSAKey(bits int) (key, ssh.PublicKey, error) {
	k, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	pub, err := ssh.NewPublicKey(&k.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	return k, pub, nil
}

func generateECDSAKey(bits int) (key, ssh.PublicKey, error) {
	var curve elliptic.Curve
	switch bits {
	case 256:
		curve = elliptic.P256()
	case 384:
		curve = elliptic.P384()
	case 521:
		curve = elliptic.P521()
	default:
		return nil, nil, fmt.Errorf("Unsupported key size. Valid sizes are '256', '384', '521'")
	}
	k, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	pub, err := ssh.NewPublicKey(&k.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	return k, pub, nil
}

func generateKey(keytype string, bits int) (key, ssh.PublicKey, error) {
	f, ok := keytypes[keytype]
	if !ok {
		var valid []string
		for k := range keytypes {
			valid = append(valid, k)
		}
		return nil, nil, fmt.Errorf("Unsupported key type %s. Valid choices are %s", keytype, valid)
	}
	return f(bits)
}
