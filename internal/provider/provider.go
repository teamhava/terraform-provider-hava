package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/teamhava/hava-sdk-go/havaclient"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"scaffolding_data_source": dataSourceScaffolding(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"scaffolding_resource": resourceScaffolding(),
				"hava_source_aws_car_resource": resourceHavaSourceAWSCAR(),
				"hava_source_aws_key_resource": resourceHavaSourceAWSKey(),
			},
			Schema: map[string]*schema.Schema{
				"api_token": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HAVA_TOKEN", nil),
      	},
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {

		token, ok := d.GetOk("api_token")

		if !ok {
			return nil, diag.Errorf("api token not found, did you set the 'HAVA_TOKEN' environment variable")
		}

		cfg := havaclient.NewConfiguration()

		cfg.UserAgent = p.UserAgent("terraform-provider-hava", version)

		cfg.DefaultHeader["Authorization"] = "Bearer " + token.(string)

		myclient := havaclient.NewAPIClient(cfg)

		return myclient, nil
	}
}
