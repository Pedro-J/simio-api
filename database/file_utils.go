package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var lock sync.Mutex

func saveEntityOnFile(filename string, object interface{}) error {
	createDirIfNotExist(getDefaultDirectory())

	filePath := getDefaultDirectory() + filename

	err := save(filePath, object)

	if err != nil {
		log.Printf("Error on saving entity in file. Details: %s", err)
		return fmt.Errorf("UNEXPECTED_ERROR_ON_SAVE")
	}
	log.Printf("Entity %s has been saved successfully", filename)

	return nil
}

func createDefaultDirectory() {
	createDirIfNotExist(getDefaultDirectory())
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func getDefaultDirectory() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("%s", err)
	}

	return currentDir + "/database/data/simios/"
}

func LoadAll(dir string) (map[string]SimioEntity, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		log.Printf("Error on loading entity in file. Details: %s", err)
		return nil, fmt.Errorf("UNEXPECTED_ERROR_ON_LOAD")
	}

	if len(files) > 1 {

		simios := make([]SimioEntity, len(files)-1)

		for current := 1; current < len(files); current++ {
			load(files[current], &simios[current-1])
		}

		data := make(map[string]SimioEntity)

		for _, simio := range simios {
			data[simio.ID] = simio
		}

		log.Printf("Files Loaded successfully. DB Size = %v", len(data))

		return data, nil
	}

	log.Printf("No files found to be loaded")

	return nil, nil
}

func save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := marshal(v)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

func load(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return unmarshal(f, v)
}

var marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

var unmarshal = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func checkFileExist(filename string) bool {
	filePath := getDefaultDirectory() + filename
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}
