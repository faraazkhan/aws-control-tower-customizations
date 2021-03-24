export TEMPLATE_OUTPUT_BUCKET = cms-cloud-custom-control-tower-template
export DIST_OUTPUT_BUCKET = cms-cloud-custom-control-tower-dist
export SOLUTION_NAME = customizations-for-aws-control-tower
export AWS_REGION = us-east-1
export AWS_DEFAULT_REGION = us-east-1

.PHONY: all test build publish

all: build publish

test:
	./deployment/run-unit-tests.sh

build: test
	./deployment/build-s3-dist.sh ${DIST_OUTPUT_BUCKET} ${TEMPLATE_OUTPUT_BUCKET} ${SOLUTION_NAME} ${VERSION}

publish:
	aws s3 cp deployment/global-s3-assets/ s3://${TEMPLATE_OUTPUT_BUCKET}/${SOLUTION_NAME}/${VERSION}/ --recursive --acl bucket-owner-full-control
	aws s3 cp deployment/global-s3-assets/ s3://${DIST_OUTPUT_BUCKET}-us-east-1/${SOLUTION_NAME}/${VERSION}/ --recursive --acl bucket-owner-full-control
