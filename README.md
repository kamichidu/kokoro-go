kokoro-go
========================================================================================================================
[kokoro-io](https://github.com/supermomonga/kokoro-io) REST API or WebSocket API client for cli tools.

How to Install
------------------------------------------------------------------------------------------------------------------------
```
go get -u -v github.com/kamichidu/kokoro-go/...
```

How to Use
------------------------------------------------------------------------------------------------------------------------
```
# Configure access_token
$ kokoro-go config token {accessToken}

# Get channel ids
$ kokoro-go request get /api/v1/channels --query '[].id'

# Subscribe events for specific channels
$ kokoro-go websocket {channelId1} {channelId2}
```
