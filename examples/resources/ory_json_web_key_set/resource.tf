# RSA signing key set
resource "ory_json_web_key_set" "signing" {
  set_id    = "token-signing-keys"
  algorithm = "RS256"
  use       = "sig"
}

# ECDSA signing key set (smaller, faster)
resource "ory_json_web_key_set" "ecdsa_signing" {
  set_id    = "ecdsa-signing-keys"
  algorithm = "ES256"
  use       = "sig"
}

# Encryption key set
resource "ory_json_web_key_set" "encryption" {
  set_id    = "encryption-keys"
  algorithm = "RS256"
  use       = "enc"
}

output "signing_key_set_id" {
  value = ory_json_web_key_set.signing.id
}
