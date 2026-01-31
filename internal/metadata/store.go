package metadata

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func WriteJSON(path string, name string, data any) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(path, name), b, 0644)
}
