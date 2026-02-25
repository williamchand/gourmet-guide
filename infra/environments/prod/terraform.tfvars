project_id             = "ameno-discovery"
environment            = "prod"
region                 = "asia-southeast2"
backend_image          = "docker.io/williamchand/gourmet-guide-backend:prod"
gemini_model           = "gemini-2.5-flash-native-audio-preview-12-2025"
cloud_run_ingress      = "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER"
allow_unauthenticated  = false
vpc_connector          = "projects/ameno-discovery/locations/asia-southeast2/connectors/prod-serverless"
