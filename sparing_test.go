package sparing

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestGetSecret(t *testing.T) {
	svr := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/secret" {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Write([]byte("secret"))
		}),
	)

	api := NewSparingAPI(svr.URL, "", "")
	_, err := api.GetSecret()
	if err != nil {
		assert.EqualError(t, err, "status code: 404")
		return
	}

	api = NewSparingAPI(svr.URL+"/secret", "", "")
	secret, err := api.GetSecret()
	if err != nil {
		assert.EqualError(t, err, "status code: 404")
	}

	assert.Equal(t, "secret", secret)
}

func TestEncodePayload(t *testing.T) {
	ph := float32(18692)
	cod := float32(5508)
	tss := float32(5466)
	nh3n := float32(16539)
	debit := float32(17006)

	payload := ApiPayload{
		UID:      1120800300014,
		DateTime: 1568630149,
		PH:       &ph,
		COD:      &cod,
		TSS:      &tss,
		NH3N:     &nh3n,
		Debit:    &debit,
	}

	secret := "enaknyaapasecret"
	// expected := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1aWQiOjExMjA4MDAzMDAwMTQsImRhdGV0aW1lIjoxNTY4NjMwMTQ5LCJwSCI6MTg2OTIsImNvZCI6NTUwOCwidHNzIjo1NDY2LCJuaDNuIjoxNjUzOSwiZGViaXQiOjE3MDA2fQ.vThMTJA2wMElKh2_uWoG45XdHj5jH3MOfUde88bSYzY"

	api := &sparing{}
	encoded, err := api.EncodePayload(payload, secret)
	if err != nil {
		assert.EqualError(t, err, "json: unsupported type: func()")
		return
	}

	t.Log(encoded)

	// assert.Equal(t, expected, encoded)
}

func TestSubmit(t *testing.T) {
	secret := "enaknyaapasecret"

	ph := float32(18692)
	cod := float32(5508)
	tss := float32(5466)
	nh3n := float32(16539)
	debit := float32(17006)

	svr := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/submit" {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			defer r.Body.Close()

			body, err := io.ReadAll(r.Body)
			if err != nil {
				assert.NoErrorf(t, err, "failed to read body: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			token := struct {
				Token string `json:"token"`
			}{}
			err = json.Unmarshal(body, &token)
			if err != nil {
				assert.NoErrorf(t, err, "failed to unmarshal body: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			payload := jwt.MapClaims{}
			_, err = jwt.ParseWithClaims(token.Token, &payload, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil {
				assert.NoErrorf(t, err, "failed to parse token: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			buffer, err := json.Marshal(payload)
			if err != nil {
				assert.NoErrorf(t, err, "failed to marshal payload: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			sparingPayload := ApiPayload{}
			err = json.Unmarshal(buffer, &sparingPayload)
			if err != nil {
				assert.NoErrorf(t, err, "failed to unmarshal payload: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			assert.Equal(t, int64(1120800300014), sparingPayload.UID)
			assert.Equal(t, int64(1568630149), sparingPayload.DateTime)
			assert.Equal(t, ph, *sparingPayload.PH)
			assert.Equal(t, cod, *sparingPayload.COD)
			assert.Equal(t, tss, *sparingPayload.TSS)
			assert.Equal(t, nh3n, *sparingPayload.NH3N)
			assert.Equal(t, debit, *sparingPayload.Debit)

			w.Write([]byte("Data Sent Successfully!"))
		}),
	)

	api := NewSparingAPI("", svr.URL, "")
	err := api.Submit(secret, ApiPayload{})
	if err != nil {
		assert.EqualError(t, err, "status code: 404")
	}

	api = NewSparingAPI("", svr.URL+"/submit", "")

	err = api.Submit(secret, ApiPayload{
		UID:      1120800300014,
		DateTime: 1568630149,
		PH:       &ph,
		COD:      &cod,
		TSS:      &tss,
		NH3N:     &nh3n,
		Debit:    &debit,
	})

	if err != nil {
		assert.EqualError(t, err, "status code: 404")
		return
	}
}
