---
page_title: "secret_ref Function - Fractal Cloud"
subcategory: ""
description: |-
  References an environment secret by short name.
---

# function: secret_ref

Returns a reference to a secret defined on the environment. Use it as a component parameter or link-setting value instead of the raw secret. The agent resolves the reference from the environment secret store at reconciliation time, so the secret itself never appears in the blueprint, the Terraform state or the control plane.

The value is the string `{"$envSecret":"<short_name>"}`. It is sent to the API as that JSON object, which is how agents recognize a reference.

## Example Usage

```terraform
locals {
  openai = provider::fc::ai_saas_unmanaged({
    id     = "openai"
    secret = provider::fc::secret_ref("openai-api-key")
  })

  # In any parameter or link setting
  gateway = provider::fc::api_management_caas_api_gateway({
    id = "gateway"
    extra_parameters = {
      licenseKey = provider::fc::secret_ref("gateway-license")
    }
  })
}
```

## Signature

```text
secret_ref(short_name string) string
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `short_name` | String | Yes | Short name of the secret on the environment. Must not be empty. |
