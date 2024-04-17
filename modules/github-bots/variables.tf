variable "extra_filter" {
  type        = map(string)
  default     = {}
  description = "Optional additional filters to include."
}

variable "extra_filter_prefix" {
  type        = map(string)
  default     = {}
  description = "Optional additional prefixes for filtering events."
}

variable "extra_filter_has_attributes" {
  type        = list(string)
  default     = []
  description = "Optional additional attributes to check for presence."
}

variable "extra_filter_not_has_attributes" {
  type        = list(string)
  default     = []
  description = "Optional additional prefixes to check for presence."
}

variable "containers" {
  description = "The containers to run in the service.  Each container will be run in each region."
  type = map(object({
    source = object({
      base_image  = optional(string, "cgr.dev/chainguard/static:latest-glibc")
      working_dir = string
      importpath  = string
    })
    args = optional(list(string), [])
    ports = optional(list(object({
      name           = optional(string, "http1")
      container_port = number
    })), [])
    resources = optional(
      object(
        {
          limits = optional(object(
            {
              cpu    = string
              memory = string
            }
          ), null)
          cpu_idle          = optional(bool, true)
          startup_cpu_boost = optional(bool, false)
        }
      ),
      {
        cpu_idle = true
      }
    )
    env = optional(list(object({
      name  = string
      value = optional(string)
      value_source = optional(object({
        secret_key_ref = object({
          secret  = string
          version = string
        })
      }), null)
    })), [])
    regional-env = optional(list(object({
      name  = string
      value = map(string)
    })), [])
    volume_mounts = optional(list(object({
      name       = string
      mount_path = string
    })), [])
  }))
}
