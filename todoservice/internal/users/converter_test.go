package users_test

import (
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/users"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertUserToModel(t *testing.T) {
	converter := users.NewConverter()
	entity := users.Entity{
		ID:        "1",
		Email:     "test@example.com",
		GithubID:  "gh123",
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expected := models.User{
		ID:        "1",
		Email:     "test@example.com",
		GithubID:  "gh123",
		Role:      "admin",
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	result := converter.ConvertUserToModel(entity)
	assert.Equal(t, expected, result)
}

func TestConvertUserToEntity(t *testing.T) {
	converter := users.NewConverter()
	user := models.User{
		ID:        "1",
		Email:     "test@example.com",
		GithubID:  "gh123",
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expected := users.Entity{
		ID:        "1",
		Email:     "test@example.com",
		GithubID:  "gh123",
		Role:      "admin",
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	result := converter.ConvertUserToEntity(user)
	assert.Equal(t, expected, result)
}
