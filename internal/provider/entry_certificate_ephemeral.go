package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResource = &EntryCertificateEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigure = &EntryCertificateEphemeralResource{}

func NewEntryCertificateEphemeralResource() ephemeral.EphemeralResource {
	return &EntryCertificateEphemeralResource{}
}

type EntryCertificateEphemeralResource struct {
	ephemeralResourceBase
}

type EntryCertificateEphemeralResourceModel = EntryCertificateDataSourceModel

func (e *EntryCertificateEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_certificate"
}

func (e *EntryCertificateEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS Certificate Entry, fetched ephemerally so the certificate file and password never land in Terraform state.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Certificate ID",
				Required:    true,
				Validators:  []validator.String{entryIdValidator{}},
			},
			"vault_id": schema.StringAttribute{
				Description: "Vault ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Certificate name",
				Computed:    true,
			},
			"folder": schema.StringAttribute{
				Description: "Certificate folder path",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Certificate description",
				Computed:    true,
			},
			"expiration": schema.StringAttribute{
				CustomType:  timetypes.RFC3339Type{},
				Description: "Certificate expiration date, in RFC3339 format (e.g. 2022-12-31T23:59:59-05:00)",
				Computed:    true,
			},
			"tags": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "Certificate tags",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Certificate password",
				Computed:    true,
				Sensitive:   true,
			},
			"file": schema.SingleNestedAttribute{
				Description: "Certificate file. Either file or url is populated.",
				Computed:    true,
				Sensitive:   true,
				Attributes: map[string]schema.Attribute{
					"content_b64": schema.StringAttribute{
						Description: "Certificate base 64 encoded string",
						Computed:    true,
						Sensitive:   true,
					},
					"name": schema.StringAttribute{
						Description: "Certificate file name",
						Computed:    true,
					},
				},
			},
			"url": schema.SingleNestedAttribute{
				Description: "Certificate url. Either file or url is populated.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Description: "Certificate url",
						Computed:    true,
					},
					"use_default_credentials": schema.BoolAttribute{
						Description: "Use default credentials",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (e *EntryCertificateEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data *EntryCertificateEphemeralResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry, content, err := fetchCertificateEntry(e.client, data.Id.ValueString())
	if err != nil {
		if errors.Is(err, errCertificateNotFound) {
			resp.Diagnostics.AddError("certificate entry not found", fmt.Sprintf("no certificate entry with id %q", data.Id.ValueString()))
			return
		}
		resp.Diagnostics.AddError("unable to read certificate entry", err.Error())
		return
	}

	resp.Diagnostics.Append(setEntryCertificateDataModel(ctx, entry, data, content)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
