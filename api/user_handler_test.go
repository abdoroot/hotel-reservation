package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	store db.UserStore
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		panic(err)
	}
	userStore := db.NewMongoUserStore(client, db.TESTDBName)
	return &testdb{
		store: userStore,
	}
}

func (tdb *testdb) tearDown(t *testing.T) {
	fmt.Println("--- dropping user collection")
	tdb.store.Drop(context.TODO())
}

func TestPostUser(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)
	userHandler := NewUserHandler(db.store)

	app := fiber.New()
	app.Post("/", userHandler.HandlePostUser)

	createUserReq := types.CreateUserRequest{
		FirstName:         "Abdelhadi",
		LastName:          "Abdelhadi",
		Email:             "abd@kk.cc",
		EncreptedPassword: "122669889",
	}
	js, err := json.Marshal(createUserReq)
	if err != nil {
		t.Errorf(err.Error())
	}

	req, _ := http.NewRequest("POST", "/", bytes.NewReader(js))
	req.Header.Add("Content-Type", "Application/json")

	repUser := types.User{}
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = json.NewDecoder(resp.Body).Decode(&repUser)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(repUser.EncreptedPassword) > 0 {
		t.Errorf("expected cncreptedpassword to not be show at the json : %v", createUserReq.EncreptedPassword)
	}

	if createUserReq.FirstName != repUser.FirstName {
		t.Errorf("expected firstName %v but got %v", createUserReq.FirstName, createUserReq.FirstName)
	}

	if createUserReq.LastName != repUser.LastName {
		t.Errorf("expected lastName %v but got %v", createUserReq.LastName, createUserReq.LastName)
	}

	if createUserReq.Email != repUser.Email {
		t.Errorf("expected email %v but got %v", createUserReq.Email, createUserReq.Email)
	}
}

func TestGetUser(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)
	userHandler := NewUserHandler(db.store)

	insertedUser, err := CreateForTestingPerposeUser(t)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println("userId :", insertedUser.ID.Hex())
	insertedUserId := insertedUser.ID.Hex()
	app := fiber.New()
	app.Get("/:id", userHandler.HandleGetUser)

	req, err := http.NewRequest("GET", strings.Join([]string{"/", insertedUserId}, ""), nil)
	req.Header.Add("Content-Type", "Application/json")
	fmt.Printf("req url =%v\n", req.URL)
	if err != nil {
		t.Errorf(err.Error())
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf(err.Error())
	}

	repUser := types.User{}
	err = json.NewDecoder(resp.Body).Decode(&repUser)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got := repUser.ID.Hex(); insertedUserId != got {
		t.Errorf("expected id %v got %v", insertedUserId, got)
	}

	if insertedUser.Email != repUser.Email {
		t.Errorf("expected email %v got %v", insertedUser.Email, repUser.Email)
	}

	if insertedUser.FirstName != repUser.FirstName {
		t.Errorf("expected FirstName %v got %v", insertedUser.FirstName, repUser.FirstName)
	}

	if insertedUser.LastName != repUser.LastName {
		t.Errorf("expected LastName %v got %v", insertedUser.LastName, repUser.LastName)
	}
}

func TestGetUsers(t *testing.T) {

	var (
		gotUsers           []*types.User
		insertedUserd      []*types.User
		insertedUserdCount = 5
	)
	for i := 0; i < insertedUserdCount; i++ {
		u, err := CreateForTestingPerposeUser(t)
		if err != nil {
			t.Error("error creating user list")
		}
		insertedUserd = append(insertedUserd, u)
	}

	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Add("Content-type", "Aplication/json")
	if err != nil {
		t.Error("error creating NewRequest")
	}

	tdb := setup(t)
	defer tdb.tearDown(t)

	userHandler := NewUserHandler(tdb.store)
	app := fiber.New()
	app.Get("/", userHandler.HandleGetUsers)

	resp, err := app.Test(req)
	if err != nil {
		t.Error("app error:", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&gotUsers)
	if err != nil {
		t.Errorf("decode error :%v", err.Error())
	}

	if len(gotUsers) != insertedUserdCount {
		t.Errorf("expected %v users got %v", insertedUserdCount, len(gotUsers))
	}
}

func CreateForTestingPerposeUser(t *testing.T) (*types.User, error) {
	db := setup(t)
	userHandler := NewUserHandler(db.store)
	insertedUser := &types.User{}
	//create user request
	u := &types.CreateUserRequest{
		FirstName:         "ahmed",
		LastName:          "mohamed",
		Email:             "ahmed@gdc.ae",
		EncreptedPassword: "64647488",
	}
	ujson, err := json.Marshal(u)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("marshal error"), err)
	}

	app := fiber.New()
	app.Post("/", userHandler.HandlePostUser)
	req, err := http.NewRequest("POST", "/", bytes.NewReader(ujson))
	req.Header.Add("Content-Type", "Application/json")
	if err != nil {
		return nil, errors.Join(fmt.Errorf("req error"), err)
	}

	resp, err := app.Test(req)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("app resp error"), err)
	}

	err = json.NewDecoder(resp.Body).Decode(insertedUser)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("NewDecoder error"), err)
	}
	//fmt.Printf("inserted user %+v", insertedUser)
	return insertedUser, err
}
