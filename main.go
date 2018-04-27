package main

import (
	"fmt"
	"github.com/chremoas/services-common/config"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"

	"context"
	"github.com/abaeve/price-cmd/command"
	"github.com/abaeve/pricing-service/proto"
	"github.com/abaeve/sde-service/proto"
	proto "github.com/chremoas/chremoas/proto"
	"github.com/chremoas/services-common/args"
)

var Version = "1.0.0"
var service micro.Service
var name = "price"

var cmdArgs *args.Args

func main() {
	cmdArgs = args.NewArg("price")

	//Makes the command's arguments available
	//TODO: Maybe this get's moved to the command/handler.go file?
	for an, ac := range command.Commands {
		cmdArgs.Add(an, ac)
	}

	service = config.NewService(Version, "cmd", name, initialize)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

// This function is a callback from the config.NewService function.  Read those docs
func initialize(config *config.Configuration) error {
	proto.RegisterCommandHandler(service.Server(),
		command.NewCommand(name, cmdArgs, &clientFactory{service.Client()}),
	)

	return nil
}

type clientFactory struct {
	c client.Client
}

func (cf *clientFactory) PricingService() pricing.PricesService {
	return &pricingService{}
}

func (cf *clientFactory) TypeService() sde.TypeQueryService {
	return &typeService{}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type pricingService struct {
}

func (ps *pricingService) GetItemPrice(ctx context.Context, in *pricing.ItemPriceRequest, opts ...client.CallOption) (*pricing.ItemPriceResponse, error) {
	return &pricing.ItemPriceResponse{
		Item: &pricing.Item{
			ItemId: in.ItemId,
			Buy: &pricing.ItemPrice{
				Min: 0.1,
				Max: 0.2,
				Avg: 0.3,
				Vol: 0.4,
				Ord: 1,
			},
			Sell: &pricing.ItemPrice{
				Min: 0.5,
				Max: 0.6,
				Avg: 0.7,
				Vol: 0.8,
				Ord: 2,
			},
		},
	}, nil
}

type typeService struct {
}

func (ts *typeService) FindTypesByTypeIds(ctx context.Context, in *sde.TypeIdRequest, opts ...client.CallOption) (*sde.TypeResponse, error) {
	return &sde.TypeResponse{}, nil
}

func (ts *typeService) FindTypesByTypeNames(ctx context.Context, in *sde.TypeNameRequest, opts ...client.CallOption) (*sde.TypeResponse, error) {
	return &sde.TypeResponse{
		Type: []*sde.Type{
			{
				TypeId: 1234,
				Name:   "150mm Rail Gun II",
			},
		},
	}, nil
}

func (ts *typeService) SearchForTypes(ctx context.Context, in *sde.TypeNameRequest, opts ...client.CallOption) (*sde.TypeNameAndIdResponse, error) {
	return &sde.TypeNameAndIdResponse{
		TypeId: []int32{1234},
		Name:   []string{"150mm Rail Gun II"},
	}, nil
}
