package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const bucket = "fintech-asset-intake"

type paymentEvent struct {
	PaymentID string `json:"payment_id"`
	AssetName string `json:"asset_name"`
	Risk      string `json:"risk"`
}

type uploadReply struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	UploadURL string `json:"upload_url"`
	AuditNote string `json:"audit_note"`
}

func decideUpload(e paymentEvent) (string, bool) {
	if e.PaymentID == "" || e.AssetName == "" || strings.EqualFold(e.Risk, "blocked") {
		return "manual review", false
	}
	return "presign issued", true
}

func main() {
	client, err := NewInfraiClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/payment-events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var event paymentEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		status, allowed := decideUpload(event)
		reply := uploadReply{PaymentID: event.PaymentID, Status: status, AuditNote: "payment event evaluated"}
		if allowed {
			digest := sha256.Sum256([]byte(event.PaymentID + ":" + event.AssetName))
			key := event.PaymentID + "/" + hex.EncodeToString(digest[:8]) + "-" + event.AssetName
			presign, err := client.PresignPut(bucket, key, event.PaymentID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
			reply.UploadURL = presign.URL
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reply)
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
