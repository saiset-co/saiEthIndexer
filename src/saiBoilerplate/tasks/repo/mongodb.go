package repo

import (
	"context"

	"github.com/webmakom-com/saiBoilerplate/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SomeRepo struct {
	Collection *mongo.Collection
}

func New(col *mongo.Collection) *SomeRepo {
	return &SomeRepo{Collection: col}
}

func (r *SomeRepo) Set(ctx context.Context, entity *types.Some) error {
	entity.ID = primitive.NewObjectID()
	_, err := r.Collection.InsertOne(ctx, entity)
	return err
}

func (r *SomeRepo) GetAll(ctx context.Context) ([]*types.Some, error) {
	var result []*types.Some

	cur, err := r.Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var s types.Some
		err := cur.Decode(&s)
		if err != nil {
			return result, err
		}

		result = append(result, &s)
	}
	if err := cur.Err(); err != nil {
		return result, err
	}
	cur.Close(ctx)

	if len(result) == 0 {
		return result, mongo.ErrNoDocuments
	}

	return result, nil
}
