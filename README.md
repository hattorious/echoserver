# echoserver
smol server who talks back

| Port | Protocols                                                                          |
|------|------------------------------------------------------------------------------------|
| 8001 | http; supports http/1.0, http/1.1 |
| 8002 | h2c (clear-text http2); supports upgrade and prior knowledge |
| ~~8003 | http2 (h2c with TLS)~~ |
| ~~8004 | clear-text grpc~~ |
| ~~8005 | grpc with TLS~~ |
| ~~8006 | clear-text TCP~~ |
| ~~8007 | UDP~~ |

## build and run docker

```bash
make run
```