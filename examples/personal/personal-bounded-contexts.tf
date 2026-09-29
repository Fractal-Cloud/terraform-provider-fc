resource "fc_personal_bounded_context" "new_bc" {
  short_name   = "new-bc"
  display_name = "New Bounded Context"
  description  = "Bounded Context created through Terraform provider"
}

data "fc_personal_bounded_context" "existing_bounded_context" {
  short_name = "existing-bounded-context"
}

output "sandbox_bc" {
  value = data.fc_personal_bounded_context.existing_bounded_context
}
