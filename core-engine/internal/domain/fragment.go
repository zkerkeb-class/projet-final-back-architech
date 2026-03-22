package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

// Fragment est l'unité monétaire atomique du système.
// Contrairement à un solde (account_money: Number), un Fragment est :
// - Immuable une fois créé
// - Cryptographiquement lié à son créateur (Signature)
// - Traçable jusqu'à son origine (ParentID = généalogie)
type Fragment struct {
	ID        string    `json:"id"`         // ID unique du fragment 
	ParentIDs []string  `json:"parent_ids"`  // ID des fragments parent , par exemple liste vide pour le premier fragment
	Value     float64   `json:"value"`      // Valeur monétaire du fragment
	OwnerPub  string    `json:"owner_pub"`  // Clé publique du propriétaire
	CreatedAt time.Time `json:"created_at"` // Timestamp de création
	Signature []byte    `json:"signature"`   // Signature cryptographique du créateur
	IsSpent   bool      `json:"is_spent"`   // Indique si le fragment a été dépensé
}

// ComputeID génère un ID déterministe à partir du contenu du fragment.
// Déterministe = le même fragment produit TOUJOURS le même ID,
// peu importe l'OS ou l'architecture. C'est la clé de la sérialisation.
func (f *Fragment) ComputeID() string {
	payload := struct {
		ParentIDs []string  `json:"parent_ids"`
		Value     float64   `json:"value"`
		OwnerPub  string    `json:"owner_pub"`
		CreatedAt time.Time `json:"created_at"`
	}{f.ParentIDs, f.Value, f.OwnerPub, f.CreatedAt}

	data, _ := json.Marshal(payload)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}