package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gen2brain/beeep"
)

func alert(title, message string) {
	err := beeep.Alert(title, message, "warn.png")
	if err != nil {
		panic(err)
	}
}

func move(filename, baseFolder, jsonFile string) {
	extMap := getExtensionMap(jsonFile)
	folders := getFolderNames(extMap)
	if isFolder(filename, folders) {
		return
	}

	folder := getFolderName(extMap, filepath.Ext(filename))
	createFolder(baseFolder, folder, folders)
	src := fmt.Sprintf("%s/%s", baseFolder, filename)
	dest := fmt.Sprintf("%s/%s/%s", baseFolder, folder, filename)
	os.Rename(src, dest)
	alert(fmt.Sprintf("New file to: %s", folder), fmt.Sprintf("Dest: %s", dest))
}

func startWatching(baseFolder, jsonFile string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					filename := filepath.Base(event.Name)
					move(filename, baseFolder, jsonFile)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()
	err = watcher.Add(baseFolder)
	if err != nil {
		log.Fatal(err)
	}
	<-make(chan struct{})
}

func createFolder(baseFolder, name string, folders []string) {
	for _, v := range folders {
		if name == v {
			newFolder := fmt.Sprintf("%s/%s", baseFolder, v)
			if _, err := os.Stat(newFolder); os.IsNotExist(err) {
				os.Mkdir(newFolder, 0755)
				fmt.Printf("Folder created: %s", newFolder)
			}
		}
	}
}

func isFolder(name string, folders []string) bool {
	for _, v := range folders {
		if name == v {
			return true
		}
	}
	return false
}

func getExtensionMap(jsonFile string) map[string][]string {
	file, err := os.Open(jsonFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	var extensionMap map[string][]string
	err = json.Unmarshal(data, &extensionMap)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		os.Exit(1)
	}
	return extensionMap

}

func getFolderName(extensionMap map[string][]string, extension string) string {
	for i := range extensionMap {
		for _, v := range extensionMap[i] {
			if v == extension {
				return i
			}
		}
	}
	return "other"
}

func getFolderNames(extensionMap map[string][]string) []string {
	folders := make([]string, 0, len(extensionMap))
	for folder := range extensionMap {
		folders = append(folders, folder)
	}
	return folders

}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Folder or Map file is missing")
		os.Exit(1)
	}

	dir := os.Args[1]
	mapFile := os.Args[2]

	info, err := os.Stat(dir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if !info.IsDir() {
		fmt.Printf("%s is not a directory\n", dir)
		os.Exit(1)
	}
	startWatching(dir, mapFile)
}
