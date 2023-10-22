package sparing

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type ApiPayload struct {
	UID int64 `json:"uid"`
	// Unix timestamp in seconds
	DateTime int64    `json:"datetime"`
	PH       *float32 `json:"pH"`
	COD      *float32 `json:"cod"`
	BOD      *float32 `json:"bod"`
	TSS      *float32 `json:"tss"`
	NH3N     *float32 `json:"nh3n"`
	Debit    *float32 `json:"debit"`
}

type SparingApi interface {
	GetSecret() (string, error)
	Submit(secret string, payload ApiPayload) error
}

type sparing struct {
	secretUrl  string
	url        string
	testingUrl string
}

func NewSparingAPI(secretUrl, url, testingUrl string) SparingApi {
	return &sparing{
		secretUrl:  secretUrl,
		url:        url,
		testingUrl: testingUrl,
	}
}

func (s *sparing) GetSecret() (string, error) {
	return s.getSecret()
}

func (s *sparing) Submit(secret string, payload ApiPayload) error {
	return s.submit(payload, secret)
}

// Get secret from sparing server with http GET method
func (s *sparing) getSecret() (string, error) {
	if s.secretUrl == "" {
		return "", fmt.Errorf("secret url is empty")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", s.secretUrl, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code: %d. Body: %s", resp.StatusCode, string(body))
		return "", err
	}

	secret := string(body)

	return secret, nil
}

func (s *sparing) buildClaims(payload ApiPayload) (jwt.MapClaims, error) {
	mapPayload := jwt.MapClaims{
		"uid":      payload.UID,
		"datetime": payload.DateTime,
	}

	if payload.PH != nil {
		mapPayload["pH"] = payload.PH
	}

	if payload.COD != nil {
		mapPayload["cod"] = payload.COD
	}

	if payload.BOD != nil {
		mapPayload["bod"] = payload.BOD
	}

	if payload.TSS != nil {
		mapPayload["tss"] = payload.TSS
	}

	if payload.NH3N != nil {
		mapPayload["nh3n"] = payload.NH3N
	}

	if payload.Debit != nil {
		mapPayload["debit"] = payload.Debit
	}

	return mapPayload, nil
}

// Convert payload to JWT token
func (s *sparing) EncodePayload(payload ApiPayload, secret string) (string, error) {
	claims, err := s.buildClaims(payload)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Submit payload to sparing server with http POST method
func (s *sparing) submit(payload ApiPayload, secret string) error {
	token, err := s.EncodePayload(payload, secret)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", s.url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	reqBody := fmt.Sprintf(`{"token": "%s"}`, token)

	req.Body = io.NopCloser(
		io.Reader(
			strings.NewReader(reqBody),
		),
	)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code: %d", resp.StatusCode)
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != "Data Sent Successfully!" {
		return fmt.Errorf("response body: %s", string(body))
	}

	return nil
}
