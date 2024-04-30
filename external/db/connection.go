package db

import (
	"context"
	"database/sql"
	"fmt"
	"go-fitness/external/config"
	"go-fitness/external/logger/sl"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/fx"
	"log/slog"
	"net/url"
)

type MongoConnection struct {
	MongoClient *mongo.Client
	DBName      string
}

func NewMysqlDatabase(lc fx.Lifecycle, log *slog.Logger, cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local",
		cfg.DB.MysqlUser,
		cfg.DB.MysqlPassword,
		cfg.DB.MysqlHost+":"+cfg.DB.MysqlPort,
		cfg.DB.MysqlDBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info("Starting database connection")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Closing database connection")
			return db.Close()
		},
	})

	return db, nil
}

func NewMongoDatabase(lc fx.Lifecycle, log *slog.Logger, cfg *config.Config) (MongoConnection, error) {
	const op = "db.NewMongoDatabase"

	log = log.With(
		sl.String("op", op),
	)

	ctx := context.Background()

	var mongoConnection MongoConnection

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authMechanism=%s&authSource=admin",
		cfg.MongoDB.User,
		url.QueryEscape(cfg.Password),
		cfg.MongoDB.Host,
		cfg.MongoDB.Port,
		cfg.MongoDB.DBName,
		cfg.MongoDB.AuthMechanism,
	)

	log.Info("Connecting to MongoDB", sl.String("uri", uri))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return mongoConnection, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error("failed to ping MongoDB", sl.Err(err))
		return mongoConnection, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	if client == nil {
		log.Error("failed to connect to MongoDB")
		return mongoConnection, fmt.Errorf("failed to connect to MongoDB")
	}

	log.Info("MongoDB connected")

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info("Starting MongoDB connection")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Closing MongoDB connection")
			return client.Disconnect(ctx)
		},
	})

	mongoConnection = MongoConnection{
		MongoClient: client,
		DBName:      cfg.MongoDB.DBName,
	}

	return mongoConnection, nil
}
