package db

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Db gives access to the Mongodb Database
type Db struct {
	client       *mongo.Client
	database     *mongo.Database
	IsConnected  bool
	DatabaseName string
	options      *options.ClientOptions
	collections  map[string]bool
}

// TrackedItem is the Item Returned from Database as a Tracked Item
type TrackedItem struct {
	InsertedID primitive.ObjectID
	Collection string
	IsFile     bool
}

// Pool is a Subset of Data in the MongoDb
type Pool struct {
	db        *Db
	c         *mongo.Collection
	Name      string
	isCreated bool
}

const (
	// DatabaseNotConnected Error
	DatabaseNotConnected = "Database not Connected"
)

var errNotConnected = errors.New("Handle not connected to Database")
var dbIDField = "_id"

// NewConnection in Context to the Database
func NewConnection(ctx context.Context, connectionString string, dataBase string) (db *Db, err error) {
	opt := options.Client()
	opt.ApplyURI(connectionString)

	client, err := mongo.NewClient(opt)
	if err != nil {
		log.Printf("Connection Error %v", err)
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	dbHandle := client.Database(dataBase)

	db = &Db{client: client, database: dbHandle, DatabaseName: dataBase, IsConnected: true, options: opt}
	db.collections = make(map[string]bool)

	db.readAllCollections(ctx)

	return db, nil
}

// Close the Database Connection
func (db *Db) Close(ctx context.Context) error {
	if !db.IsConnected {
		return errNotConnected
	}
	db.client.Disconnect(ctx)
	db.IsConnected = false
	db.database = nil
	db.client = nil
	db.DatabaseName = ""
	db.collections = make(map[string]bool)
	return nil
}

// DeleteDatabase Delete a given Database by Pool
func (p *Pool) DeleteDatabase(ctx context.Context) error {
	if !p.db.IsConnected {
		return errNotConnected
	}
	err := p.db.database.Drop(ctx)
	p.db.DatabaseName = ""
	p.db.database = nil
	p.db.IsConnected = false
	return err
}

func (db *Db) readAllCollections(ctx context.Context) error {
	collections, err := db.database.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return err
	}
	for _, col := range collections {
		db.collections[col] = true
	}
	return nil
}

// Collection creates or utilizes a Database subset
func (db *Db) Collection(ctx context.Context, collectionName string) (p *Pool, err error) {
	if !db.IsConnected {
		return nil, errNotConnected
	}

	if _, ok := db.collections[collectionName]; ok {
		p = &Pool{db: db, c: db.database.Collection(collectionName), Name: collectionName, isCreated: true}
		return p, nil
	}
	err = db.database.CreateCollection(ctx, collectionName)
	if err != nil {
		return nil, err
	}
	db.collections[collectionName] = true
	p = &Pool{db: db, c: db.database.Collection(collectionName), Name: collectionName, isCreated: true}
	return p, nil
}

// Get returns a Single Result from the search Map
func (p *Pool) Get(ctx context.Context, searchMap *bson.M) (bson.M, error) {
	if !p.db.IsConnected {
		return nil, errNotConnected
	}
	if !p.isCreated {
		return nil, fmt.Errorf("Collection %s does not exist", p.Name)
	}
	var decoded bson.M
	if err := p.c.FindOne(ctx, searchMap).Decode(&decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

// GetId finds the object id in a bson.M
func (p *Pool) GetId(m bson.M) primitive.ObjectID {
	if id, ok := m["_id"].(primitive.ObjectID); ok {
		return id
	}
	return primitive.NilObjectID
}

// Find searches for all matches in Db
func (p *Pool) Find(ctx context.Context, filter *bson.M) ([]bson.M, error) {
	if !p.db.IsConnected {
		return nil, errNotConnected
	}
	if !p.isCreated {
		return nil, fmt.Errorf("Collection %s does not exist", p.Name)
	}
	cursor, err := p.c.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var datas []bson.M
	if err = cursor.All(ctx, &datas); err != nil {
		return nil, err
	}
	return datas, nil
}

// Delete removes a tracked Item form the Database
func (p *Pool) Delete(ctx context.Context, id primitive.ObjectID) error {
	if !p.db.IsConnected {
		return errNotConnected
	}
	if !p.isCreated {
		return fmt.Errorf("Collection %s does not exist", p.Name)
	}
	deleteResult, err := p.c.DeleteOne(ctx, &bson.M{dbIDField: id})
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount != 1 {
		return fmt.Errorf("More than one Item deleted")
	}
	return nil
}

// DeleteTracked removes a tracked Item form the Database
func (p *Pool) DeleteTracked(ctx context.Context, item *TrackedItem) error {
	if !p.db.IsConnected {
		return errNotConnected
	}
	if !p.isCreated {
		return fmt.Errorf("Collection %s does not exist", p.Name)
	}
	deleteResult, err := p.c.DeleteOne(ctx, &bson.M{dbIDField: item.InsertedID})
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount != 1 {
		return fmt.Errorf("More than one Item deleted")
	}
	return nil
}

// AddM adds a Entity where the Entity is represented by bson.M to the MongoDb
func (p *Pool) AddM(ctx context.Context, document *bson.M) (item *TrackedItem, err error) {
	if !p.db.IsConnected {
		return nil, errNotConnected
	}

	if !p.isCreated {
		return nil, fmt.Errorf("Collection %s does not exist", p.Name)
	}

	insertResult, err := p.c.InsertOne(ctx, document)
	if err != nil {
		log.Printf("Failed to Add %v, err: %v", document, err)
		return nil, err
	}
	item = &TrackedItem{IsFile: false}
	item.InsertedID = insertResult.InsertedID.(primitive.ObjectID)
	item.Collection = p.Name
	return item, nil
}

// Add adds a Entity to the MongoDb
func (p *Pool) Add(ctx context.Context, document *bson.D) (item *TrackedItem, err error) {
	if !p.db.IsConnected {
		return nil, errNotConnected
	}

	if !p.isCreated {
		return nil, fmt.Errorf("Collection %s does not exist", p.Name)
	}

	insertResult, err := p.c.InsertOne(ctx, document)
	if err != nil {
		log.Printf("Failed to Add %v, err: %v", document, err)
		return nil, err
	}
	item = &TrackedItem{IsFile: false}
	item.InsertedID = insertResult.InsertedID.(primitive.ObjectID)
	item.Collection = p.Name
	return item, nil
}

// Update updates all Itemes by the filter from the given map
func (p *Pool) Update(ctx context.Context, filter *bson.M, fields map[string]interface{}) (int64, error) {
	if !p.db.IsConnected {
		return 0, errNotConnected
	}

	if !p.isCreated {
		return 0, fmt.Errorf("Collection %s does not exist", p.Name)
	}

	var updateDocument bson.D
	for key, element := range fields {
		updateDocument = append(updateDocument, bson.E{key, element})
	}

	document := bson.D{
		{"$set", updateDocument},
	}

	result, err := p.c.UpdateMany(ctx, filter, document)
	if err != nil {
		return 0, err
	}
	if result.ModifiedCount == 0 {
		return 0, fmt.Errorf("No Changes where made")
	}
	return result.ModifiedCount, nil
}
func (db *Db) ReadFile(ctx context.Context, reference string) ([]byte, error) {
	if !db.IsConnected {
		return nil, errNotConnected
	}
	bucket, err := gridfs.NewBucket(db.database)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	count, err := bucket.DownloadToStreamByName(reference, &buffer)
	if err != nil {
		return nil, err
	}
	if count != int64(buffer.Len()) {
		return nil, fmt.Errorf("Not all data loaded from buffer")
	}
	return buffer.Bytes(), nil
}

// AddFile adds a File to Mongodbs Gridfs
func (db *Db) AddFile(ctx context.Context, reference string, p []byte) (*TrackedItem, error) {
	if !db.IsConnected {
		return nil, errNotConnected
	}

	bucket, err := gridfs.NewBucket(db.database)
	if err != nil {
		return nil, err
	}

	uploadStream, err := bucket.OpenUploadStream(reference)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = uploadStream.Close(); err != nil {
			panic(err)
		}
	}()

	if err = uploadStream.SetWriteDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return nil, err
	}

	if _, err = uploadStream.Write(p); err != nil {
		return nil, err
	}

	tracked := &TrackedItem{IsFile: true}
	tracked.InsertedID = uploadStream.FileID.(primitive.ObjectID)
	tracked.Collection = bucket.GetFilesCollection().Name()
	return tracked, nil
}
