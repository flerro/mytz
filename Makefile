.PHONY: build

build:
	sam build

serve: build
	sam local start-api

package:
	sam package --s3-bucket fl-sam-deployments --output-template-file out.yaml --region=eu-west-1 --profile personal

deploy: package build
	sam deploy --template-file out.yaml --capabilities CAPABILITY_IAM --stack-name mytz --s3-bucket fl-sam-deployments --region=eu-west-1 --profile personal