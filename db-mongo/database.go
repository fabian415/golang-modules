package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database 結構體封裝所有資料庫操作
type Database struct {
	Host       string
	Port       string
	User       string
	Password   string
	DBName     string
	AuthSource string
	Client     *mongo.Client
	DB         *mongo.Database
}

// 使用 sync.Once 來確保 Singleton 實例只被創建一次
var once sync.Once
var dbInstance *Database

// GetDBInstance 獲取 Singleton 實例
func GetDBInstance() *Database {
	once.Do(func() {
		// 讀取環境變數並記錄
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "27017")
		user := getEnv("DB_USER", "")
		password := getEnv("DB_PASSWORD", "")
		dbName := getEnv("DB_NAME", "mydb")
		authSource := getEnv("DB_AUTH_SOURCE", "admin")

		// 記錄讀取到的配置（密碼只顯示長度）
		passwordMask := "***"
		if password != "" {
			passwordMask = fmt.Sprintf("*** (%d chars)", len(password))
		}
		WriteAppLog(fmt.Sprintf("Database config - Host: %s, Port: %s, User: %s, Password: %s, DB: %s, AuthSource: %s",
			host, port, user, passwordMask, dbName, authSource), true)

		dbInstance = &Database{
			Host:       host,
			Port:       port,
			User:       user,
			Password:   password,
			DBName:     dbName,
			AuthSource: authSource,
		}
		dbInstance.Initialize()
	})
	return dbInstance
}

// getEnv 獲取環境變數，如果不存在則返回預設值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Initialize 初始化資料庫
func (d *Database) Initialize() {
	if err := d.Connect(); err != nil {
		WriteAppLog(fmt.Sprintf("Failed to connect to database: %v", err), true)
		return
	}

	if err := d.runMigrations(); err != nil {
		WriteAppLog(fmt.Sprintf("Failed to run migrations: %v", err), true)
	}
}

// Connect 連接到 MongoDB
func (d *Database) Connect() error {
	WriteAppLog("Connecting to MongoDB...", true)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 構建 MongoDB 連接字串
	var uri string
	if d.User != "" && d.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s",
			d.User, d.Password, d.Host, d.Port, d.DBName, d.AuthSource)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s/%s",
			d.Host, d.Port, d.DBName)
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 測試連接
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	d.Client = client
	d.DB = client.Database(d.DBName)

	WriteAppLog("Successfully connected to MongoDB", true)
	return nil
}

// Disconnect 關閉資料庫連接
func (d *Database) Disconnect() error {
	if d.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return d.Client.Disconnect(ctx)
	}
	return nil
}

// region <-- Migration 相關函式-->
// Migration 版本記錄結構
type MigrationVersion struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Version int                `bson:"version"`
	Applied time.Time          `bson:"applied_at"`
}

// 執行資料庫 migration
func (d *Database) runMigrations() error {
	WriteAppLog("Initializing database migration...", true)
	WriteAppLog(fmt.Sprintf("DB Host: %s, Port: %s, DB: %s", d.Host, d.Port, d.DBName), true)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 獲取當前版本
	currentVersion, err := d.getCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	WriteAppLog(fmt.Sprintf("Current database version: %d", currentVersion), true)

	// 定義所有遷移步驟
	migrations := []func(context.Context) error{
		d.migration001_CreateUsersCollection, // Version 1
		// 未來可以在這裡添加更多遷移
		// d.migration002_AddNewIndex,        // Version 2
	}

	// 執行待處理的遷移
	for i, migration := range migrations {
		version := i + 1
		if version > currentVersion {
			WriteAppLog(fmt.Sprintf("Running migration %d...", version), true)
			if err := migration(ctx); err != nil {
				return fmt.Errorf("failed to run migration %d: %w", version, err)
			}

			// 記錄遷移版本
			if err := d.setVersion(ctx, version); err != nil {
				return fmt.Errorf("failed to record migration version %d: %w", version, err)
			}

			WriteAppLog(fmt.Sprintf("Migration %d completed successfully", version), true)
		} else {
			WriteAppLog(fmt.Sprintf("Migration %d already applied, skipping", version), true)
		}
	}

	WriteAppLog("All migrations completed", true)
	return nil
}

// getCurrentVersion 獲取當前資料庫版本
func (d *Database) getCurrentVersion(ctx context.Context) (int, error) {
	collection := d.DB.Collection("migrations")
	var version MigrationVersion

	err := collection.FindOne(ctx, bson.M{}).Decode(&version)
	if err == mongo.ErrNoDocuments {
		return 0, nil // 沒有記錄表示版本為 0
	}
	if err != nil {
		return 0, err
	}

	return version.Version, nil
}

// setVersion 記錄遷移版本
func (d *Database) setVersion(ctx context.Context, version int) error {
	collection := d.DB.Collection("migrations")

	// 使用 upsert 確保只有一條記錄
	_, err := collection.UpdateOne(
		ctx,
		bson.M{},
		bson.M{
			"$set": bson.M{
				"version":   version,
				"applied_at": time.Now(),
			},
		},
		options.Update().SetUpsert(true),
	)

	return err
}

// migration001_CreateUsersCollection 創建 users 集合並設置索引
func (d *Database) migration001_CreateUsersCollection(ctx context.Context) error {
	collection := d.DB.Collection("users")

	// 創建索引
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	WriteAppLog("Created users collection with indexes", true)
	return nil
}

// endregion

// InsertUser 插入用戶資料
func (d *Database) InsertUser(name, email string, age int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := d.DB.Collection("users")

	user := bson.M{
		"name":       name,
		"email":      email,
		"age":        age,
		"created_at": time.Now(),
	}

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("email already exists")
		}
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

// GetAllUsers 獲取所有用戶
func (d *Database) GetAllUsers() ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := d.DB.Collection("users")

	// 按創建時間降序排序
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []map[string]interface{}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	// 轉換 ObjectID 為字串
	for i := range users {
		if id, ok := users[i]["_id"].(primitive.ObjectID); ok {
			users[i]["id"] = id.Hex()
			delete(users[i], "_id")
		}
		// 轉換時間為字串
		if createdAt, ok := users[i]["created_at"].(primitive.DateTime); ok {
			users[i]["created_at"] = createdAt.Time().Format(time.RFC3339)
		} else if createdAt, ok := users[i]["created_at"].(time.Time); ok {
			users[i]["created_at"] = createdAt.Format(time.RFC3339)
		}
	}

	return users, nil
}

// GetUserByID 根據 ID 獲取用戶
func (d *Database) GetUserByID(id string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	collection := d.DB.Collection("users")
	var user map[string]interface{}

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// 轉換 ObjectID 為字串
	if id, ok := user["_id"].(primitive.ObjectID); ok {
		user["id"] = id.Hex()
		delete(user, "_id")
	}
	// 轉換時間為字串
	if createdAt, ok := user["created_at"].(primitive.DateTime); ok {
		user["created_at"] = createdAt.Time().Format(time.RFC3339)
	} else if createdAt, ok := user["created_at"].(time.Time); ok {
		user["created_at"] = createdAt.Format(time.RFC3339)
	}

	return user, nil
}

// UpdateUser 更新用戶資料
func (d *Database) UpdateUser(id string, name, email string, age int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	collection := d.DB.Collection("users")

	update := bson.M{
		"$set": bson.M{
			"name":  name,
			"email": email,
			"age":   age,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("email already exists")
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser 刪除用戶
func (d *Database) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	collection := d.DB.Collection("users")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// SearchUsers 搜尋用戶（支援分頁）
func (d *Database) SearchUsers(keyword string, page, pageSize int) ([]map[string]interface{}, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := d.DB.Collection("users")

	// 構建搜尋條件（支援姓名和電子郵件）
	searchFilter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": keyword, "$options": "i"}},
			{"email": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	// 計算總數
	total, err := collection.CountDocuments(ctx, searchFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// 查詢分頁資料
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := collection.Find(ctx, searchFilter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []map[string]interface{}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, fmt.Errorf("failed to decode users: %w", err)
	}

	// 轉換 ObjectID 為字串
	for i := range users {
		if id, ok := users[i]["_id"].(primitive.ObjectID); ok {
			users[i]["id"] = id.Hex()
			delete(users[i], "_id")
		}
		// 轉換時間為字串
		if createdAt, ok := users[i]["created_at"].(primitive.DateTime); ok {
			users[i]["created_at"] = createdAt.Time().Format(time.RFC3339)
		} else if createdAt, ok := users[i]["created_at"].(time.Time); ok {
			users[i]["created_at"] = createdAt.Format(time.RFC3339)
		}
	}

	return users, int(total), nil
}
