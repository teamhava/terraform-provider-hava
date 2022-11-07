package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/teamhava/hava-sdk-go/havaclient"
)

func resourceHavaSourceAWSCAR() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceSourceAWSCARCreate,
		ReadContext:   resourceSourceAWSCARRead,
		UpdateContext: resourceSourceAWSCARUpdate,
		DeleteContext: resourceSourceAWSCARDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "Display name of the source",
				Type:        schema.TypeString,
				Required:    true,
			},
			"role_arn": {
				Description: "The ARN of the role that hava will assume to access the AWS Account",
				Sensitive: true,
				Type: schema.TypeString,
				Required: true,
			},
			"external_id": {
				Description: "The external ID used by AWS for additional security when assuming the role",
				Sensitive: true,
				Type: schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSourceAWSCARCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*havaclient.APIClient)
	
	name := d.Get("name").(string)
	role := d.Get("role_arn").(string)
	awsType := "AWS::CrossAccountRole"
	externalId := d.Get("external_id").(string)
	
	sawscar := &havaclient.SourcesAWSCAR{
		Name: &name,
		RoleArn: &role,
		Type: &awsType,
		ExternalId: &externalId,
	}

	body := havaclient.SourcesAWSCARAsSourcesCreateRequest(sawscar)
	

	b,_ := json.Marshal(body)

	tflog.Info(ctx, string(b))
	tflog.Info(ctx, fmt.Sprintf("body: %+v", body))

	req := client.SourcesApi.SourcesCreate(ctx)

	req = req.SourcesCreateRequest(body)

	source, res, err := req.Execute()

	if res != nil {
		tflog.Info(ctx, res.Status)
	}

	if err != nil {
		tflog.Info(ctx, err.Error())
		tflog.Error(ctx, fmt.Sprintf("%+v", err))
		return diag.FromErr(err)
	}

	// tflog.Info(ctx, (*res).Status)
	
	// if (*res).StatusCode != 200 {
	// 	return diag.Errorf("Hava API responded with a %d status code", res.StatusCode)
	// }


	d.SetId(*source.Id)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	// return diag.Errorf("create not implemented")

	return nil
}

func resourceSourceAWSCARRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	
	client := meta.(*havaclient.APIClient)

	req := client.SourcesApi.SourcesShow(ctx, d.Id())
	source, _, err := req.Execute()


	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", source.Name)

	return nil
}

func resourceSourceAWSCARUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	
	client := meta.(*havaclient.APIClient)

	name := d.Get("name").(string)
	role := d.Get("role_arn").(string)
	awsType := "AWS::CrossAccountRole"
	externalId := d.Get("external_id").(string)
	
	sawscar := &havaclient.SourcesAWSCAR{
		Name: &name,
		RoleArn: &role,
		Type: &awsType,
		ExternalId: &externalId,
	}

	sourceUpdateRequest := havaclient.SourcesAWSCARAsSourcesUpdateRequest(sawscar)

	req := client.SourcesApi.SourcesUpdate(ctx, d.Id())

	req = req.SourcesUpdateRequest(sourceUpdateRequest)

	_, _, err := req.Execute()

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSourceAWSCARDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	
	client := meta.(*havaclient.APIClient)

	req := client.SourcesApi.SourcesDestroy(ctx, d.Id())

	_, _, err := req.Execute()

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
