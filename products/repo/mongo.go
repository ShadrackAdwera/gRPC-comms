package repo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		ProductEntry:  ProductEntry{},
		CategoryEntry: CategoryEntry{},
	}
}

type Models struct {
	ProductEntry  ProductEntry
	CategoryEntry CategoryEntry
}

type ProductEntry struct {
	ID         string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string    `bson:"name" json:"name"`
	Price      string    `bson:"price" json:"price"`
	Image      string    `bson:"image" json:"image"`
	CategoryId string    `bson:"categoryId" json:"categoryId"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

type CategoryEntry struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	CreatedBy   string    `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

func (p *ProductEntry) Insert(entry ProductEntry) error {
	collection := client.Database("products").Collection("products")

	_, err := collection.InsertOne(context.TODO(), ProductEntry{
		Name:       entry.Name,
		Price:      entry.Price,
		Image:      entry.Image,
		CategoryId: entry.CategoryId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		log.Println("Error creating product:", err)
		return err
	}

	return nil
}

func (c *CategoryEntry) AddCategory(cat CategoryEntry) error {
	collection := client.Database("products").Collection("categories")

	_, err := collection.InsertOne(context.TODO(), CategoryEntry{
		Name:        c.Name,
		Description: c.Description,
		CreatedBy:   c.CreatedBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		log.Println("Error creating category:", err)
		return err
	}

	return nil
}

func (p *ProductEntry) All() ([]*ProductEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("products").Collection("products")

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*ProductEntry

	for cursor.Next(ctx) {
		var item ProductEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding product into slice:", err)
			return nil, err
		} else {
			products = append(products, &item)
		}
	}

	return products, nil
}

func (p *ProductEntry) GetOne(id string) (*ProductEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("products").Collection("products")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry ProductEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (p *ProductEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("products").Collection("products")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (p *ProductEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("products").Collection("products")

	docID, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: p.Name},
				{Key: "price", Value: p.Price},
				{Key: "image", Value: p.Image},
				{Key: "categoryId", Value: p.CategoryId},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
