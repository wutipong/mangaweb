package mongo

import (
	"context"
	"mangaweb/meta"
	"sync"

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

func (p *Provider) IsItemExist(name string) bool {
	ctx := context.Background()
	result := p.client.Database(DatabaseName).Collection(CollectionName).FindOne(ctx, bson.D{{"name", name}})

	var item meta.Item

	err := result.Decode(&item)

	return err != nil
}

func (p *Provider) Write(i meta.Item) error {
	ctx := context.Background()

	var err error
	if p.IsItemExist(i.Name) {
		_, err = p.client.Database(DatabaseName).Collection(CollectionName).UpdateOne(ctx, bson.D{{"name", i.Name}}, i)
	} else {
		_, err = p.client.Database(DatabaseName).Collection(CollectionName).InsertOne(ctx, i)
	}

	return err
}

func (p *Provider) New(name string) (i meta.Item, err error) {
	i = meta.Item{
		Mutex: new(sync.Mutex),
	}

	return
}
func (p *Provider) Delete(i meta.Item) error {
	ctx := context.Background()
	_, err := p.client.Database(DatabaseName).Collection(CollectionName).DeleteOne(ctx, bson.D{{"name", i.Name}})

	return err
}

func (p *Provider) Read(name string) (i meta.Item, err error) {
	ctx := context.Background()
	result := p.client.Database(DatabaseName).Collection(CollectionName).FindOne(ctx, bson.D{{"name", name}})

	err = result.Decode(&i)

	return
}
func (p *Provider) Open(name string) (i meta.Item, err error) {
	ctx := context.Background()
	result := p.client.Database(DatabaseName).Collection(CollectionName).FindOne(ctx, bson.D{{"name", name}})

	err = result.Decode(&i)

	return
}

func (p *Provider) ReadAll() (items []meta.Item, err error) {
	ctx := context.Background()
	cursor, err := p.client.Database(DatabaseName).Collection(CollectionName).Find(ctx, bson.D{})
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

func (p *Provider) Close() error {
	ctx := context.Background()
	return p.client.Disconnect(ctx)
}
