package models

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	name := "张三"
	email := "zhangsan@example.com"
	
	user := NewUser(name, email)
	
	if user.Name != name {
		t.Errorf("期望用户名为 %s，实际为 %s", name, user.Name)
	}
	
	if user.Email != email {
		t.Errorf("期望邮箱为 %s，实际为 %s", email, user.Email)
	}
	
	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt 不应该为零值")
	}
	
	if user.UpdatedAt.IsZero() {
		t.Error("UpdatedAt 不应该为零值")
	}
}

func TestUserValidate(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "有效用户",
			user: &User{
				Name:  "张三",
				Email: "zhangsan@example.com",
			},
			wantErr: false,
		},
		{
			name: "用户名为空",
			user: &User{
				Name:  "",
				Email: "zhangsan@example.com",
			},
			wantErr: true,
		},
		{
			name: "邮箱为空",
			user: &User{
				Name:  "张三",
				Email: "",
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserUpdateName(t *testing.T) {
	user := NewUser("张三", "zhangsan@example.com")
	oldUpdateTime := user.UpdatedAt
	
	// 等待一小段时间确保时间戳不同
	time.Sleep(1 * time.Millisecond)
	
	newName := "李四"
	user.UpdateName(newName)
	
	if user.Name != newName {
		t.Errorf("期望用户名为 %s，实际为 %s", newName, user.Name)
	}
	
	if !user.UpdatedAt.After(oldUpdateTime) {
		t.Error("UpdatedAt 应该被更新")
	}
}

func TestUserString(t *testing.T) {
	user := &User{
		ID:    1,
		Name:  "张三",
		Email: "zhangsan@example.com",
	}
	
	expected := "User{ID: 1, Name: 张三, Email: zhangsan@example.com}"
	actual := user.String()
	
	if actual != expected {
		t.Errorf("期望字符串为 %s，实际为 %s", expected, actual)
	}
}