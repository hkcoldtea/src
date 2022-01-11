package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// IsDir reports whether d is a directory.
func IsDir(d string) (y bool) {
	if fi, err := os.Stat(d); err == nil {
		if fi.IsDir() {
			y = true
		}
	}
	return
}

// ReadJson reads and parses the JSON-encoded contents of the named file and
// stores the result in the value pointed to by v.
// Returns an error if the named file cannot be read or correctly parsed.
func ReadJson(name string, v interface{}) error {
	if b, err := ioutil.ReadFile(name); err != nil {
		return err
	} else {
		if err := json.Unmarshal(b, &v); err != nil {
			return fmt.Errorf("parse %s: %s", name, err)
		}
	}
	return nil
}
