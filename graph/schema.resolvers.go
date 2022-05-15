package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/translucent-link/owl/graph/generated"
	"github.com/translucent-link/owl/graph/model"
)

func (r *chainResolver) Protocols(ctx context.Context, obj *model.Chain) ([]*model.Protocol, error) {
	store, err := model.NewProtocolStore()
	if err != nil {
		return []*model.Protocol{}, err
	}
	return store.AllByChain(obj.ID)
}

func (r *mutationResolver) CreateChain(ctx context.Context, input model.NewChain) (*model.Chain, error) {
	store, err := model.NewChainStore()
	if err != nil {
		return &model.Chain{}, err
	}
	return store.CreateChain(input)
}

func (r *mutationResolver) CreateProtocol(ctx context.Context, input model.NewProtocol) (*model.Protocol, error) {
	store, err := model.NewProtocolStore()
	if err != nil {
		return &model.Protocol{}, err
	}
	return store.CreateProtocol(input)
}

func (r *mutationResolver) CreateProtocolInstance(ctx context.Context, input model.NewProtocolInstance) (*model.ProtocolInstance, error) {
	store, err := model.NewProtocolInstanceStore()
	if err != nil {
		return &model.ProtocolInstance{}, err
	}
	return store.CreateProtocolInstance(input)
}

func (r *mutationResolver) AddEventDefnToProtocol(ctx context.Context, input *model.NewEventDefn) (*model.EventDefn, error) {
	protocolStore, err := model.NewProtocolStore()
	if err != nil {
		return &model.EventDefn{}, err
	}
	topicSignature := []byte(input.AbiSignature)
	topicHash := crypto.Keccak256Hash(topicSignature)
	return protocolStore.AddEventDefn(input.Protocol, input.TopicName, topicHash.Hex(), input.AbiSignature)
}

func (r *protocolResolver) ScannableEvents(ctx context.Context, obj *model.Protocol) ([]*model.EventDefn, error) {
	protocolStore, err := model.NewProtocolStore()
	if err != nil {
		return []*model.EventDefn{}, err
	}
	return protocolStore.AllEventsByProtocol(obj.ID)
}

func (r *protocolInstanceResolver) Protocol(ctx context.Context, obj *model.ProtocolInstance) (*model.Protocol, error) {
	store, err := model.NewProtocolInstanceStore()
	if err != nil {
		return &model.Protocol{}, err
	}
	return store.FindProtocolById(obj.ID)
}

func (r *protocolInstanceResolver) Chain(ctx context.Context, obj *model.ProtocolInstance) (*model.Chain, error) {
	store, err := model.NewProtocolInstanceStore()
	if err != nil {
		return &model.Chain{}, err
	}
	return store.FindChainById(obj.ID)
}

func (r *queryResolver) Chains(ctx context.Context) ([]*model.Chain, error) {
	store, err := model.NewChainStore()
	if err != nil {
		return []*model.Chain{}, err
	}
	return store.All()
}

func (r *queryResolver) Protocols(ctx context.Context) ([]*model.Protocol, error) {
	store, err := model.NewProtocolStore()
	if err != nil {
		return []*model.Protocol{}, err
	}
	return store.All()
}

func (r *queryResolver) ProtocolInstances(ctx context.Context) ([]*model.ProtocolInstance, error) {
	store, err := model.NewProtocolInstanceStore()
	if err != nil {
		return []*model.ProtocolInstance{}, err
	}
	return store.All()
}

// Chain returns generated.ChainResolver implementation.
func (r *Resolver) Chain() generated.ChainResolver { return &chainResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Protocol returns generated.ProtocolResolver implementation.
func (r *Resolver) Protocol() generated.ProtocolResolver { return &protocolResolver{r} }

// ProtocolInstance returns generated.ProtocolInstanceResolver implementation.
func (r *Resolver) ProtocolInstance() generated.ProtocolInstanceResolver {
	return &protocolInstanceResolver{r}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type chainResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type protocolResolver struct{ *Resolver }
type protocolInstanceResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func Db() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
}
