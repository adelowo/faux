package vfiles

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/influx6/faux/hexwriter"
)

// ParseDir returns a new instance of all CSS files located within the provided directory.
func ParseDir(dir string, allowedExtensions []string) (map[string]string, error) {
	items := make(map[string]string)

	// Walk directory pulling contents into css items.
	if cerr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if cerr := walkDir(allowedExtensions, items, dir, path, info, err); cerr != nil {
			return cerr
		}

		return nil
	}); cerr != nil {
		return nil, cerr
	}

	return items, nil
}

// validExension returns true/false if the extension provide is a valid acceptable one
// based on the allowedExtensions string slice.
func validExtension(extensions []string, ext string) bool {
	for _, es := range extensions {
		if es != ext {
			continue
		}

		return true
	}

	return false
}

// walkDir adds the giving path if it matches certain criterias into the items map.
func walkDir(extensions []string, items map[string]string, root string, path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return nil
	}

	// Is file an exension we allow else skip.
	if len(extensions) != 0 && !validExtension(extensions, filepath.Ext(path)) {
		return nil
	}

	rel, err := filepath.Rel(root, path)
	if err != nil {
		return err
	}

	relFile, err := os.Open(path)
	if err != nil {
		return err
	}

	defer relFile.Close()

	var contents bytes.Buffer

	io.Copy(hexwriter.New(&contents), relFile)

	items[rel] = contents.String()
	return nil
}
