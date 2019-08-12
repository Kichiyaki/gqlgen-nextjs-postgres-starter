package serializer

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/sessions"
	_sessions "github.com/kichiyaki/graphql-starter/backend/sessions"
)

func NewJSONSerializer() _sessions.Serializer {
	return jsonSerializer{}
}

// jsonSerializer encode the session map to JSON.
type jsonSerializer struct{}

// Serialize to JSON. Will err if there are unmarshalable key values
func (s jsonSerializer) Serialize(ss *sessions.Session) ([]byte, error) {
	m := make(map[string]interface{}, len(ss.Values))
	for k, v := range ss.Values {
		ks, ok := k.(string)
		if !ok {
			err := fmt.Errorf("Non-string key value, cannot serialize session to JSON: %v", k)
			return nil, err
		}
		m[ks] = v
	}
	return json.Marshal(m)
}

// Deserialize back to map[string]interface{}
func (s jsonSerializer) Deserialize(d []byte, ss *sessions.Session) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(d, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		ss.Values[k] = v
	}

	return nil
}
