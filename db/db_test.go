package db

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"lutra/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestPost struct {
	Name string `json:"title" ,omitempty`
	Text string `json:"text" ,omitempty`
}

type BigTestData struct {
	Name    string
	Date    time.Time
	Desc    string
	Bool    bool
	Content string
	ID      int
}

func TestPool_Update(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	db, err := NewConnection(ctx, "mongodb+srv://dev:VmV5lBWnSPgog5xn@devcluster.pqajf.mongodb.net/defaultDatabase?retryWrites=true&w=majority", "TestDatabase")
	if err != nil {
		t.Fatal(err)
	}

	p, err := db.Collection(ctx, "TestCollection")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 15; i++ {
		d := &BigTestData{
			Name:    util.GetRandomString(5),
			Date:    time.Now().Local().Add(time.Duration(i) * time.Hour),
			Desc:    util.GetRandomString(20),
			Bool:    i%2 == 0,
			Content: util.GetRandomString(10),
			ID:      i,
		}
		bytes, err := bson.Marshal(d)
		if err != nil {
			t.Fatal(err)
		}
		var doc bson.D
		err = bson.Unmarshal(bytes, &doc)
		if err != nil {
			t.Fatal(err)
		}
		p.Add(ctx, &doc)
	}

	updated, err := p.Update(ctx, &bson.M{"date": bson.M{"$gt": time.Now().Local()}}, map[string]interface{}{"Bool": false, "Content": "22"})
	if err != nil {
		t.Fatal(err)
	}
	if updated < 14 {
		t.Fatalf("To few Items updated, updated: %d", updated)
	}

	p.DeleteDatabase(ctx)

}

func TestPool_Add(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	db, err := NewConnection(ctx, "mongodb+srv://dev:VmV5lBWnSPgog5xn@devcluster.pqajf.mongodb.net/defaultDatabase?retryWrites=true&w=majority", "TestDatabase")
	if err != nil {
		t.Fatal(err)
	}

	p, err := db.Collection(ctx, "TestCollection")
	if err != nil {
		t.Fatal(err)
	}

	if p.Name != "TestCollection" {
		t.Fatal("Pool should store Collection Name")
	}

	_, err = p.Add(ctx, &bson.D{{"test1", "test2"}})

	if err != nil {
		t.Fatal(err)
	}

	err = db.Close(ctx)
	if err != nil {
		t.Fatal(err)
	}

}

func TestDb_Collection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := NewConnection(ctx, "mongodb+srv://dev:VmV5lBWnSPgog5xn@devcluster.pqajf.mongodb.net/defaultDatabase?retryWrites=true&w=majority", "TestDatabase")
	if err != nil {
		t.Fatal(err)
	}

	p, err := db.Collection(ctx, "TestCollection")
	if err != nil {
		t.Fatal(err)
	}

	if p.Name != "TestCollection" {
		t.Fatal("Pool should store Collection Name")
	}

	err = db.Close(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := NewConnection(ctx, "wrongConnectionString", "nothing")
	if err == nil {
		t.Fatal("Wrong Connection string did not return Error")
	}

	db, err := NewConnection(ctx, "mongodb+srv://dev:VmV5lBWnSPgog5xn@devcluster.pqajf.mongodb.net/defaultDatabase?retryWrites=true&w=majority", "TestDatabase")
	if err != nil {
		t.Fatal(err)
	}

	err = db.Close(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if db.IsConnected {
		t.Fatal("Closing connection should disconnect to server")
	}
}

func TestDb_Files(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	db, err := NewConnection(ctx, "mongodb+srv://dev:VmV5lBWnSPgog5xn@devcluster.pqajf.mongodb.net/defaultDatabase?retryWrites=true&w=majority", "TestDatabase")
	if err != nil {
		t.Fatal(err)
	}

	p, err := db.Collection(ctx, "TestCollection")
	if err != nil {
		t.Fatal(err)
	}

	if p.Name != "TestCollection" {
		t.Fatal("Pool should store Collection Name")
	}

	refName := "fobarfile.dat"
	t.Log(os.Getwd())
	dat, err := ioutil.ReadFile("pools.go")
	if err != nil {
		t.Error(err)
	}
	_, err = db.AddFile(ctx, refName, dat)
	if err != nil {
		t.Fatal(err)
	}
	rDat, err := db.ReadFile(ctx, refName)
	if err != nil {
		t.Fatal(err)
	}
	if len(rDat) != len(dat) {
		t.Fatal("sizes do not match")
	}

	if !bytes.Equal(rDat, dat) {
		t.Fatal("byte Arrays don't match")
	}

	err = p.DeleteDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPool_Find(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	db, err := NewConnection(ctx, "mongodb+srv://dev:VmV5lBWnSPgog5xn@devcluster.pqajf.mongodb.net/defaultDatabase?retryWrites=true&w=majority", "TestDatabase")
	if err != nil {
		t.Fatal(err)
	}

	p, err := db.Collection(ctx, "TestCollection")
	if err != nil {
		t.Fatal(err)
	}

	if p.Name != "TestCollection" {
		t.Fatal("Pool should store Collection Name")
	}

	_, err = p.Add(ctx, &bson.D{{"n", "test1"}, {"s", 1}, {"m", 1}})
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.Add(ctx, &bson.D{{"n", "test2"}, {"s", 2}, {"m", 1}})
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.Add(ctx, &bson.D{{"n", "test3"}, {"s", 3}, {"m", 2}})
	if err != nil {
		t.Fatal(err)
	}

	m, err := p.Find(ctx, &bson.M{"m": 1})
	if err != nil {
		t.Fatal(err)
	}

	if len(m) != 2 {
		t.Fatal("Wrong number of items returned")
	}

	for i := 0; i < len(m); i++ {
		if m[i]["m"].(int32) != 1 {
			t.Fatal("Filter filed does not match")
		}
	}

	err = p.DeleteDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}

}

func TestPool_Get(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	db, err := NewConnection(ctx, "mongodb+srv://dev:VmV5lBWnSPgog5xn@devcluster.pqajf.mongodb.net/defaultDatabase?retryWrites=true&w=majority", "TestDatabase")
	if err != nil {
		t.Fatal(err)
	}

	p, err := db.Collection(ctx, "TestCollection")
	if err != nil {
		t.Fatal(err)
	}

	if p.Name != "TestCollection" {
		t.Fatal("Pool should store Collection Name")
	}

	_, err = p.Add(ctx, &bson.D{{"n", "test1"}, {"s", 1}})
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.Add(ctx, &bson.D{{"n", "test2"}, {"s", 2}})
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.Add(ctx, &bson.D{{"n", "test3"}, {"s", 3}})
	if err != nil {
		t.Fatal(err)
	}

	m, err := p.Get(ctx, &bson.M{"n": "test2"})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Result %s:%d", m["n"], m["s"])
	if m["s"].(int32) != 2 {
		t.Fatal("Document not found in Database")
	}

	err = p.DeleteDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

// TestMongoDbConnection Tests MongoDb Connection
func TestMongoDbConnection(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://dev:VmV5lBWnSPgog5xn@devcluster.pqajf.mongodb.net/defaultDatabase?retryWrites=true&w=majority"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = client.Disconnect(ctx)
		if err != nil {
			t.Fatal(err)
		}
		cancel()
	}()
}

func TestMongoDbWriteToConnection(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://dev:VmV5lBWnSPgog5xn@devcluster.pqajf.mongodb.net/defaultDatabase?retryWrites=true&w=majority"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = client.Disconnect(ctx)
		if err != nil {
			t.Fatal(err)
		}
		cancel()
	}()

	testPost := TestPost{"test1", "irgendein text"}
	collection := client.Database("GoTestDb").Collection("testposts")

	insertResult, err := collection.InsertOne(ctx, testPost)
	if err != nil {
		t.Error(err)
	}
	var deletedDocument bson.M
	err = collection.FindOneAndDelete(ctx, bson.D{{"_id", insertResult.InsertedID}}).Decode(&deletedDocument)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("deleted Document %v", deletedDocument)
	t.Logf("deleted Document %v", deletedDocument)
}
