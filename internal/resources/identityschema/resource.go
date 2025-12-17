package identityschema

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ory "github.com/ory/client-go"

	"github.com/ory/terraform-provider-orynetwork/internal/client"
)

var (
	_ resource.Resource                = &IdentitySchemaResource{}
	_ resource.ResourceWithConfigure   = &IdentitySchemaResource{}
	_ resource.ResourceWithImportState = &IdentitySchemaResource{}
)

func NewResource() resource.Resource {
	return &IdentitySchemaResource{}
}

type IdentitySchemaResource struct {
	client *client.OryClient
}

type IdentitySchemaResourceModel struct {
	ID         types.String `tfsdk:"id"`
	ProjectID  types.String `tfsdk:"project_id"`
	SchemaID   types.String `tfsdk:"schema_id"`
	Schema     types.String `tfsdk:"schema"`
	SetDefault types.Bool   `tfsdk:"set_default"`
}

func (r *IdentitySchemaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_schema"
}

func (r *IdentitySchemaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Ory Network identity schema.",
		MarkdownDescription: `
Manages an Ory Network identity schema.

**Note:** Identity schemas are immutable in Ory Network. Any changes to the schema content will require resource replacement.

**Note:** Ory Network does not support deleting identity schemas. When this resource is destroyed, the schema will remain in Ory but will no longer be managed by Terraform.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Resource ID.",
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
			"schema_id": schema.StringAttribute{
				Description: "Unique identifier for the schema (e.g., 'user', 'employee').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schema": schema.StringAttribute{
				Description: "JSON Schema definition for the identity traits (JSON string).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(), // Schemas are immutable
				},
			},
			"set_default": schema.BoolAttribute{
				Description: "Set this schema as the project's default schema.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *IdentitySchemaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IdentitySchemaResource) encodeSchema(schemaJSON string) (string, error) {
	// Validate it's valid JSON
	var schemaMap map[string]interface{}
	if err := json.Unmarshal([]byte(schemaJSON), &schemaMap); err != nil {
		return "", fmt.Errorf("invalid JSON schema: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(schemaJSON))
	return "base64://" + encoded, nil
}

func (r *IdentitySchemaResource) getSchemas(ctx context.Context, projectID string) ([]map[string]interface{}, error) {
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

	identity, _ := configMap["identity"].(map[string]interface{})
	schemas, _ := identity["schemas"].([]interface{})

	result := make([]map[string]interface{}, 0, len(schemas))
	for _, s := range schemas {
		if sm, ok := s.(map[string]interface{}); ok {
			result = append(result, sm)
		}
	}
	return result, nil
}

func (r *IdentitySchemaResource) findSchemaIndex(schemas []map[string]interface{}, schemaID string) int {
	for i, s := range schemas {
		if s["id"] == schemaID {
			return i
		}
	}
	return -1
}

func (r *IdentitySchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IdentitySchemaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := plan.ProjectID.ValueString()
	if projectID == "" {
		projectID = r.client.ProjectID()
	}

	schemaID := plan.SchemaID.ValueString()
	schemaJSON := plan.Schema.ValueString()

	schemaURL, err := r.encodeSchema(schemaJSON)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Schema", err.Error())
		return
	}

	schemas, err := r.getSchemas(ctx, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error Getting Schemas", err.Error())
		return
	}

	var patches []ory.JsonPatch

	existingIndex := r.findSchemaIndex(schemas, schemaID)
	if existingIndex >= 0 {
		// Replace existing
		patches = append(patches, ory.JsonPatch{
			Op:   "replace",
			Path: fmt.Sprintf("/services/identity/config/identity/schemas/%d", existingIndex),
			Value: map[string]string{
				"id":  schemaID,
				"url": schemaURL,
			},
		})
	} else {
		// Add new schema
		patches = append(patches, ory.JsonPatch{
			Op:   "add",
			Path: "/services/identity/config/identity/schemas/-",
			Value: map[string]string{
				"id":  schemaID,
				"url": schemaURL,
			},
		})
	}

	_, err = r.client.PatchProject(ctx, projectID, patches)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Identity Schema", err.Error())
		return
	}

	// Set as default if requested
	if plan.SetDefault.ValueBool() {
		defaultPatches := []ory.JsonPatch{{
			Op:    "replace",
			Path:  "/services/identity/config/identity/default_schema_id",
			Value: schemaID,
		}}
		_, err = r.client.PatchProject(ctx, projectID, defaultPatches)
		if err != nil {
			resp.Diagnostics.AddError("Error Setting Default Schema", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(schemaID)
	plan.ProjectID = types.StringValue(projectID)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *IdentitySchemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IdentitySchemaResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := state.ProjectID.ValueString()
	schemaID := state.SchemaID.ValueString()

	schemas, err := r.getSchemas(ctx, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Identity Schema", err.Error())
		return
	}

	index := r.findSchemaIndex(schemas, schemaID)
	if index < 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *IdentitySchemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IdentitySchemaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := plan.ProjectID.ValueString()
	if projectID == "" {
		projectID = r.client.ProjectID()
	}

	schemaID := plan.SchemaID.ValueString()

	// Only thing that can be updated without replacement is set_default
	if plan.SetDefault.ValueBool() {
		patches := []ory.JsonPatch{{
			Op:    "replace",
			Path:  "/services/identity/config/identity/default_schema_id",
			Value: schemaID,
		}}
		_, err := r.client.PatchProject(ctx, projectID, patches)
		if err != nil {
			resp.Diagnostics.AddError("Error Setting Default Schema", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(schemaID)
	plan.ProjectID = types.StringValue(projectID)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *IdentitySchemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Ory Network does not support deleting identity schemas.
	// See: https://github.com/ory/network/issues/262
	//
	// We just remove it from Terraform state. The schema will remain in Ory.
	resp.Diagnostics.AddWarning(
		"Schema Not Deleted",
		"Ory Network does not support deleting identity schemas. The schema has been removed from Terraform state but still exists in Ory Network.",
	)
}

func (r *IdentitySchemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("schema_id"), req.ID)...)
}
