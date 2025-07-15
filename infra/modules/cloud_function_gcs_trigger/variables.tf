# modules/cloud_function_gcs_trigger

variable "project_id" {
  description = "GCP project ID."
  type        = string
}

variable "region" {
  description = "Region where Cloud Function will be deployed."
  type        = string
}

variable "name" {
  description = "A unique name for the Cloud Function."
  type        = string
}

variable "source_code_bucket_name" {
  description = "Name of the GCS bucket where the Cloud Function source code will be uploaded."
  type        = string
}

variable "source_code_path" {
  description = "Local path to the zipped Cloud Function source code."
  type        = string
}

variable "runtime" {
  description = "The runtime environment for the Cloud Function (e.g., 'python311', 'go121')."
  type        = string
}

variable "entry_point" {
  description = "The name of the function (as defined in your source code) to be executed."
  type        = string
}

variable "memory" {
  description = "The amount of memory allocated to the Cloud Function (e.g., '256MiB', '1GiB')."
  type        = string
  default     = "256MiB"
}

variable "timeout_seconds" {
  description = "The maximum execution time for the Cloud Function in seconds."
  type        = number
  default     = 60
}

variable "environment_variables" {
  description = "A map of environment variables to be set for the Cloud Function."
  type        = map(string)
  default     = {}
}

variable "max_instance_count" {
  description = "The maximum number of instances for the Cloud Function."
  type        = number
  default     = 2
}

variable "min_instance_count" {
  description = "The minimum number of instances for the Cloud Function (for warm start)."
  type        = number
  default     = 0
}

variable "trigger_region" {
  description = "The region where the Eventarc trigger will be created. Can be the same as function region."
  type        = string
}

variable "event_type" {
  description = "The CloudEvent type to trigger on (e.g., 'google.cloud.storage.object.v1.finalized')."
  type        = string
}

variable "retry_policy" {
  description = "The retry policy for the Cloud Function event trigger"
  type        = string
  default     = "RETRY_POLICY_DO_NOT_RETRY"
}

variable "event_resource_id" {
  description = "The full resource ID for the event trigger (e.g., 'projects/PROJECT_ID/buckets/BUCKET_NAME')."
  type        = string
}

variable "event_attribute_filters" {
  description = "A map of attribute-value pairs to filter events (e.g., { 'name': 'raw_vtt/' })."
  type        = map(string)
  default     = {}
}
