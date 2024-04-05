// Create an identity as which the recorder service will run.
resource "google_service_account" "recorder" {
  project = var.project_id

  # This GSA doesn't need it's own audit rule because it is used in conjunction
  # with regional-go-service, which has a built-in audit rule.

  account_id   = var.name
  display_name = "Cloudevents recorder"
  description  = "Dedicated service account for our recorder service."
}

// If method is gcs, Create a Pub/Sub service identity for the project
resource "google_project_service_identity" "pubsub" {
  provider = google-beta
  project  = var.project_id
  service  = "pubsub.googleapis.com"
}

// If method is gcs, GCP subscription will write events directly to a GCS bucket.
resource "google_storage_bucket_iam_binding" "broker-writes-to-gcs-buckets" {
  for_each = var.method == "gcs" ? var.regions : {}

  bucket  = google_storage_bucket.recorder[each.key].name
  role    = "roles/storage.admin"
  members = ["serviceAccount:${google_project_service_identity.pubsub.email}"]
}

// The recorder service account is the only identity that should be writing
// to the regional GCS buckets.
resource "google_storage_bucket_iam_binding" "recorder-writes-to-gcs-buckets" {
  for_each = var.method == "trigger" ? var.regions : {}

  bucket  = google_storage_bucket.recorder[each.key].name
  role    = "roles/storage.admin"
  members = ["serviceAccount:${google_service_account.recorder.email}"]
}

module "this" {
  count      = var.method == "trigger" ? 1 : 0
  source     = "../regional-go-service"
  project_id = var.project_id
  name       = var.name
  regions    = var.regions

  service_account = google_service_account.recorder.email
  containers = {
    "recorder" = {
      source = {
        working_dir = path.module
        importpath  = "./cmd/recorder"
      }
      ports = [{ container_port = 8080 }]
      env = [{
        name  = "LOG_PATH"
        value = "/logs"
      }]
      volume_mounts = [{
        name       = "logs"
        mount_path = "/logs"
      }]
    }
    "logrotate" = {
      source = {
        working_dir = path.module
        importpath  = "./cmd/logrotate"
      }
      env = [{
        name  = "LOG_PATH"
        value = "/logs"
      }]
      regional-env = [{
        name  = "BUCKET"
        value = { for k, v in google_storage_bucket.recorder : k => v.url }
      }]
      volume_mounts = [{
        name       = "logs"
        mount_path = "/logs"
      }]
    }
  }
  volumes = [{
    name      = "logs"
    empty_dir = {}
  }]

  notification_channels = var.notification_channels
}

resource "random_id" "trigger-suffix" {
  for_each    = var.types
  byte_length = 2
}

// Create a trigger for each region x type that sends events to the recorder service.
module "triggers" {
  for_each = var.method == "trigger" ? local.regional-types : {}

  source = "../cloudevent-trigger"

  name       = "${var.name}-${random_id.trigger-suffix[each.value.type].hex}"
  project_id = var.project_id
  broker     = var.broker[each.value.region]
  filter     = { "type" : each.value.type }

  depends_on = [module.this]
  private-service = {
    region = each.value.region
    name   = var.name
  }

  notification_channels = var.notification_channels
}


module "sync_gcs" {
  for_each = var.method == "gcs" ? local.regional-types : {}

  source = "../cloudevent-sync-gcs"

  name       = "${var.name}-${random_id.trigger-suffix[each.value.type].hex}"
  project_id = var.project_id
  broker     = var.broker[each.value.region]
  filter     = { "type" : each.value.type }
  bucket     = google_storage_bucket.recorder[each.value.region].name

  private-service = {
    region = each.value.region
    name   = var.name
  }

  notification_channels = var.notification_channels
}

locals {
  alerts = tomap({
    for type, schema in var.types :
    "BQ DTS ${var.name}-${type}" => google_monitoring_alert_policy.bq_dts[type].id
    if schema.create_table == null || schema.create_table
  })
}

module "recorder-dashboard" {
  source       = "../dashboard/cloudevent-receiver"
  service_name = var.name
  project_id   = var.project_id

  labels = { for type, schema in var.types : replace(type, ".", "_") => "" }

  notification_channels = var.notification_channels

  alerts = local.alerts

  triggers = {
    for type, schema in var.types : "type: ${type}" => {
      subscription_prefix   = "${var.name}-${random_id.trigger-suffix[type].hex}"
      alert_threshold       = schema.alert_threshold
      notification_channels = schema.notification_channels
    }
  }
}
