package command

import (
	"fmt"
	"github.com/abaeve/pricing-service/proto"
	"github.com/abaeve/sde-service/proto"
	bot "github.com/chremoas/chremoas/proto"
	"github.com/chremoas/services-common/args"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
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
	Min float32 `json:"min,omitempty"`
	Max float32 `json:"max,omitempty"`
	Avg float32 `json:"avg,omitempty"`
	Vol float32 `json:"vol,omitempty"`
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
	fmt.Printf("Exec: %v\n", request)

	if len(request.Args) < 3 {
		return "Not enough arguments"
	}

	search := strings.Join(request.Args[2:], " ")

	typeId, err := findTypeId(search)
	if err != nil {
		return "Disambiguate please..."
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
		"Min: %.2f\n"+
		"Avg: %.2f\n"+
		"Max: %.2f\n"+
		"Vol: %.2f\n"+
		"Ord: %d\n"+
		"```\n"+
		"```\n"+
		"Buy\n"+
		"Min: %.2f\n"+
		"Avg: %.2f\n"+
		"Max: %.2f\n"+
		"Vol: %.2f\n"+
		"Ord: %d\n"+
		"```", ip.ItemName, ip.Sell.Min, ip.Sell.Avg, ip.Sell.Max, ip.Sell.Vol, ip.Sell.Ord, ip.Buy.Min, ip.Buy.Avg, ip.Buy.Max, ip.Buy.Vol, ip.Buy.Ord)
}

func findTypeId(search string) (int32, error) {
	fmt.Printf("Search: %s\n", search)
	typeRsp, err := CCF.TypeService().FindTypesByTypeNames(context.Background(), &sde.TypeNameRequest{
		TypeName: []string{search},
	})
	if err != nil {
		return 0, err
	}

	if len(typeRsp.Type) > 1 {
		return typeRsp.Type[0].TypeId, errors.New("too many responses")
	}

	return typeRsp.Type[0].TypeId, nil
}

func fetchPrice(ctx context.Context, regionId, typeId int32) (*item, error) {
	fmt.Printf("Region: %d, Type: %d\n", regionId, typeId)
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

func fetchTypeIdFromTypeName(ctx context.Context, typeName string) (int32, error) {
	return 0, nil
}
