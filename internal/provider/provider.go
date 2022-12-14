package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	havaclient "github.com/teamhava/hava-sdk-go"
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
			},
			ResourcesMap: map[string]*schema.Resource{
				"hava_source_aws_car_resource":           resourceHavaSourceAWSCAR(),
				"hava_source_aws_key_resource":           resourceHavaSourceAWSKey(),
				"hava_source_azure_credentials_resource": resourceHavaSourceAzureCredentials(),
				"hava_source_gcp_sa_credentials_resource": resourceHavaSourceGCPCredentials(),
			},
			Schema: map[string]*schema.Schema{
				"api_token": {
					Description: "The API token to authenticate with the Hava API. This takes precedence over the 'HAVA_TOKEN' environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("HAVA_TOKEN", nil),
				},
				"endpoint": {
					Description: "Which API endpoint to connect to. This is primarily used to support self-hosted users that does not use the default SaaS API endpoints",
					Type:	schema.TypeString,
					Optional: true,
					Default: "https://api.hava.io",
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
		endpoint := d.Get("endpoint").(string)

		if !ok {
			return nil, diag.Errorf("api token not found, did you set the 'HAVA_TOKEN' environment variable")
		}

		cfg := havaclient.NewConfiguration()
		cfg.Servers = havaclient.ServerConfigurations{
			{
				URL: endpoint,
				Description: "No description provided",
			},
		}

		cfg.UserAgent = p.UserAgent("terraform-provider-hava", version)

		cfg.DefaultHeader["Authorization"] = "Bearer " + token.(string)

		myclient := havaclient.NewAPIClient(cfg)

		return myclient, nil
	}
}
