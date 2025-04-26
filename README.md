### notion_to_pushover

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
      - PUSHOVER_NOTIFICATION_TITLE=<PUSHOVER_NOTIFICATION_TITLE> # Default = A new page has been added
      - NOTION_API_KEY=<NOTION_API_KEY>
      - NOTION_VERIFICATION_TOKEN=<NOTION_VERIFICATION_TOKEN>
    ports:
      - 8069:8069
    restart: unless-stopped
```
| Variable                     | Description                                                                             |Required |
|------------------------------|-----------------------------------------------------------------------------------------|---------|
| `PUSHOVER_USER_KEY`          | The key associated with your Pushover account.                                          | ✅      |
| `PUSHOVER_API_TOKEN`         | The token associated with your Pushover application.                                    | ✅      |
| `PUSHOVER_NOTIFICATION_TITLE`| What the title of the notification should be (defaults to `A new page has been added`.  | ❌      |
| `NOTION_API_KEY`             | The API token associated with your Notion integration.                                  | ✅      |
| `NOTION_VERIFICATION_TOKEN`  | The verification token sent when setting up your WebHook in Notion.                     | ✅      |
