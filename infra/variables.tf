variable "project_id" {
  description = "Google Cloud project ID."
  type        = string
}

variable "environment" {
  description = "Deployment environment name (dev, staging, prod)."
  type        = string
}

variable "region" {
  description = "Google Cloud region."
  type        = string
  default     = "asia-southeast2"
}

variable "backend_image" {
  description = "Container image URI for backend service."
  type        = string
}

variable "gemini_model" {
  description = "Gemini model name used for runtime responses."
  type        = string
  default     = "gemini-2.5-flash-native-audio-preview-12-2025"
}

variable "gemini_model_2" {
  description = "Gemini model name used for runtime responses."
  type        = string
  default     = "gemini-2.0-flash"
}

variable "gemini_model_3" {
  description = "Gemini model name used for runtime responses."
  type        = string
  default     = "gemini-2.0-flash-preview-image-generation"
}

variable "cloud_run_ingress" {
  description = "Ingress mode for Cloud Run service."
  type        = string
  default     = "INGRESS_TRAFFIC_ALL"
}

variable "allow_unauthenticated" {
  description = "Allow unauthenticated invocation to backend service."
  type        = bool
  default     = false
}

variable "vpc_connector" {
  description = "Optional VPC connector name for Cloud Run egress."
  type        = string
  default     = ""
}
