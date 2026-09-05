package main

import "testing"

func TestDecideUpload(t *testing.T) {
	tests := []struct {
		name        string
		event       paymentEvent
		wantStatus  string
		wantAllowed bool
	}{
		{"cleared payment", paymentEvent{PaymentID: "pay_42", AssetName: "invoice.pdf", Risk: "clear"}, "presign issued", true},
		{"blocked payment", paymentEvent{PaymentID: "pay_43", AssetName: "invoice.pdf", Risk: "blocked"}, "manual review", false},
		{"missing asset", paymentEvent{PaymentID: "pay_44", Risk: "clear"}, "manual review", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, allowed := decideUpload(tt.event)
			if status != tt.wantStatus || allowed != tt.wantAllowed {
				t.Fatalf("got (%q, %v), want (%q, %v)", status, allowed, tt.wantStatus, tt.wantAllowed)
			}
		})
	}
}
