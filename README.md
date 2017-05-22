kokoro-go
========================================================================================================================
[kokoro-io](https://github.com/supermomonga/kokoro-io) REST API or WebSocket API client for cli tools.

How to Install
------------------------------------------------------------------------------------------------------------------------
```
go get -u -v github.com/kamichidu/kokoro-go
```

How to Use
------------------------------------------------------------------------------------------------------------------------
```
# Connect default endpoint https://kokoro.io and wss://kokoro.io with {AccessToken}
$ kokoro-go {AccessToken}
{"jsonrpc":"2.0","id":"hoge","method":"http","params":{"type":"POST","url":"/api/v1/rooms","data":{"limit":100}}}
... response in line ...
```
