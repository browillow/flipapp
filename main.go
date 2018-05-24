package main

import (
	"fmt"
	"net/http"
	"os"
)

func fileForRequestExists(r *http.Request) (bool, error) {
	path := "./dist" + r.URL.Path
	if len(path) != 0 {
		lastChar := path[len(path)-1]
		if lastChar == '/' {
			path = path + "index.html"
		}
	}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func customHandler(fileHandler http.Handler, indexHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprintf("INFO [FLIPPERAPP] - Processing request [%v]", r.URL.Path))
		fileExists, err := fileForRequestExists(r)
		if err != nil {
			fmt.Println(fmt.Sprintf("ERROR [FLIPPERAPP] - Error while checking if file [%v] for request exists: %v", r.URL.Path, err.Error()))
			indexHandler.ServeHTTP(w, r)
		} else if fileExists {
			fileHandler.ServeHTTP(w, r)
		} else {
			indexHandler.ServeHTTP(w, r)
		}
	})
}

func indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/app/index.html")
	})
}

func main() {
	appPort := ":" + os.Getenv("FLIPPERAPP_SERVICE_LISTENER")
	if appPort == ":" {
		appPort = ":80"
	}
	fs := http.FileServer(http.Dir("dist"))
	mux := http.NewServeMux()
	mux.Handle("/", customHandler(fs, indexHandler()))
	err := http.ListenAndServe(appPort, mux)
	if err != nil {
		fmt.Println("SYSTEM ERROR [FLIPPERAPP]: Error while starting HTTP server -> ", err.Error())
		return
	}
}
