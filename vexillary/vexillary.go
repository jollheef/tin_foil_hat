/**
 * @file vexillary.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief work with flags
 *
 * Contain functions for work with flags, such as generate rsa key for validate
 * flag, generate flag and validate flag.
 */

package vexillary

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"
)

// GenerateKey generate rsa key
func GenerateKey() (priv *rsa.PrivateKey, err error) {
	// 128 is minimal valid length for RSA key, which can validate flag
	return rsa.GenerateKey(rand.Reader, 128)
}

// GenerateFlag generate signed flag
func GenerateFlag(priv *rsa.PrivateKey) (string, error) {

	randBuf := make([]byte, 4)

	_, err := rand.Read(randBuf)

	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, 0, randBuf)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x%x=", randBuf, signature), nil
}

// ValidFlag verify flag
func ValidFlag(flag string, pub rsa.PublicKey) (bool, error) {

	if len(flag) != 41 {
		return false, errors.New("flag length is not 41")
	}

	if flag[40] != '=' {
		return false, errors.New("no '=' at end")
	}

	randBuf, err := hex.DecodeString(flag[0:8])

	if err != nil {
		return false, err
	}

	signature, err := hex.DecodeString(flag[8:40])

	if err != nil {
		return false, err
	}

	err = rsa.VerifyPKCS1v15(&pub, 0, randBuf, signature)

	if err != nil {
		return false, err
	}

	return true, nil
}
