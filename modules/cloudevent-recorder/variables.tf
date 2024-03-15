variable "project_id" {
  type = string
}

variable "name" {
  type = string
}

variable "location" {
  default     = "US"
  description = "The location to create the BigQuery dataset in, and in which to run the data transfer jobs from GCS."
}

variable "provisioner" {
  type        = string
  description = "The identity as which this module will be applied (so it may be granted permission to 'act as' the DTS service account).  This should be in the form expected by an IAM subject (e.g. user:sally@example.com)"
}

variable "retention-period" {
  type        = number
  description = "The number of days to retain data in BigQuery."
}

variable "deletion_protection" {
  default     = true
  description = "Whether to enable deletion protection on data resources."
}

variable "regions" {
  description = "A map from region names to a network and subnetwork.  A recorder service and cloud storage bucket (into which the service writes events) will be created in each region."
  type = map(object({
    network = string
    subnet  = string
  }))
}

variable "broker" {
  type        = map(string)
  description = "A map from each of the input region names to the name of the Broker topic in that region."
}

variable "notification_channels" {
  description = "List of notification channels to alert (for service-level issues)."
  type        = list(string)
}

variable "types" {
  description = "A map from cloudevent types to the BigQuery schema associated with them, as well as an alert threshold and a list of notification channels (for subscription-level issues)."

  type = map(object({
    schema                = string
    alert_threshold       = optional(number, 50000)
    notification_channels = optional(list(string), [])
    partition_field       = optional(string)
    table_id              = optional(string)
    create_table          = optional(bool, true)
  }))
}

variable "dataset_id" {
  description = "The name of the BigQuery dataset to create."
  type        = string
  default     = null
}

variable "create_dataset" {
  description = "Whether to create the BigQuery dataset. Set to false if the dataset already exists."
  type        = bool
  default     = true
}