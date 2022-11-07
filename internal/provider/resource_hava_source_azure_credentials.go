package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/teamhava/hava-sdk-go/havaclient"
)

func resourceHavaSourceAzureCredentials() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceSourceAzureCredentialsCreate,
		ReadContext:   resourceSourceAzureCredentialsRead,
		UpdateContext: resourceSourceAzureCredentialsUpdate,
		DeleteContext: resourceSourceAzureCredentialsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "Display name of the source",
				Type:        schema.TypeString,
				Required:    true,
			},
			"subscription_id": {
				Description: "The id of the azure subscription that will be accessed to import the data",
				Sensitive:   true,
				Type:        schema.TypeString,
				Required:    true,
			},
			"tenant_id": {
				Description: "The id of the azure tenant that will be accessed to import the data",
				Sensitive:   true,
				Type:        schema.TypeString,
				Required:    true,
			},
			"client_id": {
				Description: "The id of client that will be used to access the source for import",
				Sensitive:   true,
				Type:        schema.TypeString,
				Required:    true,
			},
			"secret_key": {
				Description: "The azure secret key of the client that will be used to access the source for import",
				Sensitive:   true,
				Type:        schema.TypeString,
				Required:    true,
			},
			"state": {
				Description: "State of the Source",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		CustomizeDiff: customdiff.All(

			// if state is set to archived, it has been deleted outside of terraform and a new resource needs to be created
			customdiff.ForceNewIf("state", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				state := d.Get("state").(string)

				if state == "archived" {
					tflog.Info(ctx, fmt.Sprintf("Source '%s' was deleted outside of terraform, a new resource will be created in it's place", d.Id()))
					d.SetNewComputed("state")
					return true
				}
				return false
			}),
		),
	}
}

func resourceSourceAzureCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	tflog.Info(ctx, "creating")

	client := meta.(*havaclient.APIClient)

	name := d.Get("name").(string)
	azureType := "Azure::Credentials"
	subId := d.Get("subscription_id").(string)
	tenantId := d.Get("tenant_id").(string)
	clientId := d.Get("client_id").(string)
	secretKey := d.Get("secret_key").(string)

	azureCredentialsSource := &havaclient.SourcesAzureCredentials{
		Name:           &name,
		Type:           &azureType,
		SubscriptionId: &subId,
		TenantId:       &tenantId,
		ClientId:       &clientId,
		SecretKey:      &secretKey,
	}

	body := havaclient.SourcesAzureCredentialsAsSourcesCreateRequest(azureCredentialsSource)

	req := client.SourcesApi.SourcesCreate(ctx)

	req = req.SourcesCreateRequest(body)

	source, _, err := req.Execute()

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*source.Id)
	d.Set("state", source.State)

	return nil
}

func resourceSourceAzureCredentialsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	tflog.Info(ctx, "reading")

	client := meta.(*havaclient.APIClient)

	req := client.SourcesApi.SourcesShow(ctx, d.Id())
	source, res, err := req.Execute()

	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, res.Status)

	d.Set("name", source.Name)
	d.Set("state", source.State)

	return nil
}

func resourceSourceAzureCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	tflog.Info(ctx, "updating")
	client := meta.(*havaclient.APIClient)

	name := d.Get("name").(string)
	azureType := "Azure::Credentials"
	subId := d.Get("subscription_id").(string)
	tenantId := d.Get("tenant_id").(string)
	clientId := d.Get("client_id").(string)
	secretKey := d.Get("secret_key").(string)

	azureCredentialsSource := &havaclient.SourcesAzureCredentials{
		Name:           &name,
		Type:           &azureType,
		SubscriptionId: &subId,
		TenantId:       &tenantId,
		ClientId:       &clientId,
		SecretKey:      &secretKey,
	}

	sourceUpdateRequest := havaclient.SourcesAzureCredentialsAsSourcesUpdateRequest(azureCredentialsSource)

	req := client.SourcesApi.SourcesUpdate(ctx, d.Id())

	req = req.SourcesUpdateRequest(sourceUpdateRequest)

	_, _, err := req.Execute()

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSourceAzureCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	tflog.Info(ctx, "deleting")

	client := meta.(*havaclient.APIClient)

	req := client.SourcesApi.SourcesDestroy(ctx, d.Id())

	_, _, err := req.Execute()

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
