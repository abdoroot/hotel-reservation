package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

const (
	testUserEmail = "abdorootin@gmail.com"
	testUserPwd   = "@Aaeeerrrrp77"
)

func InsertedUser(t *testing.T, Store *db.Store) *types.User {
	r := &types.CreateUserRequest{
		FirstName:         "abdlhadi",
		LastName:          "Mohamed",
		Email:             testUserEmail,
		EncreptedPassword: testUserPwd,
	}
	u, err := r.CreateUserFromUserRequest()
	if err != nil {
		t.Fatalf("fail to create user request %v", err)
	}

	insertedUser, err := Store.User.InsertUser(context.TODO(), u)
	if err != nil {
		t.Fatalf("fail to create user request %v", err)
	}

	return insertedUser
}

func TestAuthenticatWithWrongData(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)
	insertedUser := InsertedUser(t, tdb.store)
	_ = insertedUser
	authHandler := NewAuthHandler(tdb.store)

	app := fiber.New()
	app.Post("/", authHandler.HandleAuthUser)

	//create request
	param := types.AuthUserRequest{
		Email:    testUserEmail,
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

	errResp := &types.ErrorResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		t.Fatalf("fail to decode the response %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code to be 400 got %v", resp.StatusCode)
	}

	if errResp.Type != "error" {
		t.Fatalf("expected response type to be error got %v", errResp.Type)
	}

}

func TestAuthenticatWithValidData(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)
	insertedUser := InsertedUser(t, tdb.store)
	authHandler := NewAuthHandler(tdb.store)

	app := fiber.New()
	app.Post("/", authHandler.HandleAuthUser)

	//create request
	param := types.AuthUserRequest{
		Email:    testUserEmail,
		Password: testUserPwd,
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
