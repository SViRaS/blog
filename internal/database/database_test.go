package database

import (
	"blog/internal/models"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println(".env не найден, используем системные переменные")
	}

	if err := ConnectDB(); err != nil {
		log.Fatalf("Не удалось подключиться к БД для тестов: %v", err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestUserSaveAndLoad(t *testing.T) {
	user := models.User{
		Username: "SokolFA",
		Email:    "su4ka@gmail.com",
		Password: "sokolova123",
	}

	if err := DB.Create(&user).Error; err != nil {
		t.Fatalf("Ошибка создания: %v", err)
	}

	var loaded models.User
	if err := DB.First(&loaded, user.ID).Error; err != nil {
		t.Fatalf("Ошибка чтения: %v", err)
	}

	if loaded.Username != "SokolFA" {
		t.Errorf("Ожидали 'SokolFA', получили '%s'", loaded.Username)
	}
	if loaded.Password[:4] != "$2a$" {
		t.Error("Пароль не захеширован")
	}

	DB.Delete(&loaded)
}
