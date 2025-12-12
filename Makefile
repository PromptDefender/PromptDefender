MODULES := $(shell (find .  -type f -name '*.go' -maxdepth 4 | sed -r 's|/[^/]+$$||' |cut -c 3-|sort |uniq))
TEST_MODULES := $(shell (find . -type f -name '*.go' -maxdepth 3 ! -path './test/integration_test_harness/*' | grep "_test" | sed -r 's|/[^/]+$$||' | sort | uniq))
PROJECT_DIR := $(shell pwd)
API_DIR := $(shell pwd)/api

setup-workspace:
	if [ -n "$$GITHUB_REF_NAME" ]; then \
		echo "Using branch name from GITHUB_REF_NAME env variable..." &&\
        export TF_VAR_branch_name=$$GITHUB_REF_NAME; \
	else \
		echo "Using branch name from git rev-parse..." &&\
		export TF_VAR_branch_name=$$(git rev-parse --abbrev-ref HEAD); \
	fi; \
	if [ "$$TF_VAR_branch_name" = "main" ] && [ "$$INTEGRATION_TEST" != "true" ]; then \
		echo "On 'main' branch. Using the 'default' workspace..."; \
		cd terraform/gcp && terraform init && terraform workspace select -or-create default || exit 1; \
		echo "Workspace $$TF_VAR_branch_name selected."; \
		terraform workspace show; \
		cd ../..; \
	else \
		echo "Workspace $$TF_VAR_branch_name exists. Selecting it..."; \
		workspace_name=`echo $$TF_VAR_branch_name | sed 's/[^a-zA-Z0-9-]/-/g' | cut -c 1-20` ; \
		cd terraform/gcp && terraform init && terraform workspace select -or-create  $$workspace_name; \
		echo "Workspace $$TF_VAR_branch_name selected."; \
		terraform workspace show; \
		cd ..; \
	fi

test: build
	for testable_module in $(TEST_MODULES) ; do \
	   cd $$testable_module && go test -v ./... -cover || exit 1; cd $(PROJECT_DIR) ; \
	done

build: generate
	mkdir -p $(PROJECT_DIR)/bin
	cd cmd/server && go build -o $(PROJECT_DIR)/bin/server

deploy: setup-workspace build
	export TF_VAR_commit_version=`git rev-parse --short HEAD` &&\
	cd terraform/gcp && terraform init && terraform apply -auto-approve &&\
	terraform output -json > ../../terraform_output.json

install:
	for number in  $(MODULES) ; do \
       cd $$number && go get ./... || exit 1; cd .. ; \
    done
	cd cmd/server && go get ./...

tidy:
	for number in $(MODULES); do \
		cd $$number && go mod tidy || exit 1; cd $(PROJECT_DIR) ; \
	done
	cd cmd/server && go mod tidy

upgrade:
	for number in  $(MODULES) ; do \
		printf "Upgrading dependencies for module: %s\n" $$number; \
	   cd $$number && go get -u all  || exit 1; cd .. ; \
	done
	cd cmd/server && go get -u all

clean:
	for number in  $(MODULES) ; do \
	   cd $$number && go clean -testcache || exit 1; cd .. ; \
	done
	cd test/integration_test_harness && go clean -testcache || exit 1; cd $(PROJECT_DIR) ;
	cd cmd/server && go clean -testcache

generate:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	oapi-codegen -package integration_test_harness -generate types,client $(API_DIR)/openapi.yml > test/integration_test_harness/api.gen.go
	oapi-codegen -package main -generate types,chi-server,spec -o cmd/server/api.gen.go api/openapi.yml

generate_jailbreak:
	cd builder\
	 && pip install -r requirements.txt && python3 clean_jailbreaks_into_json.py\
	 && python3 jailbreak_embeddings.py && go build -o main && ./main

integration_test:
	go install github.com/tomwright/dasel/cmd/dasel@latest
	export URL=`cd terraform/gcp && terraform output -json | dasel select -p json '.api_url.value' | tr -d '"'` &&\
	export DEFENDER_API_KEY=`cd terraform/gcp && terraform output -json | dasel select -p json '.api_key_value.value' | tr -d '"'` &&\
	echo "Defender API URL: $$URL" &&\
	cd test/integration_test_harness && go test -count=1 -v ./...

destroy: setup-workspace
	export TF_VAR_commit_version=`git rev-parse --short HEAD`;\
	current_workspace=`cd terraform/gcp && terraform workspace show`;\
	if [ "$$current_workspace" = "default" ]; then \
		echo "Skipping destruction in default workspace"; \
	else \
		cd terraform/gcp && terraform init && terraform destroy -auto-approve; \
	fi

load_test:
	export URL=`cd terraform/gcp && terraform output -json | dasel select -p json '.api_url.value' | tr -d '"'` &&\
	export DEFENDER_API_KEY=`cd terraform/gcp && terraform output -json | dasel select -p json '.api_key_value.value' | tr -d '"'` &&\
	cd test/load && k6 run wall_load.js