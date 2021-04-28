package mongo

import (
	"context"
	"mangaweb/meta"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri string

const CollectionName = "items"
const DatabaseName = "manga"

type Provider struct {
	client *mongo.Client
}

func Init(con string) error {
	uri = con

	return nil
}

func New() (p Provider, err error) {
	ctx := context.Background()

	p.client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))

	return
}

func (p *Provider) getItemCollection() *mongo.Collection {
	return p.client.Database(DatabaseName).Collection(CollectionName)
}

func (p *Provider) IsItemExist(name string) bool {
	ctx := context.Background()
	result := p.getItemCollection().FindOne(ctx, bson.D{{"name", name}})

	var item meta.Item

	err := result.Decode(&item)

	return err != nil
}

func (p *Provider) Write(i meta.Item) error {
	ctx := context.Background()

	var err error
	result, err := p.getItemCollection().UpdateOne(ctx, bson.D{{"name", bson.E{"$eq", i.Name}}}, bson.M{"$set": i})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		_, err = p.getItemCollection().InsertOne(ctx, i)
	}

	return err
}

func (p *Provider) New(name string) (i meta.Item, err error) {
	i = meta.Item{
		Name:       name,
		CreateTime: time.Now(),
		Favorite:   false,
		Mutex:      new(sync.Mutex),
	}

	i.GenerateImageIndices()
	i.GenerateThumbnail()

	err = p.Write(i)

	return
}
func (p *Provider) Delete(i meta.Item) error {
	ctx := context.Background()
	_, err := p.getItemCollection().DeleteOne(ctx, bson.D{{"name", i.Name}})

	return err
}

func (p *Provider) Read(name string) (i meta.Item, err error) {
	ctx := context.Background()
	result := p.getItemCollection().FindOne(ctx, bson.D{{"name", name}})

	err = result.Decode(&i)

	return
}
func (p *Provider) Open(name string) (i meta.Item, err error) {
	ctx := context.Background()
	result := p.getItemCollection().FindOne(ctx, bson.D{{"name", name}})

	err = result.Decode(&i)

	return
}

func (p *Provider) ReadAll() (items []meta.Item, err error) {
	ctx := context.Background()
	cursor, err := p.getItemCollection().Find(ctx, bson.D{})
	if err != nil {
		return
	}

	for cursor.Next(ctx) {
		i := meta.Item{
			Mutex: new(sync.Mutex),
		}
		err = cursor.Decode(&i)
		if err != nil {
			return
		}

		items = append(items, i)
	}

	return
}

func (p *Provider) GetCount() (count int64, err error) {
	ctx := context.Background()
	count, err = p.getItemCollection().CountDocuments(ctx, bson.D{})

	return
}

func (p *Provider) Close() error {
	ctx := context.Background()
	return p.client.Disconnect(ctx)
}
