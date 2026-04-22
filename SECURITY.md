# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability, please open an [issue](https://github.com/shing1211/futuapi4go-demo/issues) with the label `security`.

## General Security Notes

- This project connects to Futu OpenD via a local TCP socket. Ensure your OpenD instance is not exposed to untrusted networks.
- Trading operations in the demo require an unlocked trading account. Never run the demo on a production account without proper safeguards.
- Do not share your Futu API credentials or trade password in any public forum.
