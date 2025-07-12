# modules/network/main.tf

module "vpc_network" {
  source  = "terraform-google-modules/network/google"
  version = "~> 11.1"

  project_id   = var.project_id
  network_name = var.vpc_network_name
  routing_mode = "REGIONAL"

  subnets = [
    {
      subnet_name              = var.vpc_subnet_name
      subnet_ip                = var.vpc_subnet_ip_cidr_range
      subnet_region            = var.region
      private_ip_google_access = true # Required for Cloud SQL private access
    }
  ]

  ingress_rules = [
    {
      name          = "${var.vpc_network_name}-allow-internal"
      description   = "Allow internal traffic within the subnet"
      source_ranges = [var.vpc_subnet_ip_cidr_range]
      allow = [{
        protocol = "all"
      }]
    }
  ]

  egress_rules = [
    {
      name        = "${var.vpc_network_name}-allow-egress"
      description = "Allow all egress from VPC connector"
      # source_tags        = ["serverless-vpc-access"]
      destination_ranges = ["0.0.0.0/0"] # Allows outbound traffic to the internet and private IPs
      allow = [{
        protocol = "all"
      }]
    }
  ]
}

module "serverless_connector" {
  source     = "terraform-google-modules/network/google//modules/vpc-serverless-connector-beta"
  project_id = var.project_id

  vpc_connectors = [{
    name            = var.vpc_connector_name
    region          = var.region
    subnet_name     = var.vpc_subnet_name
    host_project_id = var.project_id
    machine_type    = "e2-micro"
    min_instances   = 2
    max_instances   = 3
  }]

  depends_on = [module.vpc_network]
}
