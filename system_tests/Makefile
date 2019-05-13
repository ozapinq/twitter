IMAGE_NAME=ozapinq/tweetserver-system-tests

image:
	docker build -t ${IMAGE_NAME}  .

push:
	docker push ${IMAGE_NAME}

run: image
	docker run -it -e "URL=${URL}" --rm ${IMAGE_NAME}
