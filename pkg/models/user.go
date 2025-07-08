package models

import (
	"fmt"
	"time"
)

// User 用户模型
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser 创建新用户
func NewUser(name, email string) *User {
	now := time.Now()
	return &User{
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// String 实现 Stringer 接口
func (u *User) String() string {
	return fmt.Sprintf("User{ID: %d, Name: %s, Email: %s}", u.ID, u.Name, u.Email)
}

// Validate 验证用户数据
func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if u.Email == "" {
		return fmt.Errorf("邮箱不能为空")
	}
	return nil
}

// UpdateName 更新用户名
func (u *User) UpdateName(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}