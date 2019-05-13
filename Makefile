.PHONY: image test

NAME=tweetserver
USERNAME=ozapinq
TAG?=latest

IMAGE_NAME=${USERNAME}/${NAME}:${TAG}

image:
	docker build -t ${IMAGE_NAME} .

push:
	docker push ${IMAGE_NAME}
