package logredact

import (
	"encoding/json"
	"net/url"
	"strings"
)

const RedactedPlaceholder = "[REDACTED]"

var sensitiveFields = map[string]bool{
	"token":           true,
	"access_token":    true,
	"refresh_token":   true,
	"id_token":        true,
	"idtoken":         true,
	"api_key":         true,
	"apikey":          true,
	"secret":          true,
	"password":        true,
	"pwd":             true,
	"login":           true,
	"credential":      true,
	"signature":       true,
	"upload_url":      true,
	"download_url":    true,
	"checksum_sha256": true,
}

func IsSensitive(key string) bool {
	return sensitiveFields[strings.ToLower(key)]
}

func RedactJSON(raw []byte) (redacted string, ok bool) {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", false
	}
	redactValue(v)
	b, err := json.Marshal(v)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func redactValue(v any) {
	switch t := v.(type) {
	case map[string]any:
		for key, val := range t {
			if IsSensitive(key) {
				t[key] = RedactedPlaceholder
				continue
			}
			redactValue(val)
		}
	case []any:
		for _, item := range t {
			redactValue(item)
		}
	}
}

func SanitizeQuery(rawQuery string) string {
	if rawQuery == "" {
		return ""
	}
	params, err := url.ParseQuery(rawQuery)
	if err != nil {
		return rawQuery
	}
	redacted := false
	for key := range params {
		if IsSensitive(key) {
			params.Set(key, RedactedPlaceholder)
			redacted = true
		}
	}
	if !redacted {
		return rawQuery
	}
	return params.Encode()
}
