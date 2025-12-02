---
page_title: "ory_json_web_key_set Resource - Ory"
subcategory: ""
description: |-
  Manages an Ory JSON Web Key Set (JWKS).
---

# ory_json_web_key_set (Resource)

Manages an Ory JSON Web Key Set (JWKS). JWK Sets are used for signing and verifying tokens in OAuth2/OIDC flows.

## Example Usage

### RSA Signing Key

```terraform
resource "ory_json_web_key_set" "signing" {
  set_id    = "my-signing-keys"
  algorithm = "RS256"
  use       = "sig"
}
```

### Encryption Key

```terraform
resource "ory_json_web_key_set" "encryption" {
  set_id    = "my-encryption-keys"
  algorithm = "RS256"
  use       = "enc"
}
```

### ES256 Key (ECDSA)

```terraform
resource "ory_json_web_key_set" "ecdsa" {
  set_id    = "ecdsa-signing-keys"
  algorithm = "ES256"
  use       = "sig"
}
```

## Schema

### Required

- `set_id` (String) - Identifier for the key set.
- `algorithm` (String) - Key algorithm. Values: `RS256`, `RS384`, `RS512`, `ES256`, `ES384`, `ES512`, `EdDSA`.
- `use` (String) - Key usage. Values: `sig` (signing), `enc` (encryption).

### Optional

- `project_id` (String) - The project ID. If not set, uses the provider's project_id.

### Read-Only

- `id` (String) - Resource identifier (same as set_id).
- `keys` (List of Object) - The generated keys in the set.
  - `kid` (String) - Key ID.
  - `kty` (String) - Key type.
  - `alg` (String) - Algorithm.
  - `use` (String) - Key usage.

## Import

JWK Sets can be imported using the format `project_id:set_id`:

```bash
terraform import ory_json_web_key_set.signing <project-id>:my-signing-keys
```

## Algorithm Reference

| Algorithm | Type | Description |
|-----------|------|-------------|
| RS256 | RSA | RSA with SHA-256 |
| RS384 | RSA | RSA with SHA-384 |
| RS512 | RSA | RSA with SHA-512 |
| ES256 | ECDSA | ECDSA with P-256 and SHA-256 |
| ES384 | ECDSA | ECDSA with P-384 and SHA-384 |
| ES512 | ECDSA | ECDSA with P-521 and SHA-512 |
| EdDSA | EdDSA | Edwards-curve DSA (Ed25519) |
