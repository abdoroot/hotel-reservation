package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abdoroot/hotel-reservation/db/fixtures"
	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestAuthenticatWithWrongData(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)
	insertedUser := fixtures.AddUser(tdb.store, "test", "user", false)
	_ = insertedUser
	authHandler := NewAuthHandler(tdb.store)

	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Post("/", authHandler.HandleAuthUser)

	//create request
	param := types.AuthUserRequest{
		Email:    "test@user.com",
		Password: "anyWrongPassword",
	}
	b, err := json.Marshal(param)
	if err != nil {
		t.Fatalf("fail to marshal json %v", err)
	}
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "Application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("fail to test %v", err)
	}

	errResp := &Error{}
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		t.Fatalf("fail to decode the response %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code to be 400 got %v", resp.StatusCode)
	}

	if errResp.Msg != "error email or password" {
		t.Fatalf("expected response msg to be error email or password got %v", errResp.Msg)
	}

}

func TestAuthenticatWithValidData(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)
	insertedUser := fixtures.AddUser(tdb.store, "test", "user", false)
	authHandler := NewAuthHandler(tdb.store)

	app := fiber.New()
	app.Post("/", authHandler.HandleAuthUser)

	//create request
	param := types.AuthUserRequest{
		Email:    "test@user.com",
		Password: "test_user",
	}
	b, err := json.Marshal(param)
	if err != nil {
		t.Fatalf("fail to marshal json %v", err)
	}
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "Application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("fail to test %v", err)
	}

	authResponse := &types.AuthResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatalf("fail to decode the response %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status code to be 200 got %v", resp.StatusCode)
	}

	if authResponse.User.Email != insertedUser.Email {
		t.Fatalf("expected email to be %v got %v", insertedUser.Email, authResponse.User.Email)
	}

	if len(authResponse.Token) == 0 {
		t.Fatal("expected token to be greater then zero")
	}

}
