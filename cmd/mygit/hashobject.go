package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

func HashObject(filepath string, write bool) (objectname string, err error) {
	// read file once to both compute sha-1 sum and compress using zlib
	// with larger files, using os.Open might be better
	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	// https://alblue.bandlem.com/2011/08/git-tip-of-week-objects.html
	prefix := fmt.Sprintf("blob %d\x00", len(contents))
	data := append([]byte(prefix), contents...)

	hasher := sha1.New()
	hasher.Write(data)
	objectname = hex.EncodeToString(hasher.Sum(nil))

	compressed := bytes.NewBuffer(make([]byte, 0))
	zlibWriter := zlib.NewWriter(compressed)
	if _, err := zlibWriter.Write(data); err != nil {
		return "", err
	}
	if err := zlibWriter.Close(); err != nil {
		return "", err
	}

	if write {
		err := WriteObject(objectname, compressed.Bytes())
		if err != nil {
			return objectname, nil
		}
	}

	return objectname, nil
}
