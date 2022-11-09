package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	havaclient "github.com/teamhava/hava-sdk-go"
)

func resourceHavaSourceGCPCredentials() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A Source in Hava using an Google service account to authenticate to the GCP project that will be imported.",

		CreateContext: resourceSourceGCPCredentialsCreate,
		ReadContext:   resourceSourceGCPCredentialsRead,
		UpdateContext: resourceSourceGCPCredentialsUpdate,
		DeleteContext: resourceSourceGCPCredentialsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "Display name of the source",
				Type:        schema.TypeString,
				Required:    true,
			},
			"encoded_file": {
				Description: "Base64 encoded json Service Account credentials file content",
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

func resourceSourceGCPCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	tflog.Info(ctx, "creating")

	client := meta.(*havaclient.APIClient)

	name := d.Get("name").(string)
	gcpType := "GCP::ServiceAccountCredentials"
	encodedFile := d.Get("encoded_file").(string)

	gcpCredentialsSource := &havaclient.SourcesGCPServiceAccountCredentials{
		Name:           &name,
		Type:           &gcpType,
		EncodedFile: &encodedFile,
	}

	body := havaclient.SourcesGCPServiceAccountCredentialsAsSourcesCreateRequest(gcpCredentialsSource)

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

func resourceSourceGCPCredentialsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func resourceSourceGCPCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	tflog.Info(ctx, "updating")
	client := meta.(*havaclient.APIClient)

	name := d.Get("name").(string)
	gcpType := "GCP::ServiceAccountCredentials"
	encodedFile := d.Get("encoded_file").(string)

	gcpCredentialsSource := &havaclient.SourcesGCPServiceAccountCredentials{
		Name:           &name,
		Type:           &gcpType,
		EncodedFile: &encodedFile,
	}

	sourceUpdateRequest := havaclient.SourcesGCPServiceAccountCredentialsAsSourcesUpdateRequest(gcpCredentialsSource)

	req := client.SourcesApi.SourcesUpdate(ctx, d.Id())

	req = req.SourcesUpdateRequest(sourceUpdateRequest)

	_, _, err := req.Execute()

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSourceGCPCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	tflog.Info(ctx, "deleting")

	client := meta.(*havaclient.APIClient)

	req := client.SourcesApi.SourcesDestroy(ctx, d.Id())

	_, _, err := req.Execute()

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
