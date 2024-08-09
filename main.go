package main

import (
	"fmt"
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

func move(filename, baseFolder string, folders []string) {
	if isFolder(filename, folders) {
		return
	}
	folder := getFolderName(filepath.Ext(filename))
	createFolder(baseFolder, folder, folders)
	src := fmt.Sprintf("%s/%s", baseFolder, filename)
	dest := fmt.Sprintf("%s/%s/%s", baseFolder, folder, filename)
	os.Rename(src, dest)
	alert(fmt.Sprintf("New file to: %s", folder), fmt.Sprintf("Dest: %s", dest))
}

func startWatching(baseFolder string, folders []string) {
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
					move(filename, baseFolder, folders)
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

func getFolderName(extension string) string {
	switch extension {
	case ".txt", ".pdf", ".doc", ".docx":
		return "documents"

	case ".jpg", ".jpeg", ".png":
		return "images"

	case ".xlsx", ".xls", ".csv":
		return "datasets"

	case ".iso":
		return "software"

	case ".py", ".js", ".go":
		return "code"
	default:
		return "other"
	}
}

func main() {
	folders := []string{"documents", "images", "datasets", "code", "software", "other"}

	if len(os.Args) != 2 {
		fmt.Println("Folder is missing")
		os.Exit(1)
	}

	dir := os.Args[1]

	info, err := os.Stat(dir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if !info.IsDir() {
		fmt.Printf("%s is not a directory\n", dir)
		os.Exit(1)
	}

	startWatching(dir, folders)
}
