package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/lib/pq"
	"github.com/translucent-link/owl/graph/generated"
	"github.com/translucent-link/owl/graph/model"
	"github.com/translucent-link/owl/index"
)

func (r *accountResolver) Events(ctx context.Context, obj *model.Account) ([]model.AnyEvent, error) {
	events := []model.AnyEvent{}
	db, err := model.DbConnect()
	if err != nil {
		return events, errors.Wrap(err, "Unable to connect to DB and retrieve events")
	}
	defer db.Close()
	stores := model.GenerateStores(db)
	aEvents, err := stores.Event.AllByAccount(obj.ID)
	if err != nil {
		return events, err
	}
	for _, aEvent := range aEvents {
		var event model.AnyEvent
		event, err = aEvent.AnyEvent(stores.Account, stores.Token)
		if err != nil {
			return events, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *chainResolver) Protocols(ctx context.Context, obj *model.Chain) ([]*model.Protocol, error) {
	db, err := model.DbConnect()
	if err != nil {
		return []*model.Protocol{}, errors.Wrap(err, "Unable to connect to DB and retrieve protocols")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	return stores.Protocol.AllByChain(obj.ID)
}

func (r *chainResolver) Tokens(ctx context.Context, obj *model.Chain) ([]*model.Token, error) {
	db, err := model.DbConnect()
	if err != nil {
		return []*model.Token{}, errors.Wrap(err, "Unable to connect to DB and retrieve tokens")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	return stores.Token.AllByChain(obj.ID)
}

func (r *mutationResolver) CreateChain(ctx context.Context, input model.NewChain) (*model.Chain, error) {
	db, err := model.DbConnect()
	if err != nil {
		return &model.Chain{}, errors.Wrap(err, "Unable to connect to DB and create chain")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	return stores.Chain.CreateChain(input)
}

func (r *mutationResolver) CreateProtocol(ctx context.Context, input model.NewProtocol) (*model.Protocol, error) {
	db, err := model.DbConnect()
	if err != nil {
		return &model.Protocol{}, errors.Wrap(err, "Unable to connect to DB and create protocol")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	return stores.Protocol.CreateProtocol(input)
}

func (r *mutationResolver) CreateProtocolInstance(ctx context.Context, input model.NewProtocolInstance) (*model.ProtocolInstance, error) {
	db, err := model.DbConnect()
	if err != nil {
		return &model.ProtocolInstance{}, errors.Wrap(err, "Unable to connect to DB and create protocol instance")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	return stores.ProtocolInstance.CreateProtocolInstance(input)
}

func (r *mutationResolver) AddEventDefnToProtocol(ctx context.Context, input *model.NewEventDefn) (*model.EventDefn, error) {
	db, err := model.DbConnect()
	if err != nil {
		return &model.EventDefn{}, errors.Wrap(err, "Unable to connect to DB and create event definition")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	topicSignature := []byte(input.AbiSignature)
	topicHash := crypto.Keccak256Hash(topicSignature)
	return stores.Protocol.AddEventDefn(input.Protocol, input.TopicName, topicHash.Hex(), input.AbiSignature)
}

func (r *mutationResolver) ScanProtocolInstance(ctx context.Context, input model.NewScan) (*model.ProtocolInstance, error) {
	db, err := model.DbConnect()
	if err != nil {
		return &model.ProtocolInstance{}, errors.Wrap(err, "Unable to connect to DB in preparation of scanning")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	protocol, err := stores.Protocol.FindByName(input.Protocol)
	chain, err := stores.Chain.FindByName(input.Chain)

	protocolInstance, err := stores.ProtocolInstance.FindByProtocolIdAndChainId(protocol.ID, chain.ID)
	if err != nil {
		return &model.ProtocolInstance{}, errors.Wrap(err, "Unable to connect to DB and fetch protocol instance")
	}
	scannableEvents, err := stores.Protocol.AllEventsByProtocol(protocol.ID)
	if err != nil {
		return &model.ProtocolInstance{}, errors.Wrap(err, "Retrieving list of scannable events")

	}

	client, err := index.GetClient(chain.RPCURL)
	if err != nil {
		return &model.ProtocolInstance{}, errors.Wrap(err, "Retrieving EVM client")
	}

	log.Printf("Scanning %s on %s", protocol.Name, chain.Name)

	scanRequest := index.ScanRequest{
		Client:           client,
		Chain:            chain,
		Protocol:         protocol,
		ProtocolInstance: protocolInstance,
		ScannableEvents:  scannableEvents,
	}
	index.ScanChannel <- scanRequest
	log.Println("Scan Requested")

	return protocolInstance, nil
}

func (r *protocolResolver) ScannableEvents(ctx context.Context, obj *model.Protocol) ([]*model.EventDefn, error) {
	db, err := model.DbConnect()
	if err != nil {
		return []*model.EventDefn{}, errors.Wrap(err, "Unable to connect to DB and retrieve scannable events")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	return stores.Protocol.AllEventsByProtocol(obj.ID)
}

func (r *protocolInstanceResolver) Protocol(ctx context.Context, obj *model.ProtocolInstance) (*model.Protocol, error) {
	db, err := model.DbConnect()
	if err != nil {
		return &model.Protocol{}, errors.Wrap(err, "Unable to connect to DB and retrieve scannable events")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	return stores.ProtocolInstance.FindProtocolById(obj.ID)
}

func (r *protocolInstanceResolver) Chain(ctx context.Context, obj *model.ProtocolInstance) (*model.Chain, error) {
	db, err := model.DbConnect()
	if err != nil {
		return &model.Chain{}, errors.Wrap(err, "Unable to connect to DB and retrieve scannable events")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	return stores.ProtocolInstance.FindChainById(obj.ID)
}

func (r *queryResolver) Chains(ctx context.Context) ([]*model.Chain, error) {
	db, err := model.DbConnect()
	if err != nil {
		return []*model.Chain{}, errors.Wrap(err, "Unable to connect to DB and retrieve chains")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	return stores.Chain.All()
}

func (r *queryResolver) Protocols(ctx context.Context) ([]*model.Protocol, error) {
	db, err := model.DbConnect()
	if err != nil {
		return []*model.Protocol{}, errors.Wrap(err, "Unable to connect to DB and retrieve chains")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	return stores.Protocol.All()
}

func (r *queryResolver) ProtocolInstances(ctx context.Context) ([]*model.ProtocolInstance, error) {
	db, err := model.DbConnect()
	if err != nil {
		return []*model.ProtocolInstance{}, errors.Wrap(err, "Unable to connect to DB and retrieve protocol instances")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	return stores.ProtocolInstance.All()
}

func (r *queryResolver) Accounts(ctx context.Context, address *string) ([]*model.Account, error) {
	db, err := model.DbConnect()
	if err != nil {
		return []*model.Account{}, errors.Wrap(err, "Unable to connect to DB and retrieve accounts")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	if address != nil {
		acc, err := stores.Account.FindByAddress(*address)
		return []*model.Account{acc}, err
	}
	return []*model.Account{}, err
}

func (r *queryResolver) Borrowers(ctx context.Context, top *int) ([]*model.Account, error) {
	return []*model.Account{}, fmt.Errorf("not implemented")
}

func (r *queryResolver) Liquidators(ctx context.Context, top *int) ([]*model.Account, error) {
	return []*model.Account{}, fmt.Errorf("not implemented")
}

// Account returns generated.AccountResolver implementation.
func (r *Resolver) Account() generated.AccountResolver { return &accountResolver{r} }

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

type accountResolver struct{ *Resolver }
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
