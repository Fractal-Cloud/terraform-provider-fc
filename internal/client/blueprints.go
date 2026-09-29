package fractalCloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type FractalId struct {
	ResourceGroupId ResourceGroupId
	Name            string
	Version         string
}

func (id *FractalId) ToString() string {
	return id.ResourceGroupId.ToString() + "/" + id.Name + ":" + id.Version
}

func (id *ResourceGroupId) ToString() string {
	return id.Type + "/" + id.OwnerId + "/" + id.ShortName
}

// GetBlueprint - Returns specific blueprint (fractal definition).
func (c *Client) GetBlueprint(ctx context.Context, id FractalId) (*Blueprint, error) {
	path := c.blueprintPath(id)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("building GET request for fractal %q: %w", id.ToString(), err)
	}

	resCode, body, err := c.doRequest(ctx, req, []int{200, 404})
	if err != nil {
		return nil, fmt.Errorf("fetching fractal %q: %w", id.ToString(), err)
	}

	if resCode == 404 {
		c.logDebug(fmt.Sprintf("fractal %q not found", id.ToString()))
		return nil, nil
	}

	// UseNumber keeps numeric parameters as their literal text ("1.10" stays
	// "1.10") instead of round-tripping through float64.
	blueprint := BlueprintInternal{}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(&blueprint); err != nil {
		return nil, fmt.Errorf("decoding fractal %q response: %w", id.ToString(), err)
	}

	normalizedComponents := make([]Component, len(blueprint.Components))
	for i, component := range blueprint.Components {
		links := make([]ComponentLink, len(component.Links))
		for j, link := range component.Links {
			links[j] = ComponentLink{
				ComponentId: link.ComponentId,
				Settings:    mapAnyToMapStringJSON(c.Logger, link.Settings),
			}
		}
		normalizedComponents[i] = Component{
			Id:                component.Id,
			Type:              component.Type,
			DisplayName:       &component.DisplayName,
			Description:       &component.Description,
			Version:           &component.Version,
			IsLocked:          &component.IsLocked,
			RecreateOnFailure: &component.RecreateOnFailure,
			Parameters:        mapAnyToMapStringJSON(c.Logger, component.Parameters),
			DependenciesIds:   component.DependenciesIds,
			Links:             links,
			OutputFields:      component.OutputFields,
		}
	}

	return &Blueprint{
		FractalId:   blueprint.FractalId,
		IsPrivate:   blueprint.IsPrivate,
		Status:      blueprint.Status,
		ReasonCode:  blueprint.ReasonCode,
		Description: blueprint.Description,
		Components:  normalizedComponents,
		CreatedAt:   blueprint.CreatedAt,
	}, nil
}

type ComponentLinkInternal struct {
	ComponentId string                 `json:"componentId"`
	Settings    map[string]interface{} `json:"settings"`
}

type ComponentInternal struct {
	Id                string                  `json:"id"`
	Type              string                  `json:"type"`
	DisplayName       string                  `json:"displayName"`
	Description       string                  `json:"description"`
	Version           string                  `json:"version"`
	IsLocked          bool                    `json:"locked"`
	RecreateOnFailure bool                    `json:"recreateOnFailure"`
	Parameters        map[string]interface{}  `json:"parameters"`
	DependenciesIds   []string                `json:"dependencies"`
	Links             []ComponentLinkInternal `json:"links"`
	OutputFields      []string                `json:"outputFields"`
}

type BlueprintInternal struct {
	FractalId   string              `json:"fractalId"`
	IsPrivate   bool                `json:"isPrivate"`
	Status      string              `json:"status"`
	ReasonCode  string              `json:"reasonCode"`
	Description string              `json:"description"`
	Components  []ComponentInternal `json:"components"`
	CreatedAt   string              `json:"created"`
}

type CreateBlueprintCommandRequestBody struct {
	Description string      `json:"description"`
	IsPrivate   bool        `json:"isPrivate"`
	Components  []Component `json:"components"`
}

func (c *Client) CreateBlueprint(ctx context.Context, id FractalId, description string, isPrivate bool, components []Component) error {
	requestBody := CreateBlueprintCommandRequestBody{
		Description: description,
		IsPrivate:   isPrivate,
		Components:  components,
	}

	rb, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("encoding create fractal %q request: %w", id.ToString(), err)
	}

	path := c.blueprintPath(id)

	req, err := http.NewRequest("POST", path, strings.NewReader(string(rb)))
	if err != nil {
		return fmt.Errorf("building POST request for fractal %q: %w", id.ToString(), err)
	}

	_, _, err = c.doRequest(ctx, req, []int{200, 202})
	if err != nil {
		return fmt.Errorf("creating fractal %q: %w", id.ToString(), err)
	}
	return nil
}

func (c *Client) UpdateBlueprint(ctx context.Context, id FractalId, description string, isPrivate bool, components []Component) error {
	requestBody := CreateBlueprintCommandRequestBody{
		Description: description,
		IsPrivate:   isPrivate,
		Components:  components,
	}

	rb, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("encoding update fractal %q request: %w", id.ToString(), err)
	}

	path := c.blueprintPath(id)

	req, err := http.NewRequest("PUT", path, strings.NewReader(string(rb)))
	if err != nil {
		return fmt.Errorf("building PUT request for fractal %q: %w", id.ToString(), err)
	}

	_, _, err = c.doRequest(ctx, req, []int{200, 202})
	if err != nil {
		return fmt.Errorf("updating fractal %q: %w", id.ToString(), err)
	}
	return nil
}

func (c *Client) DeleteBlueprint(ctx context.Context, id FractalId) error {
	path := c.blueprintPath(id)

	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		return fmt.Errorf("building DELETE request for fractal %q: %w", id.ToString(), err)
	}

	_, _, err = c.doRequest(ctx, req, []int{200, 202, 404})
	if err != nil {
		return fmt.Errorf("deleting fractal %q: %w", id.ToString(), err)
	}
	return nil
}

func (c *Client) blueprintPath(id FractalId) string {
	rg := id.ResourceGroupId
	return fmt.Sprintf("%s/blueprints/%s/%s/%s/%s/%s",
		c.HostURL, rg.Type, rg.OwnerId, rg.ShortName, id.Name, id.Version)
}

func mapAnyToMapStringJSON(logger *ClientLogger, in map[string]interface{}) map[string]string {
	out := make(map[string]string, len(in))

	for k, v := range in {
		if s, ok := v.(string); ok {
			out[k] = s
		} else {
			b, err := json.Marshal(v)
			if err == nil {
				out[k] = string(b)
			} else if logger != nil && logger.Warning != nil {
				logger.Warning(fmt.Sprintf("marshal key %q: %v", k, err))
			}
		}
	}

	return out
}
