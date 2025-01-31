# Cloudflare DDNS updater

Use this Golang script for dynamicly updating your IP address for multiple Records. This only works when using Cloudflare for your Domain.

This script fetches your current public ip from cloudflare (https://cloudflare.com/cdn-cgi/trace) and updates all records in your configured zones (Record needs to have the ddns_comment).

## Config setup

Enter the following details in the config.yml file. It will be "hot reloaded".

If update_interval is set to 0 the process will only run once

```yaml
api_key: API_KEY
ddns_comment: AUTO_DDNS
update_interval: 5000
zones:
    - ZONE_ID
```

Todo:

-   Pointers
-   asynchronous code
