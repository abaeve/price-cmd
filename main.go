package main

import (
	"fmt"
	"github.com/abaeve/price-cmd/command"
	"github.com/abaeve/pricing-service/proto"
	"github.com/abaeve/sde-service/proto"
	proto "github.com/chremoas/chremoas/proto"
	"github.com/chremoas/services-common/args"
	"github.com/chremoas/services-common/config"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
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

	service = config.NewService(Version, "bot.command.cmd", name, initialize)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

// This function is a callback from the config.NewService function.  Read those docs
func initialize(config *config.Configuration) error {
	proto.RegisterCommandHandler(service.Server(),
		command.NewCommand(name, cmdArgs, &clientFactory{c: service.Client(), conf: config}),
	)

	return nil
}

type clientFactory struct {
	c    client.Client
	conf *config.Configuration
}

func (cf *clientFactory) PricingService() pricing.PricesService {
	return pricing.PricesServiceClient(cf.conf.LookupService("srv", "pricing"), cf.c)
}

func (cf *clientFactory) TypeService() sde.TypeQueryService {
	return sde.TypeQueryServiceClient(cf.conf.LookupService("srv", "sde"), cf.c)
}
