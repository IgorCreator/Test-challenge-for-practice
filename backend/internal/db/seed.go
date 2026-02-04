package db

import (
	"context"
	"fmt"

	"nimble-challenge/backend/internal/crypto"
)

func (s *Store) EnsureDemoData(ctx context.Context, storeSlug, storeName, merchantUser, merchantPass, customerUser, customerPass string) error {
	var storeID int64
	err := s.pool.QueryRow(ctx, `
		INSERT INTO stores (slug, name)
		VALUES ($1, $2)
		ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, storeSlug, storeName).Scan(&storeID)
	if err != nil {
		return fmt.Errorf("upsert store: %w", err)
	}

	if err := ensureUser(ctx, s, storeID, "merchants", merchantUser, merchantPass); err != nil {
		return err
	}
	if err := ensureUser(ctx, s, storeID, "customers", customerUser, customerPass); err != nil {
		return err
	}

	var petCount int
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(1) FROM pets WHERE store_id = $1`, storeID).Scan(&petCount); err != nil {
		return fmt.Errorf("count pets: %w", err)
	}
	if petCount == 0 {
		seedPets := []Pet{
			{
				Name:         "Miso",
				Species:      SpeciesCat,
				AgeYears:     2,
				PictureURL:   "https://images.unsplash.com/photo-1518791841217-8f162f1e1131?auto=format&fit=crop&w=900&q=80",
				Description:  "Playful kitten who loves strings and sunbeams.",
				BreederName:  "Jane Doe",
				BreederEmail: "jane@example.com",
			},
			{
				Name:         "Barkley",
				Species:      SpeciesDog,
				AgeYears:     4,
				PictureURL:   "https://images.unsplash.com/photo-1507146426996-ef05306b995a?auto=format&fit=crop&w=900&q=80",
				Description:  "Friendly golden retriever who enjoys long walks.",
				BreederName:  "Tom Rivers",
				BreederEmail: "tom@example.com",
			},
			{
				Name:         "Sprout",
				Species:      SpeciesFrog,
				AgeYears:     1,
				PictureURL:   "https://images.unsplash.com/photo-1502786129293-79981df4e689?auto=format&fit=crop&w=900&q=80",
				Description:  "Tiny tree frog with a calm personality.",
				BreederName:  "Lena Moss",
				BreederEmail: "lena@example.com",
			},
		}
		for _, pet := range seedPets {
			if _, err := s.CreatePet(ctx, storeID, pet); err != nil {
				return fmt.Errorf("seed pet: %w", err)
			}
		}
	}
	return nil
}

func ensureUser(ctx context.Context, s *Store, storeID int64, table, username, password string) error {
	var count int
	err := s.pool.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(1) FROM %s WHERE username = $1`, table), username).Scan(&count)
	if err != nil {
		return fmt.Errorf("count user: %w", err)
	}
	if count > 0 {
		return nil
	}
	hash, err := crypto.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	_, err = s.pool.Exec(ctx, fmt.Sprintf(`
		INSERT INTO %s (store_id, username, password_hash)
		VALUES ($1, $2, $3)
	`, table), storeID, username, hash)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}
