package main

import (
	"fmt"
	"log"

	"baymean/core/internal/domain"
	"baymean/core/internal/engine"
	"baymean/core/internal/vault"
)

func main() {
	// 1. Créer le Vault (Mock SIM)
	sim, err := vault.NewSIMMock()
	if err != nil {
		log.Fatal("Erreur création SIM:", err)
	}

	// 2. Démarrer le moteur avec BoltDB
	e, err := engine.NewEngine("bayment.db", sim)
	if err != nil {
		log.Fatal("Erreur démarrage moteur:", err)
	}
	defer e.DB.Close()

	// 3. Créer un fragment genesis (le premier billet, créé de nulle part)
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
	// 4. Simuler un paiement : Alice paie 30€ à Bob
	bobPub := "bob_public_key_mock"
	err = e.Split(genesis.ID, 30.0, bobPub)
	if err != nil {
		log.Fatal("Erreur split:", err)
	}

	fmt.Println("Split réussi : 30€ → Bob, 20€ → Alice (mergé)")

	// 5. Tenter un double spend — doit échouer
	err = e.Split(genesis.ID, 10.0, bobPub)
	if err != nil {
		fmt.Println("Double spend bloqué ✓:", err)
	}
}