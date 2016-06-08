// Package archive handles reading assets that have been packed into a single
// file.
// TODO(hurricanerix): currently the path is hard coded, this should be changed
//                     to make it able to ready any supported file.
package archive

import (
	"archive/tar"
	"bytes"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Get file from archive
func Get(name string) ([]byte, error) {
	d, err := importPathToDir("github.com/hurricanerix/shade")
	if err != nil {
		return nil, err
	}

	filepath := fmt.Sprintf("%s/assets.tar", d)
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)

	// Open the tar archive for reading.
	r := bytes.NewReader(buf.Bytes())
	tr := tar.NewReader(r)

	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		if name == hdr.Name {
			return ioutil.ReadAll(tr)
		}
	}
	return nil, fmt.Errorf("Could not find '%s'", name)
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}
