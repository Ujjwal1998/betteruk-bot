package client

import "fmt"

type Category struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	HasChildren  bool   `json:"has_children"`
}

type Activity struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type categoriesResponse struct {
	Data []Category `json:"data"`
}

type categoryDetail struct {
	Name     string     `json:"name"`
	Slug     string     `json:"slug"`
	Children []Activity `json:"children"`
}

type categoryDetailResponse struct {
	Data categoryDetail `json:"data"`
}

// GetCategories lists bookable activity categories at a venue.
func (c *Client) GetCategories(venueSlug string) ([]Category, error) {
	var result categoriesResponse
	path := fmt.Sprintf("/api/activities/venue/%s/categories", venueSlug)
	if err := c.adminGET(path, &result); err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}
	return result.Data, nil
}

// GetCategoryActivities returns activities within a category (e.g. sports under Sports Hall).
func (c *Client) GetCategoryActivities(venueSlug, categorySlug string) ([]Activity, error) {
	var result categoryDetailResponse
	path := fmt.Sprintf("/api/activities/venue/%s/categories/%s", venueSlug, categorySlug)
	if err := c.adminGET(path, &result); err != nil {
		return nil, fmt.Errorf("get category %q: %w", categorySlug, err)
	}
	return result.Data.Children, nil
}
