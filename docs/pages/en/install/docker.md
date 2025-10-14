---
icon: material/docker
title: Installation from Docker image
---

# :material-docker: Installation from Docker image

## :material-archive-outline: Registries

LazyGate provides official Docker images hosted on multiple registries:

| Registry                                                                                            | Image                        |
|-----------------------------------------------------------------------------------------------------|------------------------------|
| [Docker Hub](https://hub.docker.com/r/kasefuchs/lazygate/)                                          | `kasefuchs/lazygate`         |
| [GitHub Container Registry](https://github.com/users/kasefuchs/packages/container/package/lazygate) | `ghcr.io/kasefuchs/lazygate` |

## :material-rocket-launch: Usage

=== "Docker Compose"

    ```yaml
    ---
    services:
      gate:
        image: kasefuchs/lazygate:latest
        ports:
          - "25565:25565"
        volumes:
          - /var/run/docker.sock:/var/run/docker.sock:rw
          - ./config.yml:/config.yml:ro
        environment:
          LAZYGATE_PLUGIN_PROVIDER: docker
          LAZYGATE_PLUGIN_NAMESPACE: default
        restart: always
    ```

=== "Docker CLI"

    ```sh
    docker run -d --name lazygate \
      -p 25565:25565 \
      -v /var/run/docker.sock:/var/run/docker.sock:rw
      -v ./config.yml:/config.yml:ro \
      -e LAZYGATE_PLUGIN_PROVIDER=docker \
      -e LAZYGATE_PLUGIN_NAMESPACE=default \
      --restart always \
      kasefuchs/lazygate:latest
    ```
