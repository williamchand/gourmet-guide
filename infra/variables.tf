variable "project_id" {
  description = "Google Cloud project ID."
  type        = string
}

variable "region" {
  description = "Google Cloud region."
  type        = string
  default     = "us-central1"
}

variable "backend_image" {
  description = "Container image URI for backend service."
  type        = string
}

variable "gemini_model" {
  description = "Gemini model name used for runtime responses."
  type        = string
  default     = "gemini-2.0-flash-live-001"
}
