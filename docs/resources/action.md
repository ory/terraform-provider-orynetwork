---
page_title: "ory_action Resource - Ory"
subcategory: ""
description: |-
  Manages an Ory Action (webhook) for identity flows.
---

# ory_action (Resource)

Manages an Ory Action (webhook) for identity flows. Actions allow you to trigger webhooks before or after identity flows like login, registration, recovery, settings, and verification.

## Example Usage

### Post-Registration Webhook

```terraform
resource "ory_action" "welcome_email" {
  flow        = "registration"
  timing      = "after"
  auth_method = "password"
  url         = "https://api.example.com/webhooks/welcome"
  method      = "POST"
  body        = <<-JSONNET
    function(ctx) {
      email: ctx.identity.traits.email,
      name: ctx.identity.traits.name,
      created_at: ctx.identity.created_at
    }
  JSONNET
}
```

### Pre-Login Validation

```terraform
resource "ory_action" "check_user_status" {
  flow         = "login"
  timing       = "before"
  url          = "https://api.example.com/webhooks/validate-login"
  method       = "POST"
  can_interrupt = true  # Allow webhook to block the flow
}
```

### Async Notification (Fire and Forget)

```terraform
resource "ory_action" "audit_log" {
  flow            = "settings"
  timing          = "after"
  auth_method     = "password"
  url             = "https://api.example.com/webhooks/audit"
  method          = "POST"
  response_ignore = true  # Don't wait for response
}
```

### Parse Response to Modify Identity

```terraform
resource "ory_action" "enrich_identity" {
  flow           = "registration"
  timing         = "after"
  auth_method    = "password"
  url            = "https://api.example.com/webhooks/enrich"
  method         = "POST"
  response_parse = true  # Parse response to modify identity
  body           = <<-JSONNET
    function(ctx) {
      email: ctx.identity.traits.email
    }
  JSONNET
}
```

## Schema

### Required

- `flow` (String) - Identity flow to hook into. Values: `login`, `registration`, `recovery`, `settings`, `verification`.
- `timing` (String) - When to trigger the webhook. Values: `before` (pre-hook), `after` (post-hook).
- `url` (String) - The webhook URL to call.

### Optional

- `project_id` (String) - The project ID. If not set, uses the provider's project_id.
- `auth_method` (String) - Authentication method to hook into. Required for `after` timing. Values: `password`, `oidc`, `code`, `webauthn`, `passkey`, `totp`, `lookup_secret`. Default: `password`.
- `method` (String) - HTTP method. Default: `POST`.
- `body` (String) - Jsonnet template for the request body. The template receives a `ctx` object with flow context.
- `response_ignore` (Boolean) - Run webhook asynchronously without waiting for response. Default: `false`.
- `response_parse` (Boolean) - Parse the webhook response to modify the identity. Default: `false`.
- `can_interrupt` (Boolean) - Allow the webhook to interrupt/block the flow by returning an error. Default: `false`.

### Read-Only

- `id` (String) - Resource identifier in format `project_id:flow:timing:auth_method:url`.

## Import

Actions can be imported using the format `project_id:flow:timing:auth_method:url`:

```bash
terraform import ory_action.welcome_email <project-id>:registration:after:password:https://api.example.com/webhooks/welcome
```

## Jsonnet Body Template

The `body` attribute accepts a Jsonnet template. The template receives a context object with information about the flow:

```jsonnet
function(ctx) {
  // Available fields depend on the flow and timing
  identity_id: ctx.identity.id,
  email: ctx.identity.traits.email,
  flow_id: ctx.flow.id,
  session_id: ctx.session.id  // Only available for some flows
}
```

### Context Variables

| Variable | Description |
|----------|-------------|
| `ctx.identity` | The identity object with traits, metadata, etc. |
| `ctx.flow` | The current flow object |
| `ctx.session` | Session object (when available) |
| `ctx.request_headers` | HTTP request headers |
| `ctx.request_url` | The request URL |
