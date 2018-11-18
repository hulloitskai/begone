package messenger

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// Config holds configuration options for the Engine.
type Config struct {
	User      string `json:"user"`
	Pass      string `json:"pass,omitempty"`
	ConvoID   string `json:"-"`
	Namespace string `json:"-"`

	// Private, cached properties.
	convoURL, fileName string
}

// FileName gets the name of the config file based on cfg.NameSpace.
//
// If cfg.NameSpace is empty, an empty string will be returned.
func (cfg *Config) FileName() string {
	if cfg.fileName == "" {
		cfg.fileName = nstail(cfg.Namespace)
	}

	return cfg.fileName
}

// NewConfig makes a new Config.
func NewConfig(user, pass, convoID string) *Config {
	return &Config{User: user, Pass: pass, ConvoID: convoID}
}

// MessengerURL constructs a 'messenger.com' URL based on the Config's ConvoID.
func (cfg *Config) MessengerURL() string {
	if cfg.convoURL == "" {
		cfg.convoURL = "https://messenger.com/t/" + cfg.ConvoID
	}

	return cfg.convoURL
}

// CheckLogin returns an error if the Config is missing the necessary
// credentials to login to FB messenger.
func (cfg *Config) CheckLogin() error {
	if cfg.User == "" {
		return errors.New("messenger: missing username (User)")
	}
	if cfg.Pass == "" {
		return errors.New("messenger: missing password (Pass)")
	}

	return nil
}

// WriteToFile writes the Config to the file located at path. Any existing data
// in the file will be overridden.
//
// If the file does not exist, it will be created (the file directory, however,
// must already exist).
func (cfg *Config) WriteToFile(path string) error {
	file, err := os.Create(path) // works with existing files, too
	if err != nil {
		return fmt.Errorf("messenger: failed to create config file: %v", err)
	}
	defer file.Close()

	// Create an encoded Config that obfuscates the password (by encoding it to
	// base64). A copy of the current Config is created in order to prevent
	// modifying the current password.
	encodedCfg := *cfg
	encodedCfg.Pass = base64.StdEncoding.EncodeToString([]byte(cfg.Pass))

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err = enc.Encode(&encodedCfg); err != nil {
		return fmt.Errorf("messenger: failed to encode config as JSON: %v", err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("messenger: could not close file: %v", err)
	}
	return nil
}

// ReadConfigFile reads the file located at path into cfg.
//
// This does not automatically fill cfg.Namespace; the caller must fill this
// field manually.
func ReadConfigFile(cfg *Config, path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}

		return fmt.Errorf("messenger: could not open config file for reading: %v",
			err)
	}
	defer file.Close()

	// Read JSON config file to cfg.
	dec := json.NewDecoder(file)
	if err = dec.Decode(cfg); err != nil {
		return fmt.Errorf("messenger: failed to decode config from JSON: %v", err)
	}

	// Decode cfg.Pass from base64.
	var (
		encPass = []byte(cfg.Pass)
		decPass = make([]byte, len(encPass))
	)
	if _, err = base64.StdEncoding.Decode(decPass, encPass); err != nil {
		return fmt.Errorf("messenger: error while decoding pass from base64: %v",
			err)
	}
	cfg.Pass = string(decPass)

	// Close file.
	if err = file.Close(); err != nil {
		return fmt.Errorf("messenger: could not close file: %v", err)
	}
	return nil
}

// configPath returns the platform-specific config file path, provided a
// namespace.
//
// On Windows, the namespace folder is created in %APPDATA%\{{namespace}} if
// it does not already exist.
func configPath(namespace string) (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("messenger: cannot get current user: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		roamingPath, ok := os.LookupEnv("APPDATA")
		if !ok {
			roamingPath = filepath.Join(u.HomeDir, "AppData", "Roaming")
		}

		namespaceDir := filepath.Join(roamingPath, namespace)
		if _, err = os.Stat(namespaceDir); os.IsNotExist(err) {
			if err = os.Mkdir(namespaceDir, 0755); err != nil {
				return "", fmt.Errorf("messenger: failed to create config "+
					"directory: %v", err)
			}
		}

		return filepath.Join(namespaceDir, nstail(namespace)), nil

	case "linux":
		configHome, ok := os.LookupEnv("XDG_CONFIG_HOME")
		if !ok {
			configHome = filepath.Join(u.HomeDir)
		}
		return filepath.Join(configHome, nstail(namespace)), nil

	default:
		return filepath.Join(u.HomeDir, nstail(namespace)), nil
	}
}

// FilePath returns the Config's expected filepath on disk, based on its
// namespace.
func (cfg *Config) FilePath() (string, error) {
	return configPath(cfg.Namespace)
}

// Save saves the Config to path specified by cfg.FilePath.
//
// Returns the resulting path to the config file.
func (cfg *Config) Save() (path string, err error) {
	if cfg.Namespace == "" {
		return "", errors.New("messenger: cannot save a Config with an empty " +
			"Namespace")
	}

	path, err = cfg.FilePath()
	if err != nil {
		return "", err
	}

	cfg.WriteToFile(path)
	return path, nil
}

// ReadConfig reads a config with the specificed name from the user config
// directory into cfg.
//
// Additionally, cfg.Namespace will be set to the provided namespace.
func ReadConfig(namespace string, cfg *Config) error {
	path, err := configPath(namespace)
	if err != nil {
		return err
	}

	if err = ReadConfigFile(cfg, path); err != nil {
		return err
	}

	cfg.Namespace = namespace
	return nil
}

// RemoveFile removes the config file(s) (if they exists), and returns the paths
// to the removed files.
//
// If no config file existed, the returned path will be empty.
func (cfg *Config) RemoveFile() (removed []string, err error) {
	path, err := cfg.FilePath()
	if err != nil {
		return nil, fmt.Errorf("messenger: failed to get config file path: %v", err)
	}

	if err = os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("messenger: failed to remove config file at "+
				"'%s': %v", path, err)
		}
	} else {
		removed = []string{path}
	}

	if runtime.GOOS == "windows" {
		dirPath := filepath.Dir(path)
		if err = os.Remove(dirPath); err != nil {
			return nil, fmt.Errorf("messenger: failed to remove config directory at "+
				"'%s': %v", path, err)
		}

		removed = append(removed, dirPath)
	}

	return removed, nil
}
