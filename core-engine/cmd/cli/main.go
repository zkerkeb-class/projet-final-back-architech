package main

import (
	"fmt"
	"log"

	"baymean/core/internal/domain"
	"baymean/core/internal/engine"
	"baymean/core/internal/vault"
	bolt "go.etcd.io/bbolt"
)

func main() {
	// 1. Ouvrir BoltDB en premier
	db, err := bolt.Open("bayment.db", 0600, nil)
	if err != nil {
		log.Fatal("Erreur ouverture BoltDB:", err)
	}
	defer db.Close()

	// 2. Créer le Vault persistant (identité stable entre sessions)
	sim, err := vault.NewPersistentVault(db, "1234")
	if err != nil {
		log.Fatal("Erreur création Vault:", err)
	}

	// 3. Démarrer le moteur
	e, err := engine.NewEngine(db, sim)
	if err != nil {
		log.Fatal("Erreur démarrage moteur:", err)
	}

	// 4. Créer un fragment genesis
	alice := string(sim.GetPublicKey())
	genesis := domain.Fragment{
		ParentIDs: []string{},
		Value:     50.0,
		OwnerPub:  alice,
	}
	genesis.ID = genesis.ComputeID()

	fmt.Println("Fragment genesis créé:", genesis.ID[:8], "- Valeur:", genesis.Value, "€")

	err = e.CreateGenesis(genesis)
	if err != nil {
		log.Fatal("Erreur sauvegarde genesis:", err)
	}

	// 5. Alice paie 30€ à Bob
	bobPub := "bob_public_key_mock"
	err = e.Split(genesis.ID, 30.0, bobPub)
	if err != nil {
		log.Fatal("Erreur split:", err)
	}
	fmt.Println("Split réussi : 30€ → Bob, 20€ → Alice (mergé)")

	// 6. Tenter un double spend — doit échouer
	err = e.Split(genesis.ID, 10.0, bobPub)
	if err != nil {
		fmt.Println("Double spend bloqué ✓:", err)
	}
}