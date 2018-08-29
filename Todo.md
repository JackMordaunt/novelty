# Novelty

## Blur Effects

MSHTML doesn't support css blur filter. One workaround for this would be to blue the image in Go and simply expose an api that will return the blurred image.

## Websocket

- [ ] Create generalised Read/Write websocket client.
    - Client takes a router to outsource the handling of the each incoming message.
- [ ] Websocket router to route incoming messages to handlers.
- [ ] Use websockets exclusively for communication (rather than mixing standard http and websockets).

## Commands

- Open Media

## Events

- Media Opened
- Media Status Update

## Core API

- Open a Show.
    - Starts download via the appropriate protocol based on the Show metadata
    - Reports download status.
    - Provides a seekable reader for the currently downloaded Show Resource.

Open(Show): Resource
    - Resource is a closeable and seekable reader.
    - Show is metadata


## "seeker can't seek" error.

- Is the torrent client "ready" for playback? 5% is downloaded.
- Is the _largest file_ being served?
- Is the http.ServeContent picking up the correct file type? 
    - The response says content-type is text/plain, why? 

```
    Content-Disposition: attachment; filename="big-buck-bunny.mp4"
    Content-Type: text/plain; charset=utf-8
    Last-Modified: Wed, 29 Aug 2018 06:00:17 GMT
    X-Content-Type-Options: nosniff
    Date: Wed, 29 Aug 2018 06:00:17 GMT
    Content-Length: 18
```