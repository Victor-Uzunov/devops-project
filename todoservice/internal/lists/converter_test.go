package lists_test

import (
	"encoding/json"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/lists"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
	constants2 "github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestConvertListToModel(t *testing.T) {
	converter := lists.NewConverter()
	creationTime := time.Now()
	updateTime := time.Now()

	tags, err := json.Marshal([]constants2.TagType{constants2.TagWork, constants2.TagFinance})
	require.NoError(t, err)

	entity := lists.Entity{
		ID:          "1",
		Name:        "Test List",
		Description: "This is a test list",
		OwnerID:     "owner1",
		SharedWith:  []string{"user1", "user2"},
		Tags:        pkg.NewValidNullableString(string(tags)),
		CreatedAt:   creationTime,
		UpdatedAt:   updateTime,
		Visibility:  constants2.VisibilityPrivate,
	}

	expectedModel := models.List{
		ID:          "1",
		Name:        "Test List",
		Description: "This is a test list",
		OwnerID:     "owner1",
		SharedWith:  []string{"user1", "user2"},
		Tags:        json.RawMessage(string(tags)),
		CreatedAt:   creationTime,
		UpdatedAt:   updateTime,
		Visibility:  constants2.VisibilityPrivate,
	}

	got := converter.ConvertListToModel(entity)

	if !reflect.DeepEqual(got, expectedModel) {
		t.Errorf("ConvertListToModel() = %v, want %v", got, expectedModel)
	}
}

func TestConvertListToEntity(t *testing.T) {
	converter := lists.NewConverter()
	creationTime := time.Now()
	updateTime := time.Now()

	tags, err := json.Marshal([]constants2.TagType{constants2.TagWork, constants2.TagFinance})
	require.NoError(t, err)

	model := models.List{
		ID:          "123",
		Name:        "Test List",
		Description: "This is a test list",
		OwnerID:     "owner123",
		SharedWith:  []string{"user1", "user2"},
		Tags:        json.RawMessage(string(tags)),
		CreatedAt:   creationTime,
		UpdatedAt:   updateTime,
		Visibility:  constants2.VisibilityPrivate,
	}

	expectedEntity := lists.Entity{
		ID:          "123",
		Name:        "Test List",
		Description: "This is a test list",
		OwnerID:     "owner123",
		SharedWith:  []string{"user1", "user2"},
		Tags:        pkg.NewValidNullableString(string(tags)),
		CreatedAt:   creationTime,
		UpdatedAt:   updateTime,
		Visibility:  constants2.VisibilityPrivate,
	}

	got := converter.ConvertListToEntity(model)

	if !reflect.DeepEqual(got, expectedEntity) {
		t.Errorf("ConvertListToEntity() = %v, want %v", got, expectedEntity)
	}
}
