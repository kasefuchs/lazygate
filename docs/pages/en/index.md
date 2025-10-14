---
title: Home
hide:
  - toc
  - navigation
---

# :material-home: Home

## :material-magnify: What is LazyGate?

**LazyGate** is a [Gate proxy](https://github.com/minekube/gate) plugin that automatically shuts down your Minecraft server when it’s idle and starts it again when players connect.

Some Minecraft servers, especially modded ones, use a huge amount of resources even when nobody is playing. LazyGate helps save resources by stopping the backend server during idle periods and waking it up automatically when a player tries to join.

As a Gate plugin, LazyGate runs inside the proxy and intercepts incoming connection attempts. When the backend server is offline, it starts it up and then seamlessly transfers the player once it’s ready, all without them noticing.

## :material-feature-search-outline: Features

- Runs entirely as a [Gate proxy](https://github.com/minekube/gate) plugin
- Automatically stops the backend Minecraft server when idle and starts it when players connect
- Supports multiple server providers out of the box: **Docker**, **Nomad**, and **PufferPanel**
- Provides multiple built-in player queues: **wait** and **kick**
- Plugin system for custom providers and queue implementations

## :material-license: License

This project is licensed under the terms of the MIT license.
