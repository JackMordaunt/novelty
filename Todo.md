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

Player handles play, pause, seek and any other playback functionality. 

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

## Backlog
- [ ] Watch button
    - Select quality
    - Select player
    - [x] Send "show.open" message
- [ ] Loading screen
    - Show name
    - Download speed
    - Progress bar where 100% = playable
- [ ] Bookmarks
    - Bookmarks page (displays all bookmarked shows)
    - Hash set of Show ID's
    - Store in json file
    - "bookmarks.add" and "bookmarks.remove" messagses with corresponding events

### How do we transition from watch button clicked to loading screen? 
Need an event listener that listens for "player.opened" and navigates to the 
loading component/page. The loading page shows the details: show name,
progress until playable, download speed, quality. Once the show is ready for
playback we open the configured player, this could be some global event handler 
or part of the loading component. 

### Loading.vue

- Mount when resource is opened.
- Display show status updates (throughput, progress, health, etc).
- UnMount when resource is closed.