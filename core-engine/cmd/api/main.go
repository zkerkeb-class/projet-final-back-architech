package main

import (
    "log"
    "net/http"

    "baymean/core/internal/api"
    "baymean/core/internal/engine"
    "baymean/core/internal/vault"
    bolt "go.etcd.io/bbolt"
)

func main() {
    //TODO : Create DB for Phone OS (probably SQLite) and function to retrieve money data from it (see for Token connection as well when offline)
    //TODO : Create function to save money data (and maybe Token) on Phone OS
    db, err := bolt.Open("bayment.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    sim, err := vault.NewPersistentVault(db, "1234")
    if err != nil {
        log.Fatal(err)
    }

    e, err := engine.NewEngine(db, sim)
    if err != nil {
        log.Fatal(err)
    }

    h := &api.Handler{Engine: e}

    http.HandleFunc("/genesis", h.Genesis)
    http.HandleFunc("/split", h.Split)
    http.HandleFunc("/balance", h.Balance)

    log.Println("Moteur BAYMEAN démarré sur :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}