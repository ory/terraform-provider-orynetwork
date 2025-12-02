// Copyright 2025 Materialize Inc. and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package project

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/jasonhernandez/terraform-provider-orynetwork/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ProjectResource{}
	_ resource.ResourceWithConfigure   = &ProjectResource{}
	_ resource.ResourceWithImportState = &ProjectResource{}
)

// NewResource returns a new Project resource.
func NewResource() resource.Resource {
	return &ProjectResource{}
}

// ProjectResource defines the resource implementation.
type ProjectResource struct {
	client *client.OryClient
}

// ProjectResourceModel describes the resource data model.
type ProjectResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Environment types.String `tfsdk:"environment"`
	Slug        types.String `tfsdk:"slug"`
	State       types.String `tfsdk:"state"`
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Ory Network project.",
		MarkdownDescription: `
Manages an Ory Network project.

Projects are the top-level resource in Ory Network. Each project has its own
identity service, OAuth2 server, and configuration.

## Example Usage

` + "```hcl" + `
resource "ory_project" "main" {
  name        = "My Application"
  environment = "prod"
}
` + "```" + `

## Import

Projects can be imported using their ID:

` + "```shell" + `
terraform import ory_project.main <project-id>
` + "```" + `
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the project.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the project.",
				Required:    true,
			},
			"environment": schema.StringAttribute{
				Description: "The environment type: prod, stage, or dev.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("prod"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"slug": schema.StringAttribute{
				Description: "The project slug (e.g., 'vibrant-moore-abc123').",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description: "The project state (e.g., 'running').",
				Computed:    true,
			},
		},
	}
}

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	oryClient, ok := req.ProviderData.(*client.OryClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.OryClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = oryClient
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProjectResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.CreateProject(ctx, plan.Name.ValueString(), plan.Environment.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Project",
			"Could not create project: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(project.GetId())
	plan.Name = types.StringValue(project.GetName())
	plan.Environment = types.StringValue(project.GetEnvironment())
	plan.Slug = types.StringValue(project.GetSlug())
	plan.State = types.StringValue(project.GetState())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProjectResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.GetProject(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Project",
			"Could not read project ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(project.GetName())
	state.Environment = types.StringValue(project.GetEnvironment())
	state.Slug = types.StringValue(project.GetSlug())
	state.State = types.StringValue(project.GetState())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProjectResourceModel
	var state ProjectResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ory doesn't have a direct update endpoint for project name
	// For now, we just read back the current state
	project, err := r.client.GetProject(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Project",
			"Could not read project for update: "+err.Error(),
		)
		return
	}

	plan.ID = state.ID
	plan.Slug = types.StringValue(project.GetSlug())
	plan.State = types.StringValue(project.GetState())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProjectResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteProject(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Project",
			"Could not delete project: "+err.Error(),
		)
		return
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
