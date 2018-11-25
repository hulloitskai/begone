package config

import (
	"encoding/base64"
	"encoding/json"
	"os"

	ess "github.com/unixpickle/essentials"
)

// Load loads a Config from a file (located at Filepath).
//
// An empty Config will be returned if no such file is found.
func Load() (*Config, error) {
	path, err := Filepath()
	if err != nil {
		return nil, err
	}

	// Open file for reading.
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return new(Config), nil
		}
		return nil, ess.AddCtx("config: opening file", err)
	}
	defer file.Close()

	// Read file.
	var (
		cfg = new(Config)
		dec = json.NewDecoder(file)
	)
	if err = dec.Decode(cfg); err != nil {
		return nil, ess.AddCtx("config: decoding file", err)
	}

	// Decode cfg.Password from base64.
	decpass, err := base64.StdEncoding.DecodeString(cfg.Password)
	if err != nil {
		return nil, ess.AddCtx("config: decoding password", err)
	}
	cfg.Password = string(decpass)

	// Close file.
	if err = file.Close(); err != nil {
		return nil, ess.AddCtx("config: closing file", err)
	}
	return cfg, nil
}
