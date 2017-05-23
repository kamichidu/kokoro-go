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

```
{"jsonrpc":"2.0","method":"websocket","params":{"command":"subscribe","identifier":"{\"channel\":\"ChatChannel\"}"}}
{"jsonrpc":"2.0","method":"websocket","params":{"command":"message","identifier":"{\"channel\":\"ChatChannel\"}","data":"{\"access_token\":\"...AccessToken...\",\"rooms\":[\"random\"],\"action\":\"subscribe\"}"}}
... ActionCable message will be shown like below format:
... {"jsonrpc":"2.0","id":"__websocket__","result":{"identifier":"{\"channel\":\"ChatChannel\"}","message":{"content":"\u003cp\u003ehi\u003c/p\u003e\n","id":781,"profile":{"avatar":"https://d23u46pnxyg3wa.cloudfront.net/attachments/a1a916b82454cdd405a3153dd3319ad00192b048/store/fill/40/40/c96e5eab967f470e0e35884690856d267117007d68db773124f192e5f1fc/avatar.png","display_name":"supermomonga","screen_name":"supermomonga","type":"user"},"published_at":"2017-05-23T02:20:08Z","raw_contet":"hi","room":{"id":3,"screen_name":"random"}}}}
```
