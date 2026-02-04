package db

import "time"

type Species string

const (
	SpeciesCat  Species = "CAT"
	SpeciesDog  Species = "DOG"
	SpeciesFrog Species = "FROG"
)

type Pet struct {
	ID           string
	StoreID      int64
	Name         string
	Species      Species
	AgeYears     int
	PictureURL   string
	Description  string
	BreederName  string
	BreederEmail string
	CreatedAt    time.Time
	PurchasedAt  *time.Time
}

type PurchaseError struct {
	PetName string
	Message string
}

type PurchaseResult struct {
	PurchasedIDs []string
	Errors       []PurchaseError
}
