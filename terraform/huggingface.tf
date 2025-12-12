module "huggingface_sagemaker" {
  source               = "philschmid/sagemaker-huggingface/aws"
  name_prefix          = "deploy-hub"
  pytorch_version      = "1.9.1"
  transformers_version = "4.12.3"
  instance_type        = "cpu"
  hf_model_id          = "deepset/deberta-v3-base-injection"
  hf_task              = "text-classification"
  serverless_config    = {
    max_concurrency   = 3
    memory_size_in_mb = 3072
  }
}