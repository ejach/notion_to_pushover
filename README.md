### notion_to_pushover
[![PyPI](https://img.shields.io/pypi/v/notion-client?logo=python&label=notion-client&style=flat-square&color=FFD43B)](https://pypi.org/project/notion-client/)
[![PyPI](https://img.shields.io/pypi/v/Requests?logo=python&label=Requests&style=flat-square&color=FFD43B)](https://pypi.org/project/Requests/)

Uses Notion WebHooks to send notifications using Pushover

#### docker-compose.yml
```yml
version: '3.7'
services:
  notion_to_pushover:
    image: ghcr.io/ejach/notion_to_pushover
    container_name: notion_to_pushover
    environment:
      - PUSHOVER_USER_KEY=<PUSHOVER_USER_KEY>
      - PUSHOVER_API_TOKEN=<PUSHOVER_API_TOKEN>
      - NOTION_API_KEY=<NOTION_API_KEY>
      - NOTION_VERIFICATION_TOKEN=<NOTION_VERIFICATION_TOKEN>
      - STRICT_MODE=<True or False> # Default = False
    ports:
      - 8069:8069
    restart: unless-stopped
```
| Variable                    | Description                                                                      | Required |
|-----------------------------|----------------------------------------------------------------------------------|--------|
| `PUSHOVER_USER_KEY`         | The key associated with your Pushover account.                                   | ✅      |
| `PUSHOVER_API_TOKEN`        | The token associated with your Pushover application.                             | ✅      |
| `NOTION_API_KEY`            | The API token associated with your Notion integration.                           | ✅      |
| `NOTION_VERIFICATION_TOKEN` | The verification token sent when setting up your WebHook in Notion.              | ✅      |
| `STRICT_MODE`               | Whether or not to verify requests using the `X-Notion-Signature` (Default False) | ❌       |
