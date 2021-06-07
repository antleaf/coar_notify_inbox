# COAR LDN Inbox
COAR Notify LDN inbox and validation test system

## Build Image
```bash
docker build -t notify_ldn_inbox .
```

## Publish Image

```bash
docker image tag notify_ldn_inbox:latest antleaf/notify_ldn_inbox:1.0
docker login 
docker push antleaf/notify_ldn_inbox:1.0
```

## Run container

With defaults:
```bash
docker run \
	-it \
	--rm \
	-p 80:80 \
	antleaf/notify_ldn_inbox:1.0
```

Specifying arguments:
```bash
docker run \
	-it \
	--rm \
	-p 80:80 \
	antleaf/notify_ldn_inbox:1.0 \
	notify_ldn_inbox -db=/opt/data/ldn_inbox.sqlite -host=http://localhost -port=1313 -debug=true
```