package client

import (
	"encoding/json"
	"testing"
)

const sampleCategoriesJSON = `{"data":[{"name":"Sports Hall Activities","slug":"sports-hall-activities","has_children":true}]}`

const sampleCategoryDetailJSON = `{"data":{"name":"Sports Hall Activities","slug":"sports-hall-activities","children":[{"name":"Badminton 40min","slug":"badminton-40min"}]}}`

func TestParseCategoriesResponse(t *testing.T) {
	var result categoriesResponse
	if err := json.Unmarshal([]byte(sampleCategoriesJSON), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 || result.Data[0].Slug != "sports-hall-activities" {
		t.Fatalf("unexpected: %+v", result.Data)
	}
}

func TestParseCategoryDetailResponse(t *testing.T) {
	var result categoryDetailResponse
	if err := json.Unmarshal([]byte(sampleCategoryDetailJSON), &result); err != nil {
		t.Fatal(err)
	}
	children := result.Data.Children
	if len(children) != 1 || children[0].Slug != "badminton-40min" {
		t.Fatalf("unexpected children: %+v", children)
	}
}
