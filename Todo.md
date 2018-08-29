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
