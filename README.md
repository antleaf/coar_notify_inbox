# coar_notify_inbox
COAR Notify LDN inbox and validation test system

# Antleaf LDN Inbox


## Build Image
`docker build -t antleaf_ldn_inbox .`

## Publish Image

```
docker image tag antleaf_ldn_inbox:latest antleaf/antleaf_ldn_inbox:1.2.1
docker login 
docker push antleaf/antleaf_ldn_inbox:1.2.1
```

## Alternativley, building and pushing from Apple Silicon:
```bash
docker buildx build --platform linux/amd64,linux/arm64 --push -t antleaf/antleaf_ldn_inbox:1.2 .

```

## Run container
````bash
docker run \
	-it \
	--rm \
	--name antleaf_ldn_inbox_instance \
	-p 1313:80 \
	antleaf_ldn_inbox
```