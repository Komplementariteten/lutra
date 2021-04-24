package pages

import (
	"context"
	"testing"
	"time"

	"lutra/db"
	"lutra/model"
)

func setupDatabase() (*db.Db, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	db, err := db.NewConnection(ctx, "mongodb+srv://dev:VmV5lBWnSPgog5xn@devcluster.pqajf.mongodb.net/defaultDatabase?retryWrites=true&w=majority", "TestDatabase")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestAuth_AddUser(t *testing.T) {
	db, err := setupDatabase()
	if err != nil {
		t.Fatal(err)
	}

	user := &model.User{
		Name:     "Test User",
		Login:    "test@test.org",
		Disabled: false,
		Email:    "test@test.org",
		Terms:    make(map[string]bool),
	}

	err = AddNewUser(context.Background(), db, user)
	if err != nil {
		t.Fatal(err)
	}

	err = AddNewUser(context.Background(), db, user)
	if err != nil {
		t.Fatal(err)
	}

	p, err := db.Collection(context.Background(), "auth")
	if err != nil {
		t.Fatal(err)
	}
	p.DeleteDatabase(context.Background())
}
