package config

import (
	"encoding/base64"
	"encoding/json"
	"os"

	ess "github.com/unixpickle/essentials"
)

// Save saves a Config to a file (located at Filepath).
// Returns the path to the resulting config file.
func Save(cfg *Config) (string, error) {
	path, err := Filepath()
	if err != nil {
		return "", err
	}

	// Create an encoded copy of cfg, with a base64-encoded password.
	var (
		encCfg  = *cfg
		encPass = []byte(encCfg.Password)
	)
	encCfg.Password = base64.StdEncoding.EncodeToString(encPass)

	// Create or open existing file for writing.
	file, err := os.Create(path)
	if err != nil {
		return "", ess.AddCtx("config: creating file", err)
	}
	defer file.Close()

	// Write to file.
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	enc.Encode(encCfg)

	// Close file.
	if err = file.Close(); err != nil {
		return "", ess.AddCtx("config: closing file", err)
	}
	return path, nil
}
