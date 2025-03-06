# Docker 이미지 빌드
client:
	docker build --platform linux/amd64 -t stradvision-client:latest -f Dockerfile.client .

consumer:
	docker build --platform linux/amd64 -t stradvision-consumer:latest -f Dockerfile.consumer .

recovery:
	docker build --platform linux/amd64 -t stradvision-recovery:latest -f Dockerfile.recovery .


build: client consumer recovery
	echo "Build Done"