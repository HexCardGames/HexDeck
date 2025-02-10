# 🎮 HexDeck *(WIP)*

[![Project Status](https://img.shields.io/badge/status-WIP-orange)]()
[![Docker Image](https://img.shields.io/badge/Docker-Coming_Soon-lightgrey?style=flat&logo=docker)]()
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

An open-source, self-hostable online multiplayer card game, built with a robust Go backend and a Svelte frontend. Developed as an Abitur project by two students in Berlin.

## 🕹️ Gameplay

HexDeck is a fast-paced card game where players match colors or numbers to empty their hand first with some unique action cards. Strategic card play and a bit of luck are key to victory! The name HexDeck comes from the idea of not only using numbers 0-9 but also including (optional) hexadecimal numbers (A-F) for additional gameplay depth.

## 🚀 Project Status

This project is currently in the early development phase. The repository does not yet contain a working codebase, but the following features are planned:

- 🎴 **Multiplayer Card Game**: A fast-paced game where players match colors or numbers to empty their hand first.
- 🖥️ **Self-Hostable**: Run your own server using Docker for easy deployment.
- 🌐 **Free Public Server**: A free public instance of HexDeck.
- 🌍 **Online Multiplayer**: Play with friends or other players over the internet.

## 🏗️ Planned Features

- **Game Lobby**: Create and join multiplayer rooms.
- **Gameplay itself**: Somewhat fair and strategic game mechanics.
- **Customizable Rules**: Flexible game settings with optional card decks.
- **Docker Deployment**: Easy setup with Docker and Docker Compose.
- **Public Server**: A free-to-play hosted version.

## 🐳 Docker Setup *(Planned)*

Once development progresses, we aim to provide an easy way to deploy HexDeck using Docker:

```bash
# Planned Docker Compose setup
# docker compose up
```

```bash
# Planned standalone Docker setup
# docker run -d --name hexdeck-server -p 3000:80 -v ./data:./data unterdrueckt/hexdeck-server:latest
```

---

## 🤝 Contributing

Contributions are welcome! Please open an issue or a pull request to improve the library.

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
