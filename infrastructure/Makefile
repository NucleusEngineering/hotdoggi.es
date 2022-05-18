init:
	gsutil mb gs://hotdoggi-es-terraform-state/ || exit 0
	terraform init

apply:
	terraform apply

destroy:
	terraform destroy

.PHONY: init apply destroy
