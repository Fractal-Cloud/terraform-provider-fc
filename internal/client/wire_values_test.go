package fractalCloud

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncodeWireValue(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain string stays a string", "etl", `"etl"`},
		{"numeric text stays a string", "8080", `"8080"`},
		{"boolean text stays a string", "true", `"true"`},
		{"empty string stays a string", "", `""`},
		{"object is sent as JSON", `{"spark.x":"1"}`, `{"spark.x":"1"}`},
		{"array is sent as JSON", `["a","b"]`, `["a","b"]`},
		{"surrounding whitespace is dropped", "  [1, 2] ", `[1, 2]`},
		{"env secret reference is sent as an object", `{"$envSecret":"openai-key"}`, `{"$envSecret":"openai-key"}`},
		{"invalid JSON stays a string", `{not json`, `"{not json"`},
		{"brace-prefixed prose stays a string", `{a: 1}`, `"{a: 1}"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(encodeWireValue(tt.in)); got != tt.want {
				t.Errorf("encodeWireValue(%q) = %s, want %s", tt.in, got, tt.want)
			}
		})
	}
}

func TestComponentMarshalJSON_TypesParametersAndLinkSettings(t *testing.T) {
	c := Component{
		Id:   "job",
		Type: "BigData.PaaS.DataProcessingJob",
		Parameters: map[string]string{
			"parameters": `["--date","today"]`,
			"secret":     `{"$envSecret":"api-key"}`,
			"maxRetries": "3",
		},
		Links: []ComponentLink{
			{ComponentId: "lake", Settings: map[string]string{"purpose": "raw", "token": `{"$envSecret":"lake-token"}`}},
		},
	}

	b, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got struct {
		Parameters map[string]any `json:"parameters"`
		Links      []struct {
			ComponentId string         `json:"componentId"`
			Settings    map[string]any `json:"settings"`
		} `json:"links"`
	}
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if _, ok := got.Parameters["parameters"].([]any); !ok {
		t.Errorf("parameters.parameters = %T, want JSON array", got.Parameters["parameters"])
	}
	secret, ok := got.Parameters["secret"].(map[string]any)
	if !ok || secret["$envSecret"] != "api-key" {
		t.Errorf("parameters.secret = %#v, want {$envSecret: api-key}", got.Parameters["secret"])
	}
	if got.Parameters["maxRetries"] != "3" {
		t.Errorf("parameters.maxRetries = %#v, want string \"3\"", got.Parameters["maxRetries"])
	}
	if len(got.Links) != 1 || got.Links[0].ComponentId != "lake" {
		t.Fatalf("links = %#v", got.Links)
	}
	if got.Links[0].Settings["purpose"] != "raw" {
		t.Errorf("link purpose = %#v, want \"raw\"", got.Links[0].Settings["purpose"])
	}
	if _, ok := got.Links[0].Settings["token"].(map[string]any); !ok {
		t.Errorf("link token = %T, want JSON object", got.Links[0].Settings["token"])
	}
}

// A structured value must come back as the exact string Terraform holds, or
// every refresh would report a diff.
func TestBlueprintRoundTrip_StructuredValuesAreStable(t *testing.T) {
	sent := map[string]string{
		"secret":    `{"$envSecret":"api-key"}`,
		"sparkConf": `{"a":"1","b":"2"}`,
		"libraries": `["x","y"]`,
		"ratio":     "1.10",
	}

	var stored []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req struct {
				Components []map[string]any `json:"components"`
			}
			body, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(body, &req); err != nil {
				t.Fatalf("decoding request: %v", err)
			}
			// The server stores the typed values and returns them, with object
			// keys in a different order, as a jsonb column would.
			params := req.Components[0]["parameters"].(map[string]any)
			params["ratio"] = json.Number("1.10")
			params["sparkConf"] = map[string]any{"b": "2", "a": "1"}
			stored, _ = json.Marshal(map[string]any{"components": req.Components})
			w.WriteHeader(http.StatusAccepted)
		case http.MethodGet:
			w.Write(stored)
		}
	}))
	defer server.Close()

	c := &Client{HostURL: server.URL, HTTPClient: server.Client()}
	id := FractalId{ResourceGroupId: ResourceGroupId{Type: "Personal", OwnerId: "o", ShortName: "bc"}, Name: "f", Version: "1"}

	if err := c.CreateBlueprint(t.Context(), id, "", false, []Component{{Id: "c", Type: "T", Parameters: sent}}); err != nil {
		t.Fatalf("create: %v", err)
	}
	bp, err := c.GetBlueprint(t.Context(), id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	got := bp.Components[0].Parameters
	for k, want := range sent {
		if got[k] != want {
			t.Errorf("parameter %q = %q after round trip, want %q", k, got[k], want)
		}
	}
}
