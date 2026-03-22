package api

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"baymean/core/internal/domain"
	"baymean/core/internal/engine"
)

type Handler struct {
	Engine *engine.Engine
}

// POST /genesis
// Body: { "amount": 50.0 }
func (h *Handler) Genesis(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "body invalide", http.StatusBadRequest)
		return
	}

	ownerPub := hex.EncodeToString(h.Engine.Vault.GetPublicKey())

	genesis := domain.Fragment{
		ParentIDs: []string{},
		Value:     body.Amount,
		OwnerPub:  ownerPub,
		CreatedAt: time.Now(),
	}
	genesis.ID = genesis.ComputeID()

	if err := h.Engine.CreateGenesis(genesis); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(genesis)
}

// POST /split
// Body: { "parent_id": "...", "amount": 30.0, "recipient_pub": "..." }
func (h *Handler) Split(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ParentID     string  `json:"parent_id"`
		Amount       float64 `json:"amount"`
		RecipientPub string  `json:"recipient_pub"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "body invalide", http.StatusBadRequest)
		return
	}

	if err := h.Engine.Split(body.ParentID, body.Amount, body.RecipientPub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET /balance
func (h *Handler) Balance(w http.ResponseWriter, r *http.Request) {
	ownerPub := hex.EncodeToString(h.Engine.Vault.GetPublicKey())
	balance := h.Engine.GetBalance(ownerPub)

	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}