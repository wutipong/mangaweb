package mongo

import (
	"context"
	"mangaweb/meta"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	p, err := New()
	if err != nil {
		return err
	}
	defer p.Close()

	model := mongo.IndexModel{
		Keys: bson.M{
			"name": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	ctx := context.Background()
	p.getItemCollection().Indexes().CreateOne(ctx, model)

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

	_, err := p.getItemCollection().UpdateOne(
		ctx, bson.D{{"name", i.Name}}, bson.M{"$set": i}, options.Update().SetUpsert(true))

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

func (p *Provider) Find(name string) (items []meta.Item, err error) {
	ctx := context.Background()
	pattern := ".*" + name + ".*"
	regex := primitive.Regex{Pattern: pattern, Options: "i"}
	cursor, err := p.getItemCollection().Find(ctx,
		bson.D{{
			"name",
			bson.D{{"$regex", regex}},
		}},
	)
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

func (p *Provider) NeedSetup() (b bool, err error) {
	ctx := context.Background()
	collectionNames, err := p.client.Database(DatabaseName).ListCollectionNames(ctx, bson.D{})

	b = len(collectionNames) == 0

	return
}

func (p *Provider) Close() error {
	ctx := context.Background()
	return p.client.Disconnect(ctx)
}
