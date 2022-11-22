package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/entity"
	"github.com/sirupsen/logrus"
	"io"
	mRand "math/rand"
	"modernc.org/strutil"
	"os"
)

const Method = "chacha20-ietf-poly1305"

type Crypto struct {
}

func (c *Crypto) Encrypt(text []byte, key []byte) ([]byte, error) {
	cpr, err := aes.NewCipher(key)
	if err != nil {
		logrus.Errorf("error created when encrypt cpr: %s", err.Error())
		return nil, err
	}

	gcm, err := cipher.NewGCM(cpr)
	if err != nil {
		logrus.Errorf("error created gcm: %s", err.Error())
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logrus.Errorf("error read nonce: %s", err.Error())
		return nil, err
	}

	return gcm.Seal(nonce, nonce, text, nil), nil
}

func (c *Crypto) Decrypt(cipherText []byte, key []byte) ([]byte, error) {
	cpr, err := aes.NewCipher(key)
	if err != nil {
		logrus.Errorf("error created cpr when decrypt: %s", err.Error())
		return nil, err
	}

	gcm, err := cipher.NewGCM(cpr)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := cipherText[:nonceSize], cipherText[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (c *Crypto) GeneratePassword(passwordLen int) string {
	passwordRunes := make([]rune, passwordLen)

	for i := range passwordRunes {
		passwordRunes[i] = randomRune()
	}

	return string(passwordRunes)
}

func (c *Crypto) GenerateConfig(conn *entity.Connection) (string, error) {
	ds, err := hex.DecodeString(conn.EncryptedSecret)
	if err != nil {
		logrus.Errorf("error decode to bytes when gen conf: %s", err.Error())
		return "", err
	}
	plainSecret, err := c.Decrypt(ds, []byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		logrus.Errorf("error decrypted when gen conf: %s", err.Error())
		return "", err
	}
	userInfo := fmt.Sprintf("%s:%s", Method, string(plainSecret))
	encodedUserInfo := strutil.Base64Encode([]byte(userInfo))
	conf := fmt.Sprintf("ss://%s@%s:%d#vpn", string(encodedUserInfo), conn.IpAddress, conn.Port)
	return conf, nil
}

func randomRune() rune {
	i := mRand.Intn(26)

	return rune('A' + i)
}
