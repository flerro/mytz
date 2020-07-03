.PHONY: build

deployEnv = --region=eu-west-1 --profile personal
bucket = fl-sam-deployments
stackName = mytz

build:
	sam build

run: build
	sam local invoke

serve: build
	sam local start-api

package:
	sam package --s3-bucket ${bucket} --output-template-file /tmp/${stackName}-out.yaml ${deployEnv}

deploy: package build
	sam deploy --template-file /tmp/${stackName}-out.yaml --capabilities CAPABILITY_IAM --stack-name mytz --s3-bucket ${bucket} ${deployEnv}