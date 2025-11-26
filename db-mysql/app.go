package main

import (
	"context"
	"fmt"
	"log"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// 初始化資料庫
	_ = GetDBInstance()
	log.Println("Database initialized via migration")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// User 結構體定義
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// CreateUser 創建新用戶
func (a *App) CreateUser(name, email string, age int) (string, error) {
	db := GetDBInstance()
	err := db.InsertUser(name, email, age)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	return "User created successfully", nil
}

// GetAllUsers 獲取所有用戶
func (a *App) GetAllUsers() ([]map[string]interface{}, error) {
	db := GetDBInstance()
	users, err := db.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

// GetUser 根據 ID 獲取用戶
func (a *App) GetUser(id int) (map[string]interface{}, error) {
	db := GetDBInstance()
	user, err := db.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateUser 更新用戶
func (a *App) UpdateUser(id int, name, email string, age int) (string, error) {
	db := GetDBInstance()
	err := db.UpdateUser(id, name, email, age)
	if err != nil {
		return "", fmt.Errorf("failed to update user: %w", err)
	}
	return "User updated successfully", nil
}

// DeleteUser 刪除用戶
func (a *App) DeleteUser(id int) (string, error) {
	db := GetDBInstance()
	err := db.DeleteUser(id)
	if err != nil {
		return "", fmt.Errorf("failed to delete user: %w", err)
	}
	return "User deleted successfully", nil
}

// SearchUsers 搜尋用戶（支援分頁）
func (a *App) SearchUsers(keyword string, page, pageSize int) (map[string]interface{}, error) {
	db := GetDBInstance()
	users, total, err := db.SearchUsers(keyword, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	
	result := map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
		"pageSize": pageSize,
	}
	
	return result, nil
}
