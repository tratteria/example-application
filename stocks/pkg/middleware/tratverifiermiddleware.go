package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

type VerifyTraTRequest struct {
	Path            string          `json:"path"`
	Method          string          `json:"method"`
	QueryParameters json.RawMessage `json:"queryParameters"`
	Headers         json.RawMessage `json:"headers"`
	Body            json.RawMessage `json:"body"`
}

type VerifyTraTResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason"`
}

func readAndReplaceBody(r *http.Request) (json.RawMessage, error) {
	if r.Body == nil {
		return []byte("{}"), nil
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body.Close()

	r.Body = io.NopCloser(bytes.NewBuffer(data))

	if len(data) == 0 {
		return []byte("{}"), nil
	}

	return data, nil
}

func convertMapToJson(data map[string]string) (json.RawMessage, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		bytes = []byte("{}")
	}

	return json.RawMessage(bytes), nil
}

// TODO: handle keys with multiple values.
func convertHeaderToJson(headers http.Header) (json.RawMessage, error) {
	headerMap := make(map[string]string)
	for key, values := range headers {
		headerMap[key] = values[0]
	}

	return convertMapToJson(headerMap)
}

func getVerifyTraTRequest(r *http.Request) (*VerifyTraTRequest, error) {
	body, err := readAndReplaceBody(r)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %w", err)
	}

	headersJson, err := convertHeaderToJson(r.Header)
	if err != nil {
		return nil, fmt.Errorf("error reading request header: %w", err)
	}

	//TODO: handle keys with multiple values.
	queryParams := make(map[string]string)
	for key, values := range r.URL.Query() {
		queryParams[key] = values[0]
	}

	queryParamsJson, err := convertMapToJson(queryParams)
	if err != nil {
		return nil, fmt.Errorf("error reading query parameters: %w", err)
	}

	details := &VerifyTraTRequest{
		Path:            r.URL.Path,
		Method:          r.Method,
		QueryParameters: queryParamsJson,
		Headers:         headersJson,
		Body:            body,
	}

	return details, nil
}

func getTraTVerifierMiddleware(traTVerifierEndpoint *url.URL, httpClient *http.Client, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			verifyTraTRequest, err := getVerifyTraTRequest(r)
			if err != nil {
				logger.Error("Error creating verify trat request.", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			requestBody, err := json.Marshal(verifyTraTRequest)
			if err != nil {
				logger.Error("Error marshalling verify trat request to JSON.", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			req, err := http.NewRequest(http.MethodPost, traTVerifierEndpoint.String(), bytes.NewBuffer(requestBody))
			if err != nil {
				logger.Error("Error creating trat verification request.", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := httpClient.Do(req)
			if err != nil {
				logger.Error("Error sending trat verification request.", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				logger.Error("Received non-OK response on trat verification.", zap.Int("status_code", resp.StatusCode))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("Error reading trat verification response body.", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			var verifyResponse VerifyTraTResponse

			err = json.Unmarshal(body, &verifyResponse)
			if err != nil {
				logger.Error("Error unmarshalling trat verification response body.", zap.Error(err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				return
			}

			if !verifyResponse.Valid {
				logger.Info("Invalid trat.", zap.String("reason", verifyResponse.Reason))
				http.Error(w, "Invalid trat", http.StatusForbidden)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
