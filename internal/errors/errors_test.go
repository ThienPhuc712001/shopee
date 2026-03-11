package errors

import (
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	err := New(ErrBadRequest, "test message")
	if err.Error() != "[BAD_REQUEST] test message" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestAppError_WithDetails(t *testing.T) {
	err := New(ErrBadRequest, "test").WithDetails(map[string]string{"field": "email"})
	if err.Details == nil {
		t.Error("expected details to be set")
	}
}

func TestAppError_ToResponse(t *testing.T) {
	err := New(ErrNotFound, "user not found").WithPath("/api/users").WithMethod("GET")
	resp := err.ToResponse()

	if !resp.Error {
		t.Error("expected Error to be true")
	}
	if resp.Code != ErrNotFound {
		t.Errorf("expected code %s, got %s", ErrNotFound, resp.Code)
	}
	if resp.Path != "/api/users" {
		t.Errorf("expected path /api/users, got %s", resp.Path)
	}
}

func TestBadRequest(t *testing.T) {
	err := BadRequest("invalid input")
	if err.Code != ErrBadRequest {
		t.Errorf("expected ErrBadRequest, got %s", err.Code)
	}
	if err.HTTPStatus != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", err.HTTPStatus)
	}
}

func TestNotFound(t *testing.T) {
	err := NotFound("User")
	if err.Code != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %s", err.Code)
	}
	if err.HTTPStatus != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", err.HTTPStatus)
	}
}

func TestInternal(t *testing.T) {
	err := Internal("something went wrong")
	if err.Code != ErrInternal {
		t.Errorf("expected ErrInternal, got %s", err.Code)
	}
	if err.HTTPStatus != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", err.HTTPStatus)
	}
}

func TestIs(t *testing.T) {
	err := NotFound("User")
	if !Is(err, ErrNotFound) {
		t.Error("expected Is to return true")
	}
	if Is(err, ErrBadRequest) {
		t.Error("expected Is to return false")
	}
}

func TestWrap(t *testing.T) {
	originalErr := &AppError{Code: ErrDatabase, Message: "original"}
	err := Wrap(originalErr, ErrDatabase, "db failed")
	if err.Err != originalErr {
		t.Error("expected underlying error to be preserved")
	}
	if err.Code != ErrDatabase {
		t.Errorf("expected ErrDatabase, got %s", err.Code)
	}
}

func TestGetHTTPStatus(t *testing.T) {
	tests := []struct {
		code     ErrorCode
		expected int
	}{
		{ErrBadRequest, http.StatusBadRequest},
		{ErrUnauthorized, http.StatusUnauthorized},
		{ErrNotFound, http.StatusNotFound},
		{ErrInternal, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		err := New(tt.code, "test")
		if GetHTTPStatus(err) != tt.expected {
			t.Errorf("expected %d for %s, got %d", tt.expected, tt.code, GetHTTPStatus(err))
		}
	}
}
