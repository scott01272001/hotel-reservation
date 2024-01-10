package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/scott/hotel-reservation/db"
	"github.com/scott/hotel-reservation/types"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testdburi = "mongodb+srv://Cluster28936:YkhTcXV+bHBl@hotel-reservation.jtvk28l.mongodb.net/?retryWrites=true&w=majority"
	dbname    = "hotel-reservation-test"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbname),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	param := types.CreateUserparam{
		FirstName: "test",
		LastName:  "test",
		Email:     "test@gmail.com",
		Password:  "testaaa",
	}
	b, _ := json.Marshal(param)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, int(2*time.Second))
	if err != nil {
		log.Fatalf("Error testing POST request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equalf(t, 200, resp.StatusCode, "status code should be 200")
}
