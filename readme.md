# Cloudflare DDNS updater

Use this Python script for dynamicly updating your IP address for multiple Records. This only works when using Cloudflare for your Domain.

This script fetches your current public ip from cloudflare (https://cloudflare.com/cdn-cgi/trace) and updates all defined type "A" records in your Cloudflare DNS settings.

## Config setup

Enter the following details in the config.json file.

```json
{
  "token": "YOUR_CLOUDFLARE_API_TOKEN",
  "zone": "CLOUDFLARE_DOMAIN_ZONE_ID",
  "skipUpdate": "COMMENT TO IGNORE THIS RECORD"
}
```

# Crontab Setup for automatic execution

Add this line to your crontab file (run crontab -e) for running the script every 10 minutes.

```
*/10    *       *       *       *       python3 /your/path/ddns.py
```
