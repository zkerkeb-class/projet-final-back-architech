package engine

import (
	"encoding/json"
	"errors"
	"time"

	"baymean/core/internal/domain"
	"baymean/core/internal/vault"
	bolt "go.etcd.io/bbolt"
)

var bucketName = []byte("Fragments")

type Engine struct {
	DB    *bolt.DB
	Vault vault.SecureVault
}

func (e *Engine) Split(parentID string, amount float64, recipientPub string) error {
	return e.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		// 1. Lire le fragment parent
		data := b.Get([]byte(parentID))
		if data == nil {
			return errors.New("fragment introuvable")
		}

		var parent domain.Fragment
		if err := json.Unmarshal(data, &parent); err != nil {
			return err
		}

		// 2. Vérifications
		if parent.IsSpent {
			return errors.New("double spend détecté")
		}
		if amount <= 0 || amount >= parent.Value {
			return errors.New("montant invalide")
		}

		// 3. Créer les deux fragments enfants
		now := time.Now()
		change := parent.Value - amount

		recipient := domain.Fragment{
			ParentIDs: []string{parent.ID},
			Value:     amount,
			OwnerPub:  recipientPub,
			CreatedAt: now,
		}
		recipient.ID = recipient.ComputeID()

		// 4. Merge automatique : le change fusionne avec les fragments
		// existants non dépensés du propriétaire
		existingIDs, existingTotal := e.findUnspentFragments(b, parent.OwnerPub, parentID)

		parentIDs := append([]string{parent.ID}, existingIDs...)
		merged := domain.Fragment{
			ParentIDs: parentIDs,
			Value:     change + existingTotal,
			OwnerPub:  parent.OwnerPub,
			CreatedAt: now,
		}
		merged.ID = merged.ComputeID()

		// 5. Signer les deux fragments avec le Vault (clé privée reste dans la SIM)
		sig, err := e.Vault.Sign([]byte(recipient.ID))
		if err != nil {
			return err
		}
		recipient.Signature = sig

		sig, err = e.Vault.Sign([]byte(merged.ID))
		if err != nil {
			return err
		}
		merged.Signature = sig

		// 6. Écriture atomique — tout ou rien
		parent.IsSpent = true
		if err := put(b, parent); err != nil {
			return err
		}
		// Marquer les fragments fusionnés comme dépensés
		for _, id := range existingIDs {
			if err := markSpent(b, id); err != nil {
				return err
			}
		}
		if err := put(b, recipient); err != nil {
			return err
		}
		return put(b, merged)
	})
}

// findUnspentFragments retourne les IDs et la valeur totale
// des fragments non dépensés d'un propriétaire (hors le parent courant)
func (e *Engine) findUnspentFragments(b *bolt.Bucket, ownerPub, excludeID string) ([]string, float64) {
	var ids []string
	var total float64

	b.ForEach(func(k, v []byte) error {
		var f domain.Fragment
		if err := json.Unmarshal(v, &f); err != nil {
			return nil
		}
		if f.OwnerPub == ownerPub && !f.IsSpent && f.ID != excludeID {
			ids = append(ids, f.ID)
			total += f.Value
		}
		return nil
	})

	return ids, total
}

func put(b *bolt.Bucket, f domain.Fragment) error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	return b.Put([]byte(f.ID), data)
}

func markSpent(b *bolt.Bucket, id string) error {
	data := b.Get([]byte(id))
	if data == nil {
		return errors.New("fragment introuvable pour markSpent")
	}
	var f domain.Fragment
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	f.IsSpent = true
	return put(b, f)
}

func (e *Engine) CreateGenesis(f domain.Fragment) error {
	return e.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return put(b, f)
	})
}

func NewEngine(db *bolt.DB, v vault.SecureVault) (*Engine, error) {
    err := db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists(bucketName)
        return err
    })
    return &Engine{DB: db, Vault: v}, err
}

func (e *Engine) GetBalance(ownerPub string) float64 {
    var total float64
    e.DB.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(bucketName)
        _, total = e.findUnspentFragments(b, ownerPub, "")
        return nil
    })
    return total
}