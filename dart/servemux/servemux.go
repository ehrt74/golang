package dart

import (
	"net/http"
	"os"
	"bufio"
	"strings"
	"log"
)

var pubCache map[string]string
func init() {
	initPubCache()
}

func initPubCache() {
	pubCache = make(map[string]string)
	file, err := os.Open(".packages")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		pubCache[parts[0]] = parts[len(parts)-1]
	}
}

// BasePath returns the path to the web/ directory. If the Dartium browser is used, this
// is "web/", otherwise "build/web/"
func BasePath(r *http.Request) (path string) {
	if !strings.Contains(r.UserAgent(), "Dart") {
		path = "build/web/"
	} else {
		path = "web/"
	}
	return
}

// NewServeMux returns a *http.ServeMux with handlers for dart packages
func NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/packages/", packagesHandler)
	for _, dir := range []string{"web", "lib"} {
		mux.Handle("/"+dir+"/", http.FileServer(http.Dir(dir)))
	}
	return mux
}

func packagesHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	fileName := strings.Join(parts[2:len(parts)], "/")
	packageName := parts[1]
	dirName, ok := pubCache[packageName]
	if ok {
		http.ServeFile(w, r, dirName+fileName)
		return
	}
	path := "lib/" + fileName
	if _, err := os.Stat(path); err == nil {
		http.ServeFile(w, r, path)
		return
	}
	log.Printf("%v not found", r.URL.Path)
}
