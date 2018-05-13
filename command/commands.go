package command

import (
	"fmt"
	"github.com/abaeve/pricing-service/proto"
	"github.com/abaeve/sde-service/proto"
	bot "github.com/chremoas/chremoas/proto"
	"github.com/chremoas/services-common/args"
	"github.com/dustin/go-humanize"
	"golang.org/x/net/context"
	"math"
	"strings"
)

var Commands map[string]*args.Command

func init() {
	Commands = make(map[string]*args.Command)

	Commands["jita"] = &args.Command{
		Funcptr: jita,
		Help:    "Fetch price's from Jita's region The Forge",
	}
	Commands["amarr"] = &args.Command{
		Funcptr: amarr,
		Help:    "Fetch price's from Amarr's region Domain",
	}
	Commands["rens"] = &args.Command{
		Funcptr: rens,
		Help:    "Fetch price's from Rens' region Heimatar",
	}
	Commands["dodi"] = &args.Command{
		Funcptr: dodi,
		Help:    "Fetch prices from Dodixie's region Sinq Laison",
	}
}

type item struct {
	RegionId int32     `json:"region_id,omitempty"`
	ItemName string    `json:"item_name,omitempty"`
	ItemId   int32     `json:"item_id,omitempty"`
	Buy      itemPrice `json:"buy,omitempty"`
	Sell     itemPrice `json:"sell,omitempty"`
}

type itemPrice struct {
	Min float64 `json:"min,omitempty"`
	Max float64 `json:"max,omitempty"`
	Avg float64 `json:"avg,omitempty"`
	Vol int64   `json:"vol,omitempty"`
	Ord int64   `json:"ord,omitempty"`
}

func jita(ctx context.Context, request *bot.ExecRequest) string {
	//10000002
	return execute(ctx, request, 10000002)
}

func amarr(ctx context.Context, request *bot.ExecRequest) string {
	//10000043
	return execute(ctx, request, 10000043)
}

func rens(ctx context.Context, request *bot.ExecRequest) string {
	//10000032
	return execute(ctx, request, 10000032)
}

func dodi(ctx context.Context, request *bot.ExecRequest) string {
	//10000030
	return execute(ctx, request, 10000030)
}

func execute(ctx context.Context, request *bot.ExecRequest, regionId int32) string {
	if len(request.Args) < 3 {
		return "Not enough arguments"
	}

	search := strings.Join(request.Args[2:], " ")

	typeId, err := findTypeId(search)
	if err != nil {
		return fmt.Sprintf("No type with name %s found", search)
	}

	return executeContinued(ctx, request, search, regionId, typeId)
}

func executeContinued(ctx context.Context, request *bot.ExecRequest, itemName string, regionId, typeId int32) string {
	ip, err := fetchPrice(ctx, regionId, typeId)
	if err != nil {
		return "Had an issue: " + err.Error()
	}

	ip.ItemName = itemName

	return format(ip)
}

func format(ip *item) string {
	return fmt.Sprintf("Item: %s\n"+
		"```\n"+
		"Sell\n"+
		"Min: %s\n"+
		"Avg: %s\n"+
		"Max: %s\n"+
		"Vol: %s\n"+
		"Ord: %s\n"+
		"```\n"+
		"```\n"+
		"Buy\n"+
		"Max: %s\n"+
		"Avg: %s\n"+
		"Min: %s\n"+
		"Vol: %s\n"+
		"Ord: %s\n"+
		"```",
		ip.ItemName,
		humanize.Commaf(Round(ip.Sell.Min)),
		humanize.Commaf(Round(ip.Sell.Avg)),
		humanize.Commaf(Round(ip.Sell.Max)),
		humanize.Comma(ip.Sell.Vol),
		humanize.Comma(ip.Sell.Ord),
		humanize.Commaf(Round(ip.Buy.Max)),
		humanize.Commaf(Round(ip.Buy.Avg)),
		humanize.Commaf(Round(ip.Buy.Min)),
		humanize.Comma(ip.Buy.Vol),
		humanize.Comma(ip.Buy.Ord))
}

func findTypeId(search string) (int32, error) {
	typeRsp, err := CCF.TypeService().FindTypeByTypeName(context.Background(), &sde.TypeNameRequest{
		TypeName: search,
	})
	if err != nil {
		return 0, err
	}

	return typeRsp.Type.TypeId, nil
}

func fetchPrice(ctx context.Context, regionId, typeId int32) (*item, error) {
	psItemPrice, err := CCF.PricingService().GetItemPrice(ctx, &pricing.ItemPriceRequest{RegionId: regionId, ItemId: typeId})
	if err != nil {
		return nil, err
	}

	return &item{
		RegionId: regionId,
		ItemId:   typeId,
		Sell: itemPrice{
			Min: psItemPrice.Item.Sell.Min,
			Avg: psItemPrice.Item.Sell.Avg,
			Max: psItemPrice.Item.Sell.Max,
			Vol: psItemPrice.Item.Sell.Vol,
			Ord: psItemPrice.Item.Sell.Ord,
		},
		Buy: itemPrice{
			Min: psItemPrice.Item.Buy.Min,
			Avg: psItemPrice.Item.Buy.Avg,
			Max: psItemPrice.Item.Buy.Max,
			Vol: psItemPrice.Item.Buy.Vol,
			Ord: psItemPrice.Item.Buy.Ord,
		},
	}, nil
}

func Round(x float64) float64 {
	return math.Round(x*100) / 100
}
