terraform {
  backend "s3" {
    bucket  = "linkpulse-mnq-terraform-state"
    key     = "linkpulse/dev/terraform.tfstate"
    region  = "ap-south-1"
    encrypt = true

    use_lockfile = true
  }
}
