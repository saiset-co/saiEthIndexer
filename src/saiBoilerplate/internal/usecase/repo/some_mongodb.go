package repo

import (
	"context"

	"github.com/webmakom-com/saiBoilerplate/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SomeRepo struct {
	col *mongo.Collection
}

func New(col *mongo.Collection) *SomeRepo {
	return &SomeRepo{}
}

func (r *SomeRepo) Set(ctx context.Context, entity *entity.Some) error {
	_, err := r.col.InsertOne(ctx, entity)
	return err
}

func (r *SomeRepo) GetAll(ctx context.Context) ([]*entity.Some, error) {
	var result []*entity.Some

	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var s entity.Some
		err := cur.Decode(&s)
		if err != nil {
			return result, err
		}

		result = append(result, &s)

		if err := cur.Err(); err != nil {
			return result, err
		}
		cur.Close(ctx)

		if len(result) == 0 {
			return result, mongo.ErrNoDocuments
		}

	}
	return result, nil
}
