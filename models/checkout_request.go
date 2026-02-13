package models

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}
