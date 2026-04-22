# Security Policy

## Supported Versions

| Version | Supported |
| ------- | --------- |
| 0.1.x   | Yes       |

## Reporting a Vulnerability

If you discover a security vulnerability, please **do not** open a public issue. Instead, please open a [private security advisory](https://github.com/shing1211/futuapi4go-demo/security/advisories/new) or contact the maintainer directly.

We ask that you:

- Give us reasonable time to investigate and fix the issue before public disclosure.
- Provide enough detail to reproduce the vulnerability.
- Do not access or modify data that does not belong to you.

We will acknowledge your report within 48 hours and aim to provide a fix within 7 days for confirmed vulnerabilities.

## General Security Notes

- This project connects to Futu OpenD via a local TCP socket. Ensure your OpenD instance is not exposed to untrusted networks.
- Trading operations in the demo require an unlocked trading account. **Never run the demo on a production account without proper safeguards.**
- Do not share your Futu API credentials or trade password in any public forum.
- Do not commit `.env` files or any files containing credentials to the repository.
