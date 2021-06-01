package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BucketMerkle struct {
	col *mongo.Collection
}

func NewBucketMerkle(_ context.Context, db *mongo.Database) (*BucketMerkle, error) {
	s := &BucketMerkle{col: db.Collection("bucketmerkle")}
	return s, nil
}

type Merkle struct {
	BucketKey string `bson:"_id"`
	Value     string `bson:"value"`
}

func (k *BucketMerkle) Create(ctx context.Context, bucketKey, value string) (*Merkle, error) {
	bm := &Merkle{
		BucketKey: bucketKey,
		Value:     value,
	}
	if _, err := k.col.InsertOne(ctx, bm); err != nil {
		return nil, err
	}
	return bm, nil
}

func (k *BucketMerkle) Get(ctx context.Context, key string) (*Merkle, error) {
	res := k.col.FindOne(ctx, bson.M{"_id": key})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var raw bson.M
	if err := res.Decode(&raw); err != nil {
		return nil, err
	}
	return &Merkle{
		BucketKey: raw["_id"].(string),
		Value:     raw["value"].(string),
	}, nil
}
