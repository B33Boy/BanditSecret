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

  # Ingress Firewall Rules
  # These rules allow incoming traffic to resources within your VPC network.
  ingress_rules = [
    {
      name        = "${var.vpc_network_name}-allow-internal-from-subnet"
      description = "Allow all protocols and ports for internal traffic originating from within the primary subnet to targets within the VPC network."
      # Applies to all targets in the network (if target_tags/service_accounts are not set)
      source_ranges = [var.vpc_subnet_ip_cidr_range]
      allow = [{
        protocol = "all"
      }]
    },
    # You would add other specific ingress rules here, 
    # e.g. specific ports from external sources if your application needs them.
  ]

  # Egress Firewall Rules 
  # These rules allow outgoing traffic from resources within your VPC network.
  egress_rules = [
    {
      name        = "${var.vpc_network_name}-allow-vpc-connector-egress"
      description = "Allow all egress traffic from VMs associated with the Cloud Run VPC Connector to any destination."
      # The 'serverless-vpc-access' tag is automatically applied to VMs used by the VPC connector.
      source_tags        = ["serverless-vpc-access"]
      destination_ranges = ["0.0.0.0/0"] # Allows outbound traffic to the internet and private IPs
      allow = [{
        protocol = "all"
      }]
    },
    # You would add other specific egress rules here, 
    # e.g., restricting outbound to specific IPs/ports.
  ]

  depends_on = [var.api_dependencies]
}

// Enable cloud run to connect to resources inside VPC such as cloud sql
module "serverless_connector" {
  source     = "terraform-google-modules/network/google//modules/vpc-serverless-connector-beta"
  project_id = var.project_id
  vpc_connectors = [{
    name            = var.vpc_connector_name
    region          = var.region
    subnet_name     = var.vpc_subnet_name
    host_project_id = var.project_id
    machine_type    = "e2-small"
    min_instances   = 2
    max_instances   = 3
  }]

  # Ensure the VPC network is created and relevant APIs are enabled before creating the connector
  depends_on = [
    module.vpc_network,
    var.api_dependencies
  ]
}

# resource "google_project_service" "vpcaccess_api" {
#   project            = var.project_id
#   service            = "vpcaccess.googleapis.com"
#   disable_on_destroy = false
# }

# resource "google_project_service" "sqladmin_api" {
#   project            = var.project_id
#   service            = "sqladmin.googleapis.com"
#   disable_on_destroy = false
# }
