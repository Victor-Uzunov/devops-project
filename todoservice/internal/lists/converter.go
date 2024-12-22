package lists

import (
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
)

type Converter struct{}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertListToModel(entity Entity) models.List {
	return models.List{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		OwnerID:     entity.OwnerID,
		SharedWith:  entity.SharedWith,
		Tags:        pkg.JSONRawMessageFromNullableString(entity.Tags),
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Visibility:  entity.Visibility,
	}
}

func (c *Converter) ConvertListToEntity(list models.List) Entity {
	return Entity{
		ID:          list.ID,
		Name:        list.Name,
		Description: list.Description,
		OwnerID:     list.OwnerID,
		SharedWith:  list.SharedWith,
		Tags:        pkg.NewNullableStringFromJSONRawMessage(list.Tags),
		CreatedAt:   list.CreatedAt,
		UpdatedAt:   list.UpdatedAt,
		Visibility:  list.Visibility,
	}
}

func (c *Converter) ConvertAccessToModel(entity AccessEntity) models.Access {
	return models.Access{
		ListID: entity.ListID,
		UserID: entity.UserID,
		Role:   entity.Role,
		Status: entity.Status,
	}
}

func (c *Converter) ConvertAccessToEntity(access models.Access) AccessEntity {
	return AccessEntity{
		ListID: access.ListID,
		UserID: access.UserID,
		Role:   access.Role,
		Status: access.Status,
	}
}
