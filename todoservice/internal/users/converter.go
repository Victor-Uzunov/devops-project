package users

import (
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
)

type Converter struct{}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertUserToModel(entity Entity) models.User {
	return models.User{
		ID:        entity.ID,
		Email:     entity.Email,
		GithubID:  entity.GithubID,
		Role:      entity.Role,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (c *Converter) ConvertUserToEntity(user models.User) Entity {
	return Entity{
		ID:        user.ID,
		Email:     user.Email,
		GithubID:  user.GithubID,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
