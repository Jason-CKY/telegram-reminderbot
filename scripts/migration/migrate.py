import json
import requests
from pydantic import BaseModel
from pathlib import Path

# ["file_id":"AgACAgQAAxkBAAIek2Is_EZgR_T1eb8pTmgJ11o5h2qFAAIcuTEbQq9oUQ5b14L7haxeAQADAgADeAADIwQ","time":"14:00","frequency":"Weekly-1","reminder_id":"f4ead1aa-2472-48ac-8753-ac496201f7c5","job_id":"5ee9e73f-f454-4204-8fc4-e3b6a35beeaa","timezone":"Asia/Singapore"},{"user_id":5226530735,"reminder_text":"To be greatful for you god","time":"13:30","frequency":"Daily","reminder_id":"2cbed323-ad1a-4ed9-b799-47eb03119504","job_id":"466284f7-22b6-4948-9634-1d9219200386","timezone":"Asia/Singapore"}]},{"_id":{"$oid":"622e2bc26e23be1a2adb2bf5"},"chat_id":1374999435,"timezone":"Asia/Singapore","update_settings":false,"reminders_in_construction":[{"user_id":1374999435,"reminder_text":"Hey"}],"reminders":[{"user_id":1374999435,"reminder_text":"Aşk bezi what she was about to tell you","time":"23:00","frequency":"Once 2022-03-16","reminder_id":"e14b840b-d392-4597-b2a2-ada20456bf20","job_id":"d2e1dea1-c74f-4e95-adbf-f68e7f8f8fd5","timezone":"Asia/Singapore"}]},{"_id":{"$oid":"622efeb86e23be1a2adb2f8b"},"chat_id":84962687,"timezone":"Asia/Singapore","update_settings":true,"reminders_in_construction":[],"reminders":[{"user_id":84962687,"reminder_text":"","file_id":"AgACAgQAAxkBAAIe_mIu_1vY1g60McGD5up4g_xyu8KgAAI7uDEb8vh5UZB4Pe2_g2fvAQADAgADeQADIwQ","time":"02:00","frequency":"Weekly-6","reminder_id":"4fb3393f-a1f]
class ChatRow(BaseModel):
    reminder_id: str
    user_id: int
    reminder_text: str
    file_id: str | None
    time: str
    frequency: str
    timezone: str


def get_all_chat_ids(chat_collection):
    chat_ids = []
    for chat in chat_collection:
        if chat["chat_id"] not in chat_ids:
            chat_ids.append(chat["chat_id"])
        else:
            print(chat)
    return chat_ids

def insert_directus_chat_settings(chat_collection_row):
    pass


def main():
    with open(Path(__file__).parent / 'chat_collection.json', "r") as f:
        chat_collection = json.load(f)
    chat_ids = get_all_chat_ids(chat_collection)
    # print(len(chat_ids))
    # print(json.dumps(chat_collection[0], indent=True))

if __name__ == '__main__': 
    main()