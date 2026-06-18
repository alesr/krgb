package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alesr/krgb/xerrors"
)

type Settings struct {
	EffectIdx  int `json:"effect_idx"`
	Hue        int `json:"hue"`
	Sat        int `json:"sat"`
	Brightness int `json:"brightness"`
	Speed      int `json:"speed"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", xerrors.Wrap(xerrors.ConfigHomeDir, "cannot determine home directory", err)
	}
	return filepath.Join(home, ".config", "krgb", "settings.json"), nil
}

func Defaults() *Settings {
	return &Settings{
		EffectIdx:  0,
		Hue:        180,
		Sat:        200,
		Brightness: 128,
		Speed:      64,
	}
}

func Load() *Settings {
	s := Defaults()

	p, err := configPath()
	if err != nil {
		return s
	}

	f, err := os.Open(p)
	if err != nil {
		return s
	}

	defer f.Close()
	if err := json.NewDecoder(f).Decode(s); err != nil {
		return Defaults()
	}
	return s
}

func (s *Settings) Save() error {
	p, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return xerrors.Wrap(xerrors.ConfigSave, "create config directory failed", err)
	}

	f, err := os.Create(p)
	if err != nil {
		return xerrors.Wrap(xerrors.ConfigSave, "create config file failed", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(s); err != nil {
		return xerrors.Wrap(xerrors.ConfigSave, "write config file failed", err)
	}
	return nil
}
