package adapter

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/samber/mo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"todoe/domain/task/domain"
	"todoe/domain/task/port"
)

type MongoViewRepository struct {
	clientIO    mo.IOEither[*mongo.Client]
	once        sync.Once
	cached      mo.Either[error, *mongo.Client]
	initialized atomic.Bool
}

var _ port.ViewRepository = (*MongoViewRepository)(nil)

func NewMongoViewRepository(clientIO mo.IOEither[*mongo.Client]) *MongoViewRepository {
	return &MongoViewRepository{clientIO: clientIO}
}

func (r *MongoViewRepository) getClient() mo.Either[error, *mongo.Client] {
	r.once.Do(func() {
		r.cached = r.clientIO.Run()
		r.initialized.Store(true)
	})
	return r.cached
}

func (r *MongoViewRepository) collection() (*mongo.Collection, error) {
	either := r.getClient()
	if either.IsLeft() {
		return nil, either.MustLeft()
	}
	return either.MustRight().Database("todoe").Collection("tasks_view"), nil
}

func (r *MongoViewRepository) Upsert(ctx context.Context, task domain.Task) mo.Result[struct{}] {
	col, err := r.collection()
	if err != nil {
		return mo.Err[struct{}](err)
	}

	opts := options.Replace().SetUpsert(true)
	_, err = col.ReplaceOne(ctx, bson.M{"_id": task.ID}, task, opts)
	if err != nil {
		return mo.Err[struct{}](err)
	}

	return mo.Ok(struct{}{})
}

func (r *MongoViewRepository) FindAll(ctx context.Context) mo.Result[[]domain.Task] {
	col, err := r.collection()
	if err != nil {
		return mo.Err[[]domain.Task](err)
	}

	cursor, err := col.Find(ctx, bson.D{})
	if err != nil {
		return mo.Err[[]domain.Task](err)
	}
	defer cursor.Close(ctx)

	var tasks []domain.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return mo.Err[[]domain.Task](err)
	}

	return mo.Ok(tasks)
}

func (r *MongoViewRepository) FindByID(ctx context.Context, id bson.ObjectID) mo.Result[domain.Task] {
	col, err := r.collection()
	if err != nil {
		return mo.Err[domain.Task](err)
	}

	var task domain.Task
	err = col.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return mo.Err[domain.Task](err)
	}

	return mo.Ok(task)
}
