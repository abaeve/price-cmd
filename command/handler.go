package command

import (
	"github.com/abaeve/pricing-service/proto"
	"github.com/abaeve/sde-service/proto"
	bot "github.com/chremoas/chremoas/proto"
	"github.com/chremoas/services-common/args"
	"golang.org/x/net/context"
)

var CCF ClientFactory

type ClientFactory interface {
	PricingService() pricing.PricesService
	TypeService() sde.TypeQueryService
}

type Command struct {
	//Store anything you need the Help or Exec functions to have access to here
	name string
	args *args.Args
}

func (c *Command) Help(ctx context.Context, req *bot.HelpRequest, rsp *bot.HelpResponse) error {
	rsp.Usage = c.name
	rsp.Description = "Fetch prices of things from various places"
	return nil
}

func (c *Command) Exec(ctx context.Context, req *bot.ExecRequest, rsp *bot.ExecResponse) error {
	return c.args.Exec(ctx, req, rsp)
}

func NewCommand(name string, args *args.Args, cf ClientFactory) *Command {
	CCF = cf
	return &Command{name: name, args: args}
}
