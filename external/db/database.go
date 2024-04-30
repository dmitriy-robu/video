package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SqlInterface interface {
	DoInTransaction(do func(tx *sql.Tx) error) error
	Begin() error
	Commit() error
	Rollback() error
	GetTx() *sql.Tx
	GetExecer() QueryExecer
}

type QueryExecer interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type DataBase struct {
	db    *sql.DB
	tx    *sql.Tx
	mongo MongoConnection
}

func NewMysql(
	db *sql.DB,
) *DataBase {
	return &DataBase{
		db: db,
	}
}

func (u *DataBase) Begin() error {
	var err error
	u.tx, err = u.db.Begin()
	return err
}

func (u *DataBase) DoInTransaction(do func(tx *sql.Tx) error) (err error) {
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = rollbackErr
				return
			}
			panic(p)
		} else if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = rollbackErr
				return
			}
		} else {
			err = tx.Commit()
		}
	}()

	err = do(tx)
	return err
}

func (u *DataBase) Commit() error {
	if u.tx == nil {
		return nil
	}

	err := u.tx.Commit()
	if err != nil {
		return err
	}

	u.tx = nil
	return err
}

func (u *DataBase) Rollback() error {
	if u.tx == nil {
		return fmt.Errorf("transaction is not started: %w", errors.New("rollback"))
	}

	err := u.tx.Rollback()
	if err != nil {
		return err
	}

	u.tx = nil
	return nil
}

func (u *DataBase) GetDB() *sql.DB {
	return u.db
}

func (u *DataBase) GetTx() *sql.Tx {
	return u.tx
}

func (u *DataBase) GetExecer() QueryExecer {
	if u.tx != nil {
		return u.tx
	}
	return u.db
}

func (u *DataBase) CheckConnection(db *sql.DB) error {
	conn, err := db.Conn(context.Background())
	if err != nil {
		return err
	}
	defer func(conn *sql.Conn) {
		if err := conn.Close(); err != nil {
			return
		}
	}(conn)

	if err = conn.PingContext(context.Background()); err != nil {
		return err
	}

	return nil
}

type MongoDBInterface interface {
	Find(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	InsertOne(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error)
	Aggregate(ctx context.Context, collection string, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	CountDocuments(ctx context.Context, collection string, filter interface{}, opts ...*options.CountOptions) (int64, error)
}

func NewMongo(
	mongoConnection MongoConnection,
) *DataBase {
	return &DataBase{
		mongo: mongoConnection,
	}
}

func (u *DataBase) Find(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	col := u.mongo.MongoClient.Database(u.mongo.DBName).Collection(collection)
	return col.Find(ctx, filter, opts...)
}

func (u *DataBase) FindOne(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	col := u.mongo.MongoClient.Database(u.mongo.DBName).Collection(collection)
	return col.FindOne(ctx, filter, opts...)
}

func (u *DataBase) InsertOne(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	col := u.mongo.MongoClient.Database(u.mongo.DBName).Collection(collection)
	return col.InsertOne(ctx, document)
}

func (u *DataBase) UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	col := u.mongo.MongoClient.Database(u.mongo.DBName).Collection(collection)
	return col.UpdateOne(ctx, filter, update)
}

func (u *DataBase) UpdateMany(ctx context.Context, collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	col := u.mongo.MongoClient.Database(u.mongo.DBName).Collection(collection)
	return col.UpdateMany(ctx, filter, update)
}

func (u *DataBase) DeleteOne(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	col := u.mongo.MongoClient.Database(u.mongo.DBName).Collection(collection)
	return col.DeleteOne(ctx, filter)
}

func (u *DataBase) DeleteMany(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	col := u.mongo.MongoClient.Database(u.mongo.DBName).Collection(collection)
	return col.DeleteMany(ctx, filter)
}

func (u *DataBase) Aggregate(ctx context.Context, collection string, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	col := u.mongo.MongoClient.Database(u.mongo.DBName).Collection(collection)
	return col.Aggregate(ctx, pipeline, opts...)
}

func (u *DataBase) CountDocuments(ctx context.Context, collection string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	col := u.mongo.MongoClient.Database(u.mongo.DBName).Collection(collection)
	return col.CountDocuments(ctx, filter, opts...)
}
