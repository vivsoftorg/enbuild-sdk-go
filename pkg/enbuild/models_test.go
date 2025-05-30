package enbuild_test

import (
	"encoding/json"
	"testing"

	"github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
)

func TestCatalogDataTimestampUnmarshal(t *testing.T) {
	t.Run("NumericTimestamps", func(t *testing.T) {
		jsonData := `{"createdOn": 1678886400, "updatedOn": 1678886401, "name": "TestCatalog", "_id": "cat1"}`
		var catalogData enbuild.CatalogData
		err := json.Unmarshal([]byte(jsonData), &catalogData)

		if err != nil {
			t.Fatalf("Failed to unmarshal CatalogData with numeric timestamps: %v", err)
		}

		expectedCreatedOn := int64(1678886400)
		if catalogData.CreatedOn != expectedCreatedOn {
			t.Errorf("Expected CreatedOn to be %d, got %d", expectedCreatedOn, catalogData.CreatedOn)
		}

		expectedUpdatedOn := int64(1678886401)
		if catalogData.UpdatedOn != expectedUpdatedOn {
			t.Errorf("Expected UpdatedOn to be %d, got %d", expectedUpdatedOn, catalogData.UpdatedOn)
		}
		if catalogData.Name != "TestCatalog" {
			t.Errorf("Expected Name to be TestCatalog, got %s", catalogData.Name)
		}
	})

	t.Run("OmittedTimestamps", func(t *testing.T) {
		jsonData := `{"name": "TestCatalogOmitted", "_id": "cat2"}` // Timestamps omitted
		var catalogData enbuild.CatalogData
		err := json.Unmarshal([]byte(jsonData), &catalogData)

		if err != nil {
			t.Fatalf("Failed to unmarshal CatalogData with omitted timestamps: %v", err)
		}

		if catalogData.CreatedOn != 0 {
			t.Errorf("Expected CreatedOn to be 0 (default for int64) when omitted, got %d", catalogData.CreatedOn)
		}

		if catalogData.UpdatedOn != 0 {
			t.Errorf("Expected UpdatedOn to be 0 (default for int64) when omitted, got %d", catalogData.UpdatedOn)
		}
	})

	t.Run("NullTimestamps", func(t *testing.T) {
		jsonData := `{"createdOn": null, "updatedOn": null, "name": "TestCatalogNull", "_id": "cat3"}`
		var catalogData enbuild.CatalogData
		err := json.Unmarshal([]byte(jsonData), &catalogData)

		if err != nil {
			t.Fatalf("Failed to unmarshal CatalogData with null timestamps: %v", err)
		}

		if catalogData.CreatedOn != 0 {
			t.Errorf("Expected CreatedOn to be 0 (default for int64) when null, got %d", catalogData.CreatedOn)
		}

		if catalogData.UpdatedOn != 0 {
			t.Errorf("Expected UpdatedOn to be 0 (default for int64) when null, got %d", catalogData.UpdatedOn)
		}
	})
}

func TestStackTimestampUnmarshal(t *testing.T) {
	t.Run("NumericTimestamps", func(t *testing.T) {
		// Ensure other potentially required fields for Stack are present for valid unmarshaling
		jsonData := `{"id": "stack1", "name": "TestStack", "createdOn": 1678886400, "updatedOn": 1678886401, "created_on": 1678886402}`
		var stack enbuild.Stack
		err := json.Unmarshal([]byte(jsonData), &stack)

		if err != nil {
			t.Fatalf("Failed to unmarshal Stack with numeric timestamps: %v", err)
		}

		expectedCreatedOn := int64(1678886400)
		if stack.CreatedOn != expectedCreatedOn {
			t.Errorf("Expected Stack CreatedOn to be %d, got %d", expectedCreatedOn, stack.CreatedOn)
		}

		expectedUpdatedOn := int64(1678886401)
		if stack.UpdatedOn != expectedUpdatedOn {
			t.Errorf("Expected Stack UpdatedOn to be %d, got %d", expectedUpdatedOn, stack.UpdatedOn)
		}

		expectedCreated_on := int64(1678886402)
		if stack.Created_on != expectedCreated_on {
			t.Errorf("Expected Stack Created_on to be %d, got %d", expectedCreated_on, stack.Created_on)
		}
		if stack.Name != "TestStack" {
			t.Errorf("Expected Name to be TestStack, got %s", stack.Name)
		}
	})

	t.Run("OmittedTimestamps", func(t *testing.T) {
		jsonData := `{"id": "stack2", "name": "TestStackOmitted"}` // Timestamps omitted
		var stack enbuild.Stack
		err := json.Unmarshal([]byte(jsonData), &stack)

		if err != nil {
			t.Fatalf("Failed to unmarshal Stack with omitted timestamps: %v", err)
		}

		if stack.CreatedOn != 0 {
			t.Errorf("Expected Stack CreatedOn to be 0 when omitted, got %d", stack.CreatedOn)
		}
		if stack.UpdatedOn != 0 {
			t.Errorf("Expected Stack UpdatedOn to be 0 when omitted, got %d", stack.UpdatedOn)
		}
		if stack.Created_on != 0 {
			t.Errorf("Expected Stack Created_on to be 0 when omitted, got %d", stack.Created_on)
		}
	})

	t.Run("NullTimestamps", func(t *testing.T) {
		jsonData := `{"id": "stack3", "name": "TestStackNull", "createdOn": null, "updatedOn": null, "created_on": null}`
		var stack enbuild.Stack
		err := json.Unmarshal([]byte(jsonData), &stack)

		if err != nil {
			t.Fatalf("Failed to unmarshal Stack with null timestamps: %v", err)
		}

		if stack.CreatedOn != 0 {
			t.Errorf("Expected Stack CreatedOn to be 0 when null, got %d", stack.CreatedOn)
		}
		if stack.UpdatedOn != 0 {
			t.Errorf("Expected Stack UpdatedOn to be 0 when null, got %d", stack.UpdatedOn)
		}
		if stack.Created_on != 0 {
			t.Errorf("Expected Stack Created_on to be 0 when null, got %d", stack.Created_on)
		}
	})
}
