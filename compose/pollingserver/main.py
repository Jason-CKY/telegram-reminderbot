import time
import logging
import os
import requests
import json

logging.basicConfig(encoding='utf-8', level=logging.INFO)

def main():
    OFFSET = 0
    bot_token = os.environ.get("BOT_TOKEN")
    polling_interval = float(os.environ.get("POLLING_INTERVAL", "0.2"))
    while True:
        logging.info("Checking for updates...")
        updates = requests.get(f"https://api.telegram.org/bot{bot_token}/getUpdates",
                                params={"offset": OFFSET}
                            )
        if updates.status_code == 200:
            updates = updates.json()['result']
            if len(updates) > 0:
                print(json.dumps(updates, indent=4))
            for update in updates:
                requests.post(f'http://{os.environ.get("APP_SERVER")}/{bot_token}', json=update)
            OFFSET = max([update['update_id'] for update in updates]) + 1 if len(updates) > 0 else OFFSET

        time.sleep(polling_interval)

if __name__ == '__main__':
    main()