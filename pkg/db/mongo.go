package db

import (
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

var (
	MongoClient *MongoDBClient
)

type MongoDBClient struct {
	client     *mongo.Client
	selectedDb *mongo.Database
	prefix     string
}

func NewMongoDBClient(dsnConf DsnConfig) (*MongoDBClient, error) {
	// 创建一个15秒超时的上下文
	connectCtx, connectCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer connectCancel() // 确保在函数结束时调用cancel，避免资源泄露

	// 构建连接到MongoDB的Uri
	dsn := fmt.Sprintf("mongodb://%s:%s",
		dsnConf.Host,
		strconv.Itoa(dsnConf.Port),
	)

	// 设置连接选项
	option := options.Client().ApplyURI(dsn)
	if dsnConf.User != "" {
		option = option.SetAuth(options.Credential{
			Username:   dsnConf.User,
			Password:   dsnConf.Pass,
			AuthSource: dsnConf.Dbname,
		})
	}

	client, connErr := mongo.Connect(connectCtx, option)
	if connErr != nil {
		return nil, connErr
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 释放上下文

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	// 选择数据库
	mongoDb := client.Database(dsnConf.Dbname)

	return &MongoDBClient{
		client:     client,
		selectedDb: mongoDb,
		prefix:     dsnConf.Prefix,
	}, nil
}

func (r *MongoDBClient) Close() error {
	return r.client.Disconnect(context.Background())
}

func (r *MongoDBClient) GetCollection(collection string) *mongo.Collection {
	return r.selectedDb.Collection(r.prefix + collection)
}

// InsertOne Mongo插入一条数据
func (r *MongoDBClient) InsertOne(collection string, data bson.M) error {
	c := r.GetCollection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bsonData := convertor.MapToSlice(data, func(key string, value interface{}) bson.E {
		return bson.E{Key: key, Value: value}
	})

	_, err := c.InsertOne(ctx, bsonData)

	return err
}
