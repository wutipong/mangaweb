package mongo

import (
	"context"
	"golang.org/x/exp/slices"

	"github.com/wutipong/mangaweb/meta"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri string
var databaseName string

const CollectionName = "items"

type Provider struct {
	client *mongo.Client
}

func Init(con string, db string) error {
	uri = con
	databaseName = db

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
	return p.client.Database(databaseName).Collection(CollectionName)
}

func (p *Provider) IsItemExist(name string) bool {
	ctx := context.Background()
	result := p.getCollection().FindOne(ctx, bson.D{{Key: "name", Value: name}})

	var item meta.Item

	err := result.Decode(&item)

	return err != nil
}

func (p *Provider) Write(i meta.Item) error {
	ctx := context.Background()

	_, err := p.getCollection().UpdateOne(
		ctx, bson.D{{Key: "name", Value: i.Name}}, bson.M{"$set": i}, options.Update().SetUpsert(true))

	return err
}

func (p *Provider) Delete(i meta.Item) error {
	ctx := context.Background()
	_, err := p.getCollection().DeleteOne(ctx, bson.D{{Key: "name", Value: i.Name}})

	return err
}

func (p *Provider) Read(name string) (i meta.Item, err error) {
	ctx := context.Background()
	result := p.getCollection().FindOne(ctx, bson.D{{Key: "name", Value: name}})

	err = result.Decode(&i)

	return
}
func (p *Provider) Open(name string) (i meta.Item, err error) {
	ctx := context.Background()
	result := p.getCollection().FindOne(ctx, bson.D{{Key: "name", Value: name}})

	err = result.Decode(&i)

	return
}

func (p *Provider) ReadAll() (items []meta.Item, err error) {
	ctx := context.Background()
	cursor, err := p.getCollection().Find(ctx, bson.D{})
	if err != nil {
		return
	}

	for cursor.Next(ctx) {
		i := meta.Item{}
		err = cursor.Decode(&i)
		if err != nil {
			return
		}

		items = append(items, i)
	}

	return
}

func (p *Provider) Search(criteria []meta.SearchCriteria, sort meta.SortField, order meta.SortOrder, pageSize int, page int) (items []meta.Item, err error) {
	ctx := context.Background()

	filter := createFilter(criteria)
	opts := options.Find().SetAllowDiskUse(true)

	orderInt := 0
	switch order {
	case meta.SortOrderAscending:
		orderInt = 1

	case meta.SortOrderDescending:
		orderInt = -1
	}
	switch sort {
	case meta.SortFieldName:
		opts.SetSort(bson.D{{"name", orderInt}})

	case meta.SortFieldCreateTime:
		opts.SetSort(bson.D{{"create_time", orderInt}})
	}

	opts.SetSkip(int64(pageSize * page)).SetLimit(int64(pageSize))

	cursor, err := p.getCollection().Find(ctx, filter, opts)
	if err != nil {
		return
	}

	for cursor.Next(ctx) {
		i := meta.Item{}
		err = cursor.Decode(&i)
		if err != nil {
			return
		}

		items = append(items, i)
	}

	return
}

func (p *Provider) Count(criteria []meta.SearchCriteria) (count int64, err error) {
	ctx := context.Background()

	opts := options.Count()

	filter := createFilter(criteria)
	count, err = p.getCollection().CountDocuments(ctx, filter, opts)

	return
}

func createNameRegex(name string) primitive.Regex {
	pattern := ".*" + name + ".*"
	regex := primitive.Regex{Pattern: pattern, Options: "i"}
	return regex
}

func createFilter(criteria []meta.SearchCriteria) bson.D {
	output := make(bson.D, 0)
	for _, c := range criteria {
		switch c.Field {
		case meta.SearchFieldName:
			{
				name := c.Value.(string)
				regex := createNameRegex(name)

				output = append(output, bson.E{"name", regex})
				break
			}

		case meta.SearchFieldFavorite:
			{
				output = append(output, bson.E{"favorite", c.Value})
				break
			}
		}
	}
	return output
}

func (p *Provider) NeedSetup() (b bool, err error) {
	ctx := context.Background()
	collectionNames, err := p.client.Database(databaseName).ListCollectionNames(ctx, bson.D{})

	b = slices.Contains(collectionNames, CollectionName)

	return
}

func (p *Provider) Close() error {
	ctx := context.Background()
	return p.client.Disconnect(ctx)
}
