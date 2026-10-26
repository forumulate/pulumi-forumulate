package main

import (
	forumulate "github.com/forumulate/pulumi-forumulate/sdk/go/forumulate"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// FORUMULATE_TOKEN and FORUMULATE_ORG_ID come from the environment (or
// `pulumi config set --secret forumulate:token ...`).
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		spaceID := config.Require(ctx, "statusSpaceId")

		groups, err := forumulate.NewStatusGroups(ctx, "page", &forumulate.StatusGroupsArgs{
			SpaceId: pulumi.String(spaceID),
			Groups: forumulate.StatusGroupsGroupArray{
				&forumulate.StatusGroupsGroupArgs{Name: pulumi.String("Core"), Description: pulumi.String("What every customer touches.")},
				&forumulate.StatusGroupsGroupArgs{Name: pulumi.String("Integrations"), Layout: pulumi.String("grid-2"), HideUptime: pulumi.Bool(true)},
			},
		})
		if err != nil {
			return err
		}

		api, err := forumulate.NewStatusComponent(ctx, "api", &forumulate.StatusComponentArgs{
			SpaceId:         pulumi.String(spaceID),
			Name:            pulumi.String("API"),
			Url:             pulumi.String("https://api.example.com/health"),
			Group:           pulumi.String("Core"),
			Regions:         pulumi.StringArray{pulumi.String("eu-west"), pulumi.String("us-east")},
			ConfirmRegions:  pulumi.Int(2),
			BarDays:         pulumi.Int(90),
			IntervalSeconds: pulumi.Int(60),
		}, pulumi.DependsOn([]pulumi.Resource{groups}))
		if err != nil {
			return err
		}

		ctx.Export("apiComponent", api.ID())
		return nil
	})
}
