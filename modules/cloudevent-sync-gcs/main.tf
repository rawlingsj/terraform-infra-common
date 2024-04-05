resource "random_string" "suffix" {
  length  = 4
  upper   = false
  special = false
}

// A dedicated service account for this subscription.
resource "google_service_account" "this" {
  project = var.project_id

  account_id   = "${var.name}-${random_string.suffix.result}"
  display_name = "Delivery account for ${var.private-service.name} in ${var.private-service.region}."
}

// Lookup the identity of the pubsub service agent.
resource "google_project_service_identity" "pubsub" {
  provider = google-beta
  project  = var.project_id
  service  = "pubsub.googleapis.com"
}

locals {
  // See https://cloud.google.com/pubsub/docs/subscription-message-filter#filtering_syntax
  filter-elements = concat(
    [for key, value in var.filter : "attributes.ce-${key}=\"${value}\""],
    [for key, value in var.filter_prefix : "hasPrefix(attributes.ce-${key}, \"${value}\")"],
  )
}

resource "google_pubsub_topic" "dead-letter" {
  name = "${var.name}-dlq-${random_string.suffix.result}"

  message_storage_policy {
    allowed_persistence_regions = [var.private-service.region]
  }
}

// Grant the pubsub service account the ability to send to the dead-letter topic.
resource "google_pubsub_topic_iam_binding" "allow-pubsub-to-send-to-dead-letter" {
  topic = google_pubsub_topic.dead-letter.name

  role    = "roles/pubsub.publisher"
  members = ["serviceAccount:${google_project_service_identity.pubsub.email}"]
}

// Configure the subscription to deliver the events matching our filter to this service
// using the above identity to authorize the delivery..
resource "google_pubsub_subscription" "this" {
  name  = "${var.name}-${random_string.suffix.result}"
  topic = var.broker

  ack_deadline_seconds = var.ack_deadline_seconds

  filter = join(" AND ", local.filter-elements)

  cloud_storage_config {
    bucket = var.bucket
    filename_prefix = "${var.filter["type"]}/"
    max_bytes = 1000000000 // 1GB
    max_duration = "600s"
  }

  expiration_policy {
    ttl = "" // This does not expire.
  }

  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.dead-letter.id
    max_delivery_attempts = var.max_delivery_attempts
  }
}

// Grant the pubsub service account the ability to Acknowledge messages on this "this" subscription.
resource "google_pubsub_subscription_iam_binding" "allow-pubsub-to-ack" {
  subscription = google_pubsub_subscription.this.name

  role    = "roles/pubsub.subscriber"
  members = ["serviceAccount:${google_project_service_identity.pubsub.email}"]
}
