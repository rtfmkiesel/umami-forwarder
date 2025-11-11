<h1 align="center">umami-forwarder</h1>

<p align="center">
  A simple tool to use <a href="https://github.com/umami-software/umami">Umami</a> invisible.<br>Benefit from statistics on JavaScript-Free sites, or even if users block your tracker. Docker image only 2.33MB!
</p>

<p align="center">
  <a href="https://github.com/rtfmkiesel/umami-forwarder">
    <img src="https://img.shields.io/github/stars/rtfmkiesel/umami-forwarder" alt="Stars">
  </a>
  <a href="https://github.com/rtfmkiesel/umami-forwarder/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/rtfmkiesel/umami-forwarder" alt="LICENSE">
  </a>
  <a href="https://github.com/rtfmkiesel/umami-forwarder/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/rtfmkiesel/umami-forwarder/ghcr.yaml" alt="Build Status" />
  </a>
</p>

---

## How does this work?

`umami-forwarder` works by receiving mirrored/shadowed HTTP requests from a reverse proxy, parsing them, and sending them to the collection endpoint of an Umami instance. This requires a reverse proxy capable of mirroring/shadowing HTTP requests.

> [!NOTE]  
> This is very much beta software. Expect all the bugs.

## Setup

**One instance of `umami-forwarder` per site is required.**

### Environment Variables

`umami-forwarder` is configured through the following environment variables.

| Environment Variable | Description                                                  | Default Value |
|----------------------|--------------------------------------------------------------|---------------|
| `WEBSITE_ID`         | The website ID (from Umami dashboard)                        | -             |
| `COLLECTION_URL`     | The absolute URL to the Umami collection endpoint            | -             |
| `IP_HEADER`          | Which header contains the real client IP-address             | `X-Real-IP`   |
| `HTTP_TIMEOUT`       | HTTP timeout in seconds when connecting to Umami             | `5`           |
| `HTTP_RETRIES`       | HTTP retries when connecting to Umami                        | `3`           |
| `HTTP_MAX_REQUESTS`  | Limit on how many concurrent HTTP requests are made to Umami | `25`          |
| `HTTP_IGNORE_TLS`    | Ignore TLS errors when connecting to Umami                   | `false`       |
| `DEBUG`              | Enable debug output                                          | -             |

### Docker

This is meant to be used inside a Docker environment. See the `umami-forwarder` block from below or check out the [examples](./examples). 

```yaml
services:
  umami-forwarder:
    image: ghcr.io/rtfmkiesel/umami-forwarder
    environment:
      COLLECTION_URL: http://umami:3000/api/send
      WEBSITE_ID:  ca7f3ee6-e396-4bdb-858f-983978179794 
      IP_HEADER: X-Real-IP
    depends_on:
      umami:
        condition: service_healthy
    restart: unless-stopped
```

### Reverse Proxy

You need to configure your reverse proxy for mirroring. [Here](./examples/) are some examples.

## Contributing

Improvements in the form of PRs are welcome.

## Legal

This project is not affiliated with [Umami](https://github.com/umami-software/umami).