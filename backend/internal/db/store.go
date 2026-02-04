package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"nimble-challenge/backend/internal/auth"
	"nimble-challenge/backend/internal/crypto"
)

func (s *Store) Authenticate(ctx context.Context, username, password string) (*auth.Principal, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("empty username")
	}

	var (
		userID     int64
		storeID    int64
		storeSlug  string
		passHash   string
		isMerchant bool
	)

	err := s.pool.QueryRow(ctx, `
		SELECT m.id, m.store_id, s.slug, m.password_hash
		FROM merchants m
		JOIN stores s ON s.id = m.store_id
		WHERE m.username = $1
	`, username).Scan(&userID, &storeID, &storeSlug, &passHash)
	if err == nil {
		ok, err := crypto.VerifyPassword(password, passHash)
		if err != nil || !ok {
			return nil, errors.New("invalid credentials")
		}
		isMerchant = true
	} else {
		err = s.pool.QueryRow(ctx, `
			SELECT c.id, c.store_id, s.slug, c.password_hash
			FROM customers c
			JOIN stores s ON s.id = c.store_id
			WHERE c.username = $1
		`, username).Scan(&userID, &storeID, &storeSlug, &passHash)
		if err != nil {
			return nil, errors.New("invalid credentials")
		}
		ok, err := crypto.VerifyPassword(password, passHash)
		if err != nil || !ok {
			return nil, errors.New("invalid credentials")
		}
	}

	role := auth.RoleCustomer
	if isMerchant {
		role = auth.RoleMerchant
	}

	return &auth.Principal{
		Role:      role,
		UserID:    userID,
		StoreID:   storeID,
		StoreSlug: storeSlug,
		Username:  username,
	}, nil
}

func (s *Store) CreatePet(ctx context.Context, storeID int64, input Pet) (Pet, error) {
	if input.Name == "" {
		return Pet{}, errors.New("name is required")
	}
	if input.AgeYears < 0 {
		return Pet{}, errors.New("age must be positive")
	}
	if input.PictureURL == "" {
		return Pet{}, errors.New("picture url is required")
	}
	if input.Description == "" {
		return Pet{}, errors.New("description is required")
	}
	if input.BreederName == "" {
		return Pet{}, errors.New("breeder name is required")
	}
	if input.Species != SpeciesCat && input.Species != SpeciesDog && input.Species != SpeciesFrog {
		return Pet{}, errors.New("invalid species")
	}
	if input.BreederEmail == "" {
		return Pet{}, errors.New("breeder email is required")
	}
	if !strings.Contains(input.BreederEmail, "@") {
		return Pet{}, errors.New("breeder email is invalid")
	}

	encEmail, nonce, err := s.crypto.Encrypt(input.BreederEmail)
	if err != nil {
		return Pet{}, fmt.Errorf("encrypt email: %w", err)
	}

	var petID string
	var createdAt time.Time
	err = s.pool.QueryRow(ctx, `
		INSERT INTO pets (
			store_id, name, species, age_years, picture_url, description,
			breeder_name, breeder_email_enc, breeder_email_nonce
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at
	`, storeID, input.Name, input.Species, input.AgeYears, input.PictureURL, input.Description,
		input.BreederName, encEmail, nonce).Scan(&petID, &createdAt)
	if err != nil {
		return Pet{}, fmt.Errorf("insert pet: %w", err)
	}

	input.ID = petID
	input.StoreID = storeID
	input.CreatedAt = createdAt
	return input, nil
}

func (s *Store) ListMerchantPets(ctx context.Context, storeID int64) ([]Pet, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, store_id, name, species, age_years, picture_url,
		       description, breeder_name, breeder_email_enc, breeder_email_nonce,
		       created_at, purchased_at
		FROM pets
		WHERE store_id = $1
		ORDER BY created_at DESC
	`, storeID)
	if err != nil {
		return nil, fmt.Errorf("query pets: %w", err)
	}
	defer rows.Close()

	var pets []Pet
	for rows.Next() {
		var pet Pet
		var emailEnc []byte
		var emailNonce []byte
		if err := rows.Scan(
			&pet.ID, &pet.StoreID, &pet.Name, &pet.Species, &pet.AgeYears,
			&pet.PictureURL, &pet.Description, &pet.BreederName,
			&emailEnc, &emailNonce, &pet.CreatedAt, &pet.PurchasedAt,
		); err != nil {
			return nil, fmt.Errorf("scan pet: %w", err)
		}
		email, err := s.crypto.Decrypt(emailEnc, emailNonce)
		if err != nil {
			return nil, fmt.Errorf("decrypt email: %w", err)
		}
		pet.BreederEmail = email
		pets = append(pets, pet)
	}
	return pets, nil
}

func (s *Store) ListAvailablePets(ctx context.Context, storeID int64) ([]Pet, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, store_id, name, species, age_years, picture_url,
		       description, breeder_name, breeder_email_enc, breeder_email_nonce,
		       created_at, purchased_at
		FROM pets
		WHERE store_id = $1 AND purchased_at IS NULL
		ORDER BY created_at DESC
	`, storeID)
	if err != nil {
		return nil, fmt.Errorf("query pets: %w", err)
	}
	defer rows.Close()

	var pets []Pet
	for rows.Next() {
		var pet Pet
		var emailEnc []byte
		var emailNonce []byte
		if err := rows.Scan(
			&pet.ID, &pet.StoreID, &pet.Name, &pet.Species, &pet.AgeYears,
			&pet.PictureURL, &pet.Description, &pet.BreederName,
			&emailEnc, &emailNonce, &pet.CreatedAt, &pet.PurchasedAt,
		); err != nil {
			return nil, fmt.Errorf("scan pet: %w", err)
		}
		email, err := s.crypto.Decrypt(emailEnc, emailNonce)
		if err != nil {
			return nil, fmt.Errorf("decrypt email: %w", err)
		}
		pet.BreederEmail = email
		pets = append(pets, pet)
	}
	return pets, nil
}

func (s *Store) ListPurchasedPets(ctx context.Context, storeID int64, customerID int64) ([]Pet, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, store_id, name, species, age_years, picture_url,
		       description, breeder_name, breeder_email_enc, breeder_email_nonce,
		       created_at, purchased_at
		FROM pets
		WHERE store_id = $1 AND purchased_by_customer_id = $2
		ORDER BY purchased_at DESC
	`, storeID, customerID)
	if err != nil {
		return nil, fmt.Errorf("query pets: %w", err)
	}
	defer rows.Close()

	var pets []Pet
	for rows.Next() {
		var pet Pet
		var emailEnc []byte
		var emailNonce []byte
		if err := rows.Scan(
			&pet.ID, &pet.StoreID, &pet.Name, &pet.Species, &pet.AgeYears,
			&pet.PictureURL, &pet.Description, &pet.BreederName,
			&emailEnc, &emailNonce, &pet.CreatedAt, &pet.PurchasedAt,
		); err != nil {
			return nil, fmt.Errorf("scan pet: %w", err)
		}
		email, err := s.crypto.Decrypt(emailEnc, emailNonce)
		if err != nil {
			return nil, fmt.Errorf("decrypt email: %w", err)
		}
		pet.BreederEmail = email
		pets = append(pets, pet)
	}
	return pets, nil
}

func (s *Store) PurchasePets(ctx context.Context, storeID int64, customerID int64, petIDs []string) (PurchaseResult, error) {
	if len(petIDs) == 0 {
		return PurchaseResult{}, errors.New("no pets in cart")
	}
	result := PurchaseResult{}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return result, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	rows, err := tx.Query(ctx, `
		SELECT id, name, purchased_at
		FROM pets
		WHERE store_id = $1 AND id = ANY($2)
		FOR UPDATE
	`, storeID, petIDs)
	if err != nil {
		return result, fmt.Errorf("select pets: %w", err)
	}
	defer rows.Close()

	available := make(map[string]string)
	seen := make(map[string]bool)
	for rows.Next() {
		var id string
		var name string
		var purchasedAt *time.Time
		if err := rows.Scan(&id, &name, &purchasedAt); err != nil {
			return result, fmt.Errorf("scan pet: %w", err)
		}
		seen[id] = true
		if purchasedAt != nil {
			result.Errors = append(result.Errors, PurchaseError{
				PetName: name,
				Message: "already purchased",
			})
			continue
		}
		available[id] = name
	}

	for _, id := range petIDs {
		if !seen[id] {
			result.Errors = append(result.Errors, PurchaseError{
				PetName: id,
				Message: "not found",
			})
		}
	}

	if len(available) > 0 {
		ids := make([]string, 0, len(available))
		for id := range available {
			ids = append(ids, id)
		}
		_, err = tx.Exec(ctx, `
			UPDATE pets
			SET purchased_at = NOW(), purchased_by_customer_id = $1
			WHERE store_id = $2 AND id = ANY($3) AND purchased_at IS NULL
		`, customerID, storeID, ids)
		if err != nil {
			return result, fmt.Errorf("update pets: %w", err)
		}
		result.PurchasedIDs = ids
	}

	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit: %w", err)
	}
	return result, nil
}
