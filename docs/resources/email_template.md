---
page_title: "ory_email_template Resource - Ory"
subcategory: ""
description: |-
  Manages an Ory email template.
---

# ory_email_template (Resource)

Manages an Ory email template. Email templates customize the content of emails sent during recovery, verification, and login flows.

~> **Note**: Deleting an email template resets it to Ory's default template rather than removing it entirely.

## Example Usage

### Recovery Code Email

```terraform
resource "ory_email_template" "recovery_code" {
  template_type = "recovery_code_valid"
  subject       = "Reset your password"

  body_html = <<-HTML
    <!DOCTYPE html>
    <html>
    <body>
      <h1>Password Reset</h1>
      <p>Hi {{ .Identity.traits.name.first }},</p>
      <p>Your recovery code is: <strong>{{ .RecoveryCode }}</strong></p>
      <p>This code expires in 1 hour.</p>
    </body>
    </html>
  HTML

  body_plaintext = <<-TEXT
    Password Reset

    Hi {{ .Identity.traits.name.first }},

    Your recovery code is: {{ .RecoveryCode }}

    This code expires in 1 hour.
  TEXT
}
```

### Verification Email

```terraform
resource "ory_email_template" "verification" {
  template_type = "verification_code_valid"
  subject       = "Verify your email address"

  body_html = <<-HTML
    <!DOCTYPE html>
    <html>
    <body>
      <h1>Email Verification</h1>
      <p>Welcome! Please verify your email address.</p>
      <p>Your verification code is: <strong>{{ .VerificationCode }}</strong></p>
    </body>
    </html>
  HTML

  body_plaintext = <<-TEXT
    Email Verification

    Welcome! Please verify your email address.

    Your verification code is: {{ .VerificationCode }}
  TEXT
}
```

### Login Code Email

```terraform
resource "ory_email_template" "login_code" {
  template_type = "login_code_valid"
  subject       = "Your login code"

  body_html = <<-HTML
    <!DOCTYPE html>
    <html>
    <body>
      <h1>Login Code</h1>
      <p>Your login code is: <strong>{{ .LoginCode }}</strong></p>
      <p>This code expires in 15 minutes.</p>
    </body>
    </html>
  HTML

  body_plaintext = <<-TEXT
    Login Code

    Your login code is: {{ .LoginCode }}

    This code expires in 15 minutes.
  TEXT
}
```

## Schema

### Required

- `template_type` (String) - Type of email template. See Template Types below.
- `subject` (String) - Email subject line. Supports Go template syntax.
- `body_html` (String) - HTML email body. Supports Go template syntax.
- `body_plaintext` (String) - Plain text email body. Supports Go template syntax.

### Optional

- `project_id` (String) - The project ID. If not set, uses the provider's project_id.

### Read-Only

- `id` (String) - Resource identifier.

## Template Types

| Type | Description |
|------|-------------|
| `recovery_code_valid` | Password recovery with code |
| `recovery_code_invalid` | Invalid/expired recovery code |
| `recovery_valid` | Password recovery with magic link |
| `recovery_invalid` | Invalid/expired recovery link |
| `verification_code_valid` | Email verification with code |
| `verification_code_invalid` | Invalid/expired verification code |
| `verification_valid` | Email verification with magic link |
| `verification_invalid` | Invalid/expired verification link |
| `login_code_valid` | Passwordless login with code |
| `login_code_invalid` | Invalid/expired login code |
| `registration_code_valid` | Registration with code |
| `registration_code_invalid` | Invalid/expired registration code |

## Template Variables

Templates use Go's `text/template` syntax. Available variables depend on the template type:

### Recovery Templates
- `{{ .RecoveryCode }}` - The recovery code
- `{{ .RecoveryURL }}` - The recovery link (for magic link templates)
- `{{ .Identity }}` - The identity object
- `{{ .Identity.traits.email }}` - Access identity traits

### Verification Templates
- `{{ .VerificationCode }}` - The verification code
- `{{ .VerificationURL }}` - The verification link
- `{{ .Identity }}` - The identity object

### Login Templates
- `{{ .LoginCode }}` - The login code
- `{{ .Identity }}` - The identity object

### Registration Templates
- `{{ .RegistrationCode }}` - The registration code
- `{{ .Identity }}` - The identity object (partial, during registration)

## Import

Email templates can be imported using the format `project_id:template_type`:

```bash
terraform import ory_email_template.recovery_code <project-id>:recovery_code_valid
```
