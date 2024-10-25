package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mergermarket/go-pkcs7"
)

func SecureClose[T any](c chan T) {
	defer func() {
		_ = recover()
	}()
	close(c)
}

func SecureCloseIn[T any](c chan<- T) {
	defer func() {
		_ = recover()
	}()
	close(c)
}

func GetHost(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	cs := strings.Split(IPAddress, ":")
	return cs[0]
}

func GetPort(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	cs := strings.Split(r.RemoteAddr, ":")

	if len(cs) == 2 {
		return cs[1]
	} else {
		return ""
	}
}

// EncryptAES256CBC Encrypt encrypts plain text string into cipher text string
func EncryptAES256CBC(aesKey, unencrypted string) (string, error) {
	key := []byte(aesKey)
	plainText := []byte(unencrypted)
	plainText, err := pkcs7.Pad(plainText, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf(`plainText: "%s" has error`, plainText)
	}
	if len(plainText)%aes.BlockSize != 0 {
		err := fmt.Errorf(`plainText: "%s" has the wrong block size`, plainText)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return fmt.Sprintf("%x", cipherText), nil
}

// DecryptAES256CBC Decrypt decrypts cipher text string into plain text string
func DecryptAES256CBC(aesKey string, encrypted string) (string, error) {
	key := []byte(aesKey)
	cipherText, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("cipherText too short")
		return "", err
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		err = errors.New("cipherText is not a multiple of the block size")
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, _ = pkcs7.Unpad(cipherText, aes.BlockSize)
	return fmt.Sprintf("%s", cipherText), nil
}

func SplitText(text string, chunkSize int) []string {
	var chunks []string

	for len(text) > chunkSize {
		chunks = append(chunks, text[:chunkSize])
		text = text[chunkSize:]
	}

	chunks = append(chunks, text)

	return chunks
}

func RFC1123Z(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(time.RFC1123Z)
}

func RFC3339(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func ParseStringToRFC3339(t string) (*time.Time, error) {
	if t == "" {
		return nil, errors.New("no time specified")
	}
	tm, err := time.Parse(time.RFC3339, t)
	if err != nil {
		tm, err = time.Parse("2006-01-02T15:04:05.000Z07:00", t) // RFC3339Mili
		if err != nil {
			tm, err = time.Parse(time.RFC3339Nano, t)
			if err != nil {
				return nil, errors.New("use RFC3339 format string for datetime")
			}
		}
	}
	return &tm, nil
}

func MustParseStringToRFC3339(t string) *time.Time {
	tm, _ := ParseStringToRFC3339(t)
	return tm
}

func CurrentTimeAsRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func CurrentTimeAsRFC822() string {
	return time.Now().UTC().Format(time.RFC822)
}

func SliceUintTo64(xs []uint) []uint64 {
	ys := make([]uint64, len(xs))
	for i, x := range xs {
		ys[i] = uint64(x)
	}
	return ys
}

func Slice64ToUint(xs []uint64) []uint {
	ys := make([]uint, len(xs))
	for i, x := range xs {
		ys[i] = uint(x)
	}
	return ys
}

func MarshalToString(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func GenerateSignature(secret string) string {
	if strings.Trim(secret, " ") == "" {
		return ""
	}
	mac := hmac.New(sha256.New, []byte(secret))
	unixTs := time.Now().Unix()
	ts := unixTs - (unixTs % 30)
	mac.Write([]byte(fmt.Sprintf("%d", ts)))
	sign := mac.Sum(nil)
	return hex.EncodeToString(sign)
}
