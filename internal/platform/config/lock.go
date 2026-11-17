package config

import (
	_ "embed"
	"encoding/json"
	"errors"
)

//go:embed sources.lock
var defaultLock []byte

// ArchivePin is one downloaded archive.
type ArchivePin struct {
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}

// BinaryPin is one per-arch official binary pin.
type BinaryPin struct {
	Version string     `json:"version"`
	AMD64   ArchivePin `json:"amd64"`
	ARM64   ArchivePin `json:"arm64"`
}

// SwiftBuildPin is the idev Swift Build source pin.
type SwiftBuildPin struct {
	Repo string `json:"repo"`
	SHA  string `json:"sha"`
}

// AppleDevLock is the maintainer provenance lock.
type AppleDevLock struct {
	Repo      string    `json:"repo"`
	Commit    string    `json:"commit"`
	XtoolRepo string    `json:"xtool_repo"`
	XtoolSHA  string    `json:"xtool_sha"`
	OamRepo   string    `json:"oam_repo"`
	OamSHA    string    `json:"oam_sha"`
	SDKStamp  string    `json:"sdk_stamp"`
	Rcodesign BinaryPin `json:"rcodesign"`
	Ipsw      BinaryPin `json:"ipsw"`
}

// IdevLock is the maintainer pin for the React Native and Expo stack.
type IdevLock struct {
	Pymobiledevice3 string        `json:"pymobiledevice3"`
	Cocoapods       string        `json:"cocoapods"`
	Erb             string        `json:"erb"`
	SwiftBuild      SwiftBuildPin `json:"swift_build"`
}

// Lock is the pinned upstream provenance and idev stack.
type Lock struct {
	AppleDev AppleDevLock `json:"apple_dev"`
	Idev     IdevLock     `json:"idev"`
}

// LoadLock reads the embedded maintainer lock.
func LoadLock() (Lock, error) {
	var lock Lock
	if err := json.Unmarshal(defaultLock, &lock); err != nil {
		return Lock{}, err
	}
	if lock.AppleDev.Commit == "" || lock.AppleDev.XtoolSHA == "" {
		return Lock{}, errors.New("lock: apple-dev commit or xtool SHA is empty")
	}
	if lock.Idev.Pymobiledevice3 == "" || lock.Idev.Cocoapods == "" || lock.Idev.Erb == "" {
		return Lock{}, errors.New("lock: idev component pin is empty")
	}
	if lock.AppleDev.SDKStamp == "" || lock.AppleDev.Rcodesign.Version == "" || lock.AppleDev.Ipsw.Version == "" {
		return Lock{}, errors.New("lock: apple-dev binary pin is empty")
	}
	if lock.Idev.SwiftBuild.SHA == "" {
		return Lock{}, errors.New("lock: Swift Build pin is empty")
	}
	return lock, nil
}
