package library

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func EnvArray(envName string) []string {
	val := os.Getenv(envName)
	ar := strings.Split(val, ",")
	return ar
}
func EnvInt(envName string) int {
	val, _ := strconv.Atoi(os.Getenv(envName))
	return val
}
func EnvBool(envName string) bool {
	val := os.Getenv(envName)
	return strings.ToUpper(val) == "TRUE"
}
func EnvString(envName string) string {
	return os.Getenv(envName)
}

var ctxB = context.Background()

func RandomNumber(min, max int) int {
	randomNumberPin := rand.Intn(max-min) + min
	return randomNumberPin
}

var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyz")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func PadLeft(val, length int) string {
	return fmt.Sprintf("%0*d", length, val)
}

func ToJSONString(val interface{}) string {
	temp, _ := json.Marshal(val)
	return string(temp)
}

func ToBytes(val interface{}) []byte {
	temp, _ := json.Marshal(val)
	return temp
}

func NewHTTPClient(rt http.RoundTripper, timeOut time.Duration) HttpClient {
	obj := HttpClient{
		Client: &http.Client{
			Transport: rt,
			Timeout:   timeOut,
		},
	}
	return obj
}

func StructToUrlValue(data interface{}) (url.Values, error) {
	return query.Values(data)
}

func (hc HttpClient) GET(header http.Header, url string) ([]byte, error) {
	handleErr := func(err error) ([]byte, error) {
		return nil, err
	}

	load, err := request(hc.Client, url, "GET", header, nil)
	if err != nil {
		return handleErr(err)
	}

	return load, nil
}

func Random() string {
	return strconv.Itoa(100 + rand.Intn(1000-100))
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
func UUID() string {
	id := uuid.New()
	return id.String()
}
