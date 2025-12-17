package socialprovider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ory "github.com/ory/client-go"

	"github.com/ory/terraform-provider-orynetwork/internal/client"
)

var (
	_ resource.Resource                = &SocialProviderResource{}
	_ resource.ResourceWithConfigure   = &SocialProviderResource{}
	_ resource.ResourceWithImportState = &SocialProviderResource{}
)

func NewResource() resource.Resource {
	return &SocialProviderResource{}
}

type SocialProviderResource struct {
	client *client.OryClient
}

type SocialProviderResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ProjectID    types.String `tfsdk:"project_id"`
	ProviderID   types.String `tfsdk:"provider_id"`
	ProviderType types.String `tfsdk:"provider_type"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	IssuerURL    types.String `tfsdk:"issuer_url"`
	Scope        types.List   `tfsdk:"scope"`
	MapperURL    types.String `tfsdk:"mapper_url"`
	AuthURL      types.String `tfsdk:"auth_url"`
	TokenURL     types.String `tfsdk:"token_url"`
	Tenant       types.String `tfsdk:"tenant"`
}

func (r *SocialProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_social_provider"
}

func (r *SocialProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Ory Network social sign-in provider (Google, GitHub, etc.).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Resource ID (same as provider_id).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "Project ID. If not set, uses provider's project_id.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"provider_id": schema.StringAttribute{
				Description: "Unique identifier for the provider (used in callback URLs).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"provider_type": schema.StringAttribute{
				Description: "Provider type (google, github, microsoft, apple, generic, etc.).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "OAuth2 client ID from the provider.",
				Required:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "OAuth2 client secret from the provider.",
				Required:    true,
				Sensitive:   true,
			},
			"issuer_url": schema.StringAttribute{
				Description: "OIDC issuer URL (required for generic providers).",
				Optional:    true,
			},
			"scope": schema.ListAttribute{
				Description: "OAuth2 scopes to request.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"mapper_url": schema.StringAttribute{
				Description: "Jsonnet mapper URL for claims mapping.",
				Optional:    true,
			},
			"auth_url": schema.StringAttribute{
				Description: "Custom authorization URL (for non-standard providers).",
				Optional:    true,
			},
			"token_url": schema.StringAttribute{
				Description: "Custom token URL (for non-standard providers).",
				Optional:    true,
			},
			"tenant": schema.StringAttribute{
				Description: "Tenant ID (for Microsoft/Azure providers).",
				Optional:    true,
			},
		},
	}
}

func (r *SocialProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	oryClient, ok := req.ProviderData.(*client.OryClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.OryClient, got: %T", req.ProviderData))
		return
	}
	r.client = oryClient
}

func (r *SocialProviderResource) buildProviderConfig(ctx context.Context, plan *SocialProviderResourceModel) map[string]interface{} {
	config := map[string]interface{}{
		"id":            plan.ProviderID.ValueString(),
		"provider":      plan.ProviderType.ValueString(),
		"client_id":     plan.ClientID.ValueString(),
		"client_secret": plan.ClientSecret.ValueString(),
	}

	if !plan.IssuerURL.IsNull() && !plan.IssuerURL.IsUnknown() {
		config["issuer_url"] = plan.IssuerURL.ValueString()
	}
	if !plan.Scope.IsNull() && !plan.Scope.IsUnknown() {
		var scope []string
		plan.Scope.ElementsAs(ctx, &scope, false)
		config["scope"] = scope
	}
	if !plan.MapperURL.IsNull() && !plan.MapperURL.IsUnknown() {
		config["mapper_url"] = plan.MapperURL.ValueString()
	}
	if !plan.AuthURL.IsNull() && !plan.AuthURL.IsUnknown() {
		config["auth_url"] = plan.AuthURL.ValueString()
	}
	if !plan.TokenURL.IsNull() && !plan.TokenURL.IsUnknown() {
		config["token_url"] = plan.TokenURL.ValueString()
	}
	if !plan.Tenant.IsNull() && !plan.Tenant.IsUnknown() {
		config["microsoft_tenant"] = plan.Tenant.ValueString()
	}

	return config
}

func (r *SocialProviderResource) getProviders(ctx context.Context, projectID string) ([]map[string]interface{}, error) {
	project, err := r.client.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if project.Services.Identity == nil {
		return nil, nil
	}

	configMap := project.Services.Identity.Config
	if configMap == nil {
		return nil, nil
	}

	selfservice, _ := configMap["selfservice"].(map[string]interface{})
	methods, _ := selfservice["methods"].(map[string]interface{})
	oidc, _ := methods["oidc"].(map[string]interface{})
	oidcConfig, _ := oidc["config"].(map[string]interface{})
	providers, _ := oidcConfig["providers"].([]interface{})

	result := make([]map[string]interface{}, 0, len(providers))
	for _, p := range providers {
		if pm, ok := p.(map[string]interface{}); ok {
			result = append(result, pm)
		}
	}
	return result, nil
}

func (r *SocialProviderResource) findProviderIndex(providers []map[string]interface{}, providerID string) int {
	for i, p := range providers {
		if p["id"] == providerID {
			return i
		}
	}
	return -1
}

func (r *SocialProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SocialProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := plan.ProjectID.ValueString()
	if projectID == "" {
		projectID = r.client.ProjectID()
	}

	providerConfig := r.buildProviderConfig(ctx, &plan)

	// Get current providers
	providers, err := r.getProviders(ctx, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error Getting Providers", err.Error())
		return
	}

	var patches []ory.JsonPatch

	existingIndex := r.findProviderIndex(providers, plan.ProviderID.ValueString())
	if existingIndex >= 0 {
		// Replace existing
		patches = append(patches, ory.JsonPatch{
			Op:    "replace",
			Path:  fmt.Sprintf("/services/identity/config/selfservice/methods/oidc/config/providers/%d", existingIndex),
			Value: providerConfig,
		})
	} else {
		// Enable OIDC if not already
		if len(providers) == 0 {
			patches = append(patches, ory.JsonPatch{
				Op:    "add",
				Path:  "/services/identity/config/selfservice/methods/oidc/enabled",
				Value: true,
			})
		}
		// Add new provider
		patches = append(patches, ory.JsonPatch{
			Op:    "add",
			Path:  "/services/identity/config/selfservice/methods/oidc/config/providers/-",
			Value: providerConfig,
		})
	}

	_, err = r.client.PatchProject(ctx, projectID, patches)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Social Provider", err.Error())
		return
	}

	plan.ID = plan.ProviderID
	plan.ProjectID = types.StringValue(projectID)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SocialProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SocialProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := state.ProjectID.ValueString()
	providers, err := r.getProviders(ctx, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Social Provider", err.Error())
		return
	}

	index := r.findProviderIndex(providers, state.ProviderID.ValueString())
	if index < 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	provider := providers[index]
	state.ProviderType = types.StringValue(fmt.Sprintf("%v", provider["provider"]))
	state.ClientID = types.StringValue(fmt.Sprintf("%v", provider["client_id"]))
	// Don't read back client_secret for security

	if issuer, ok := provider["issuer_url"].(string); ok {
		state.IssuerURL = types.StringValue(issuer)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SocialProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SocialProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := plan.ProjectID.ValueString()
	if projectID == "" {
		projectID = r.client.ProjectID()
	}

	providers, err := r.getProviders(ctx, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error Getting Providers", err.Error())
		return
	}

	index := r.findProviderIndex(providers, plan.ProviderID.ValueString())
	if index < 0 {
		resp.Diagnostics.AddError("Provider Not Found",
			fmt.Sprintf("Provider '%s' not found", plan.ProviderID.ValueString()))
		return
	}

	providerConfig := r.buildProviderConfig(ctx, &plan)
	patches := []ory.JsonPatch{{
		Op:    "replace",
		Path:  fmt.Sprintf("/services/identity/config/selfservice/methods/oidc/config/providers/%d", index),
		Value: providerConfig,
	}}

	_, err = r.client.PatchProject(ctx, projectID, patches)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Social Provider", err.Error())
		return
	}

	plan.ID = plan.ProviderID
	plan.ProjectID = types.StringValue(projectID)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SocialProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SocialProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := state.ProjectID.ValueString()
	providers, err := r.getProviders(ctx, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error Getting Providers", err.Error())
		return
	}

	index := r.findProviderIndex(providers, state.ProviderID.ValueString())
	if index < 0 {
		return // Already deleted
	}

	patches := []ory.JsonPatch{{
		Op:   "remove",
		Path: fmt.Sprintf("/services/identity/config/selfservice/methods/oidc/config/providers/%d", index),
	}}

	_, err = r.client.PatchProject(ctx, projectID, patches)
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Social Provider", err.Error())
		return
	}
}

func (r *SocialProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("provider_id"), req.ID)...)
}
