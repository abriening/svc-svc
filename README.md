# Service Service

## Build

Build the service:
```bash
make
# or
make all
# or
go build
```

## Build

Run the service on port 4444:
```bash
bin/svc-svc -port 4444
```

```bash
curl 'http://localhost:4444/ping '-v
curl 'http://localhost:4444/random'
```

## Run Chain

Run these in separate terminals:

```bash
bin/svc-svc -port 4440
bin/svc-svc -port 4441 -upstream 'http://localhost:4440/random'
bin/svc-svc -port 4442 -upstream 'http://localhost:4441/upstream?name=level1'
bin/svc-svc -port 4443 -upstream 'http://localhost:4442/upstream?name=level2'
bin/svc-svc -port 4444 -upstream 'http://localhost:4443/upstream?name=level3'
```

```bash
curl 'http://localhost:4444/upstream?name=level4'
curl 'http://localhost:4444/ping' -v
curl 'http://localhost:4443/ping' -v
curl 'http://localhost:4442/ping' -v
curl 'http://localhost:4441/ping' -v
curl 'http://localhost:4440/random'
```
