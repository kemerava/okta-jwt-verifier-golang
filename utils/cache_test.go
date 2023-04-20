package utils_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/kemerava/okta-jwt-verifier-golang/utils"
)

type Value struct {
	key string
}

func TestNewDefaultCache(t *testing.T) {
	lookup := func(key string, client *http.Client) (interface{}, error) {
		return &Value{key: key}, nil
	}
	cache, err := utils.NewDefaultCache(lookup, 5*time.Minute, 10*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	first, firstErr := cache.Get("first", nil)
	if firstErr != nil {
		t.Fatalf("Expected no error, got %v", firstErr)
	}
	if _, ok := first.(*Value); !ok {
		t.Error("Expected first to be a *Value")
	}

	second, secondErr := cache.Get("second", nil)
	if secondErr != nil {
		t.Fatalf("Expected no error, got %v", secondErr)
	}
	if _, ok := second.(*Value); !ok {
		t.Error("Expected second to be a *Value")
	}

	if first == second {
		t.Error("Expected first and second to be different")
	}

	firstAgain, firstAgainErr := cache.Get("first", nil)
	if firstAgainErr != nil {
		t.Fatalf("Expected no error, got %v", firstAgainErr)
	}
	if first != firstAgain {
		t.Error("Expected cached value to be the same")
	}
}
