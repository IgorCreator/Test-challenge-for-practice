package graphql

import (
	"context"
	"errors"

	gql "github.com/graph-gophers/graphql-go"

	"nimble-challenge/backend/internal/auth"
	"nimble-challenge/backend/internal/db"
)

type Resolver struct {
	Store *db.Store
}

type CreatePetInput struct {
	Name         string
	Species      db.Species
	AgeYears     int32
	PictureURL   string
	Description  string
	BreederName  string
	BreederEmail string
}

type PurchasePetsInput struct {
	StoreSlug string
	PetIDs    []gql.ID
}

func (r *Resolver) MerchantPets(ctx context.Context) ([]*PetResolver, error) {
	principal, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if principal.Role != auth.RoleMerchant {
		return nil, errors.New("merchant access required")
	}
	pets, err := r.Store.ListMerchantPets(ctx, principal.StoreID)
	if err != nil {
		return nil, err
	}
	return wrapPets(pets), nil
}

func (r *Resolver) StorePets(ctx context.Context, args struct{ StoreSlug string }) ([]*PetResolver, error) {
	principal, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if principal.Role != auth.RoleCustomer {
		return nil, errors.New("customer access required")
	}
	if principal.StoreSlug != args.StoreSlug {
		return nil, errors.New("store access denied")
	}
	pets, err := r.Store.ListAvailablePets(ctx, principal.StoreID)
	if err != nil {
		return nil, err
	}
	return wrapPets(pets), nil
}

func (r *Resolver) PurchasedPets(ctx context.Context, args struct{ StoreSlug string }) ([]*PetResolver, error) {
	principal, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if principal.Role != auth.RoleCustomer {
		return nil, errors.New("customer access required")
	}
	if principal.StoreSlug != args.StoreSlug {
		return nil, errors.New("store access denied")
	}
	pets, err := r.Store.ListPurchasedPets(ctx, principal.StoreID, principal.UserID)
	if err != nil {
		return nil, err
	}
	return wrapPets(pets), nil
}

func (r *Resolver) CreatePet(ctx context.Context, args struct{ Input CreatePetInput }) (*PetResolver, error) {
	principal, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if principal.Role != auth.RoleMerchant {
		return nil, errors.New("merchant access required")
	}
	pet, err := r.Store.CreatePet(ctx, principal.StoreID, db.Pet{
		Name:         args.Input.Name,
		Species:      args.Input.Species,
		AgeYears:     int(args.Input.AgeYears),
		PictureURL:   args.Input.PictureURL,
		Description:  args.Input.Description,
		BreederName:  args.Input.BreederName,
		BreederEmail: args.Input.BreederEmail,
	})
	if err != nil {
		return nil, err
	}
	return &PetResolver{pet: pet}, nil
}

func (r *Resolver) PurchasePets(ctx context.Context, args struct{ Input PurchasePetsInput }) (*PurchaseResultResolver, error) {
	principal, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if principal.Role != auth.RoleCustomer {
		return nil, errors.New("customer access required")
	}
	if principal.StoreSlug != args.Input.StoreSlug {
		return nil, errors.New("store access denied")
	}

	ids := make([]string, 0, len(args.Input.PetIDs))
	for _, id := range args.Input.PetIDs {
		ids = append(ids, string(id))
	}

	result, err := r.Store.PurchasePets(ctx, principal.StoreID, principal.UserID, ids)
	if err != nil {
		return nil, err
	}
	return &PurchaseResultResolver{result: result}, nil
}

type PetResolver struct {
	pet db.Pet
}

func wrapPets(pets []db.Pet) []*PetResolver {
	resolvers := make([]*PetResolver, 0, len(pets))
	for _, pet := range pets {
		resolvers = append(resolvers, &PetResolver{pet: pet})
	}
	return resolvers
}

func (p *PetResolver) ID() gql.ID           { return gql.ID(p.pet.ID) }
func (p *PetResolver) Name() string         { return p.pet.Name }
func (p *PetResolver) Species() db.Species  { return p.pet.Species }
func (p *PetResolver) AgeYears() int32      { return int32(p.pet.AgeYears) }
func (p *PetResolver) PictureUrl() string   { return p.pet.PictureURL }
func (p *PetResolver) Description() string  { return p.pet.Description }
func (p *PetResolver) BreederName() string  { return p.pet.BreederName }
func (p *PetResolver) BreederEmail() string { return p.pet.BreederEmail }
func (p *PetResolver) CreatedAt() gql.Time  { return gql.Time{Time: p.pet.CreatedAt} }
func (p *PetResolver) PurchasedAt() *gql.Time {
	if p.pet.PurchasedAt == nil {
		return nil
	}
	t := gql.Time{Time: *p.pet.PurchasedAt}
	return &t
}

type PurchaseErrorResolver struct {
	err db.PurchaseError
}

func (p *PurchaseErrorResolver) PetName() string { return p.err.PetName }
func (p *PurchaseErrorResolver) Message() string { return p.err.Message }

type PurchaseResultResolver struct {
	result db.PurchaseResult
}

func (r *PurchaseResultResolver) PurchasedIds() []gql.ID {
	ids := make([]gql.ID, 0, len(r.result.PurchasedIDs))
	for _, id := range r.result.PurchasedIDs {
		ids = append(ids, gql.ID(id))
	}
	return ids
}

func (r *PurchaseResultResolver) Errors() []*PurchaseErrorResolver {
	errs := make([]*PurchaseErrorResolver, 0, len(r.result.Errors))
	for _, err := range r.result.Errors {
		errs = append(errs, &PurchaseErrorResolver{err: err})
	}
	return errs
}
