Watchtower provides an HTTP API mode that enables an HTTP endpoint that can be requested to trigger container updating. The current available endpoint list is:

-   `/v1/update` - triggers an update for all of the containers monitored by this Watchtower instance.

---

To enable this mode, use the flag `--http-api-update`. For example, in a Docker Compose config file:

```yaml
version: '3'

services:
  app-monitored-by-watchtower:
    image: myapps/monitored-by-watchtower
    labels:
      - "com.centurylinklabs.watchtower.enable=true"

  watchtower:
    image: marrrrrrrrry/watchtower
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: --debug --http-api-update
    environment:
      - WATCHTOWER_HTTP_API_TOKEN=mytoken
    labels:
      - "com.centurylinklabs.watchtower.enable=false"
    ports:
      - 8080:8080
```

By default, enabling this mode prevents periodic polls (i.e. what is specified using `--interval` or `--schedule`). To run periodic updates regardless, pass `--http-api-periodic-polls`.

Notice that there is an environment variable named WATCHTOWER_HTTP_API_TOKEN. To prevent external services from accidentally triggering image updates, all of the requests have to contain a "Token" field, valued as the token defined in WATCHTOWER_HTTP_API_TOKEN, in their headers. In this case, there is a port bind to the host machine, allowing to request localhost:8080 to reach Watchtower. The following `curl` command would trigger an image update:

```bash
curl -H "Authorization: Bearer mytoken" localhost:8080/v1/update
```

---

### Changing the listen address

By default, the HTTP API listens on `:8080`. Use the `--http-api-listen-address` flag (or the `WATCHTOWER_HTTP_API_LISTEN_ADDRESS` environment variable) to change the address and port.

The value accepts the following forms, which control which interfaces the API is reachable on:

```text
127.0.0.1:9789   only reachable from the local machine
192.168.1.2:9789 reachable only on that specific interface
:9789            reachable on all interfaces
9789             same as :9789 (a bare port is prefixed with a colon)
0.0.0.0:9789     same as :9789 (an explicit form of all interfaces)
```

Changing the listen address is required when running with host networking, since port mappings such as `-p 8080:8080` have no effect in `--network host` mode and the API would otherwise always bind to port 8080 on the host:

```yaml
services:
  watchtower:
    image: marrrrrrrrry/watchtower
    network_mode: host
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: --http-api-update
    environment:
      - WATCHTOWER_HTTP_API_TOKEN=mytoken
      - WATCHTOWER_HTTP_API_LISTEN_ADDRESS=127.0.0.1:9789
```

The example above also binds the API to the loopback interface (`127.0.0.1`), so that the endpoint is reachable only from the local machine instead of every interface of the host. With the API listening on `127.0.0.1:9789`, the following `curl` command would trigger an image update:

```bash
curl -H "Authorization: Bearer mytoken" localhost:9789/v1/update
```

---

In order to update only certain images, the image names can be provided as URL query parameters. The following `curl` command would trigger an update for the images `foo/bar` and `foo/baz`:

```bash
curl -H "Authorization: Bearer mytoken" localhost:8080/v1/update?image=foo/bar,foo/baz
```
