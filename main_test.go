package main

import (
	"bytes"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/apeyroux/gosm"
)

func TestNewVectorFromMapParams(t *testing.T) {
	tt := []struct {
		vectorParams map[string]string
		vector       *gosm.Tile
		mustBeValid  bool
	}{
		{map[string]string{"x": "1", "y": "1", "z": "1"}, &gosm.Tile{X: 1, Y: 1, Z: 1}, true},
		{map[string]string{"y": "1", "z": "1"}, nil, false},
		{map[string]string{"x": "", "z": "1"}, nil, false},
		{map[string]string{"x": "1", "y": "1"}, nil, false},
		{map[string]string{"x": "", "y": "1", "z": "1"}, nil, false},
		{map[string]string{"x": "1", "y": "", "z": "1"}, nil, false},
		{map[string]string{"x": "1", "y": "1", "z": ""}, nil, false},
	}

	for _, testData := range tt {
		_, err := newVectorFromMapParams(testData.vectorParams)

		if err != nil && testData.mustBeValid == true {
			t.Errorf("Error during vector creation, %s", err.Error())
		}
	}
}

func TestLoadTileFromMapProvider(t *testing.T) {
	payload := []byte{1, 2, 3}

	osmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(payload)
	}))

	b, err := loadTileFromMapProvider(gosm.Tile{X: 1, Y: 1, Z: 1}, osmServer.URL+"/%d/%d/%d")

	if err != nil {
		t.Errorf("Error during tile loading from provider. Cause: %s", err.Error())
	}

	if bytes.Compare(b, payload) != 0 {
		t.Errorf("Tile payload mismatch")
	}
}

func TestWriteTileInDisk(t *testing.T) {
	dir := "test-tiles"
	filename := dir + "/" + "test" + strconv.Itoa(rand.Intn(1000))

	createDirIfNotExists(dir)

	defer os.RemoveAll(dir)

	err := writeTileInDisk(filename, []byte{})

	if err != nil {
		t.Errorf("Error during tile writing. Cause %s", err.Error())
		return
	}

	fileExists := checkIfFileExists(filename)

	if !fileExists {
		t.Errorf("File not found")
	}
}

func TestCheckIfFileExists(t *testing.T) {
	filenames := []string{
		"test" + strconv.Itoa(rand.Intn(1000)),
		"test" + strconv.Itoa(rand.Intn(1000)),
		"test" + strconv.Itoa(rand.Intn(1000)),
		"test" + strconv.Itoa(rand.Intn(1000)),
		"test" + strconv.Itoa(rand.Intn(1000)),
		"test" + strconv.Itoa(rand.Intn(1000)),
		"test" + strconv.Itoa(rand.Intn(1000)),
	}

	for _, f := range filenames {
		defer os.Remove(f)

		os.WriteFile(f, []byte{}, 0644)

		if !checkIfFileExists(f) {
			t.Errorf("File %s not found", f)
		}
	}
}
