resource "hava_source_gcp_sa_credentials_resource" "example" {
  name         = "Example Source"
  encoded_file = filebase64("./credentials.json")
}