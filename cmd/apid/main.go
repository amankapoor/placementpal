// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// This program provides a sample web service that implements a
// RESTFul CRUD API against a MongoDB database.
package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/amankapoor/goth/gothic"
	"github.com/amankapoor/placementpal/cmd/apid/handlers"
	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/amankapoor/placementpal/internal/platform/web"
	"github.com/gorilla/sessions"
)

// init is called before main. We are using init to customize logging output.
func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	web.ParseTemplates("templates", funcs)

	// goth settings
	key := "goth-example" // Replace with your SESSION_SECRET or similar
	maxAge := 86400       // 1 day
	isProd := false       // Set to true when serving over https
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd
	//store := sessions.NewFilesystemStore(os.TempDir(), []byte("goth-example"))
	//store.MaxLength(math.MaxInt64)
	gothic.Store = store
}

var funcs = template.FuncMap{"add": add}

func add(x, y int) int {
	return x + y
}

// main is the entry point for the application.
func main() {

	// Initialising logging to file
	f, err := os.OpenFile("logs/logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Printf("\n\n\nmain : Started on %v", time.Now())

	// Check the environment for a configured port value.
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:27017"
		//dbHost = "mongodb://aman:password@ds119685.mlab.com:19685/pp"
		//dbHost = "got:got2015@ds039441.mongolab.com:39441/gotraining"
	}

	// Register the Master Session for the database.
	log.Println("main : Started : Capturing Master DB...")
	masterDB, err := db.NewMGO(dbHost, 25*time.Second)
	if err != nil {
		log.Fatalf("startup : Register DB : %v", err)
	}
	defer masterDB.MGOClose()

	// Check the environment for a configured port value.
	host := os.Getenv("HOST")
	if host == "" {
		host = ":8080"
	}

	// Create a new server and set timeout values.
	server := http.Server{
		Addr:           host,
		Handler:        handlers.API(masterDB),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// We want to report the listener is closed.
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the listener.
	go func() {
		log.Printf("startup : Listening %s", host)
		log.Printf("shutdown : Listener closed : %v", server.ListenAndServe())
		wg.Done()
	}()

	// Listen for an interrupt signal from the OS.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)

	// Wait for a signal to shutdown.
	<-osSignals

	// Create a context to attempt a graceful 5 second shutdown.
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Attempt the graceful shutdown by closing the listener and
	// completing all inflight requests.
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown : Graceful shutdown did not complete in %v : %v", timeout, err)

		// Looks like we timedout on the graceful shutdown. Kill it hard.
		if err := server.Close(); err != nil {
			log.Printf("shutdown : Error killing server : %v", err)
		}
	}

	// Wait for the listener to report it is closed.
	wg.Wait()
	log.Println("main : Completed")
}
