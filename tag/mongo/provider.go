package mongo

import (
	"context"
	"github.com/wutipong/mangaweb/tag"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
)

var uri string
var database string

const CollectionName = "tags"

type Provider struct {
	client *mongo.Client
}

func Init(con string, db string) error {
	uri = con
	database = db

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
	p.getCollection().Indexes().CreateOne(ctx, model)

	return nil
}

func New() (p Provider, err error) {
	ctx := context.Background()

	p.client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))

	return
}

func (p *Provider) getCollection() *mongo.Collection {
	return p.client.Database(database).Collection(CollectionName)
}

func (p *Provider) IsTagExist(name string) bool {
	ctx := context.Background()
	result := p.getCollection().FindOne(ctx, bson.D{{Key: "name", Value: name}})

	var t tag.Tag

	err := result.Decode(&t)

	return err != nil
}

func (p *Provider) Write(t tag.Tag) error {
	ctx := context.Background()

	_, err := p.getCollection().UpdateOne(
		ctx, bson.D{{Key: "name", Value: t.Name}}, bson.M{"$set": t}, options.Update().SetUpsert(true))

	return err
}

func (p *Provider) Delete(t tag.Tag) error {
	ctx := context.Background()
	_, err := p.getCollection().DeleteOne(ctx, bson.D{{Key: "name", Value: t.Name}})

	return err
}

func (p *Provider) Read(name string) (t tag.Tag, err error) {
	ctx := context.Background()
	result := p.getCollection().FindOne(ctx, bson.D{{Key: "name", Value: name}})

	if result.Err() != nil {
		err = tag.ErrTagNotFound.Wrap(result.Err()).Format(name)
		return
	}

	err = result.Decode(&t)

	return
}

func (p *Provider) ReadAll() (tags []tag.Tag, err error) {
	ctx := context.Background()
	cursor, err := p.getCollection().Find(ctx, bson.D{})
	if err != nil {
		return
	}

	for cursor.Next(ctx) {
		i := tag.Tag{}
		err = cursor.Decode(&i)
		if err != nil {
			return
		}

		tags = append(tags, i)
	}

	return
}

func (p *Provider) NeedSetup() (b bool, err error) {
	ctx := context.Background()
	collectionNames, err := p.client.Database(database).ListCollectionNames(ctx, bson.D{})

	b = slices.Contains(collectionNames, CollectionName)

	return
}

func (p *Provider) Close() error {
	ctx := context.Background()
	return p.client.Disconnect(ctx)
}
