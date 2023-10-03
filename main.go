package main

import (
    "database/sql"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/lib/pq"
)

var (
    db          *sql.DB
    portNumber  int
    postgresURL string
)

func main() {
    flag.IntVar(&portNumber, "port", 8080, "Port number for the server")
    flag.StringVar(&postgresURL, "postgres-url", "user=your_username dbname=your_database sslmode=disable", "PostgreSQL connection URL")
    flag.Parse()

    // Initialize the database
    if err := initDatabase(); err != nil {
        log.Fatalf("Failed to initialize the database: %v", err)
    }

    // Create the HTTP server
    srv := &http.Server{
        Addr:    fmt.Sprintf(":%d", portNumber),
        Handler: createHandler(),
    }

    // Start the HTTP server
    go func() {
        log.Printf("Server listening on port %d\n", portNumber)
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Handle shutdown gracefully
    waitForShutdown(srv)
}

func initDatabase() error {
    var err error
    db, err = sql.Open("postgres", postgresURL)
    if err != nil {
        return err
    }

    err = db.Ping()
    if err != nil {
        return err
    }

    // Create the 'aliases' table if it doesn't exist
    createTable()

    return nil
}

func createHandler() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/aliases", handleAliases)
    return mux
}

func handleAliases(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        aliases, err := getAliases()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        jsonBytes, err := json.Marshal(aliases)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(jsonBytes)

    case http.MethodPost:
        var newAlias Alias
        err := json.NewDecoder(r.Body).Decode(&newAlias)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        err = saveAlias(&newAlias)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(fmt.Sprintf("Alias for \"%s\" has been saved!\n", newAlias.Name)))
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func createTable() {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS aliases (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            text TEXT
        )
    `)
    if err != nil {
        log.Fatal(err)
    }
}

func getAliases() ([]Alias, error) {
    rows, err := db.Query("SELECT id, name, text FROM aliases")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var aliases []Alias
    for rows.Next() {
        var alias Alias
        err := rows.Scan(&alias.ID, &alias.Name, &alias.Text)
        if err != nil {
            return nil, err
        }
        aliases = append(aliases, alias)
    }

    return aliases, nil
}

func saveAlias(alias *Alias) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    var existingID int
    err = tx.QueryRow("SELECT id FROM aliases WHERE name = $1", alias.Name).Scan(&existingID)
    if err != nil && err != sql.ErrNoRows {
        return err
    }

    if existingID > 0 {
        _, err = tx.Exec("UPDATE aliases SET text = $2 WHERE id = $1", existingID, alias.Text)
        if err != nil {
            return err
        }
    } else {
        _, err = tx.Exec("INSERT INTO aliases (name, text) VALUES ($1, $2)", alias.Name, alias.Text)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func waitForShutdown(srv *http.Server) {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    <-sigChan

    log.Println("Shutting down...")

    // Create a context with a timeout to force shutdown after a period
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server shutdown error: %v", err)
    }

    log.Println("Server gracefully stopped")
}
