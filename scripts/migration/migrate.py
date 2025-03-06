import json
import requests
import os
from datetime import datetime, time, timedelta
from dateutils import relativedelta
import pytz
from pydantic import BaseModel
from pathlib import Path
from tqdm import tqdm

DIRECTUS_ACCESS_TOKEN = os.environ.get("DIRECTUS_ACCESS_TOKEN", "test-token")
DIRECTUS_URL = os.environ.get("DIRECTUS_URL", "http://localhost:8055")
DIRECTUS_PUBLIC_URL = os.environ.get("DIRECTUS_PUBLIC_URL", "http://localhost:8055")
REMINDER_ONCE = 'Once'
REMINDER_DAILY = 'Daily'
REMINDER_WEEKLY = 'Weekly'
REMINDER_MONTHLY = 'Monthly'
REMINDER_YEARLY = 'Yearly'

request_headers = {
    "Authorization": DIRECTUS_ACCESS_TOKEN,
    "Host": DIRECTUS_URL,
}

class ChatSettings(BaseModel):
    chat_id: int
    timezone: str
    updating: bool

class Reminder(BaseModel):
    from_user_id: int
    file_id: str
    reminder_text: str
    frequency: str
    time: str
    in_construction: bool
    next_trigger_time: str | None
    chat_id: int

def get_all_chat_ids(chat_collection):
    chat_ids = []
    for chat in chat_collection:
        if chat["chat_id"] not in chat_ids:
            chat_ids.append(chat["chat_id"])
        else:
            print(chat)
    return chat_ids

def insert_directus_chat_settings(chat_settings: ChatSettings):
    resp = requests.post(
        f"{DIRECTUS_PUBLIC_URL}/items/reminderbot_chat_settings",
        headers=request_headers,
        json=chat_settings.model_dump()
    )
    if resp.status_code != 200 and not (resp.status_code == 400 and "unique" in resp.text):
        raise AssertionError(resp.text)

def insert_reminder(reminder: Reminder):
    resp = requests.post(
        f"{DIRECTUS_PUBLIC_URL}/items/reminderbot_reminder",
        headers=request_headers,
        json=reminder.model_dump()
    )
    if resp.status_code != 200 and not (resp.status_code == 400 and "unique" in resp.text):
        raise AssertionError(resp.text)

def convert_time(chat_settings: ChatSettings, reminder: Reminder):
    hour, minute = [int(t) for t in reminder.time.split(":")]
    _time = pytz.utc.localize(
        datetime.now()).replace(hour=hour,
                                minute=minute).astimezone(
                                    pytz.timezone(chat_settings.timezone)).strftime("%H:%M")
    return _time


def convert_frequency(reminder: Reminder) -> str:
    if REMINDER_ONCE in reminder.frequency:
        time_str = f"{reminder.frequency.split()[1]}"
        return f'Once-{pytz.utc.localize(datetime.strptime(time_str, "%Y-%m-%d")).strftime("%Y/%m/%d")}'
    elif REMINDER_WEEKLY in reminder.frequency:
        day_of_week = int(reminder.frequency.split('-')[1])
        return f'Weekly-{day_of_week%7}'
    elif REMINDER_YEARLY in reminder.frequency:
        time_str = '-'.join(reminder.frequency.split('-')[1:])
        return f'Yearly-{pytz.utc.localize(datetime.strptime(time_str, "%Y-%m-%d")).strftime("%Y/%m/%d")}'
    
    return reminder.frequency

def calculate_next_trigger_time(chat_settings: ChatSettings, reminder: Reminder) -> str:
    hour, minute = [int(t) for t in reminder.time.split(":")]
    if REMINDER_ONCE in reminder.frequency:
        time_str = f"{'-'.join(reminder.frequency.split('-')[1:])}-{hour}-{minute}"
        run_date = datetime.strptime(time_str, "%Y/%m/%d-%H-%M").astimezone(pytz.timezone(chat_settings.timezone))
        run_date = run_date.astimezone(pytz.utc)
        return run_date.strftime("%Y-%m-%dT%H:%M:00")
    elif REMINDER_DAILY in reminder.frequency:
        # extract hour and minute
        run_date = pytz.timezone(chat_settings.timezone).localize(
            datetime.combine(
                pytz.timezone(chat_settings.timezone).localize(datetime.today()), 
                time(hour, minute)
            )
        )
        run_date = run_date.astimezone(pytz.utc)
        if datetime.now().astimezone(pytz.utc) > run_date:
            run_date = run_date + timedelta(days=1)
        return run_date.strftime("%Y-%m-%dT%H:%M:00")
    elif REMINDER_WEEKLY in reminder.frequency:
        day_of_week = int(reminder.frequency.split('-')[1])
        run_date = pytz.timezone(chat_settings.timezone).localize(
            datetime.combine(
                pytz.timezone(chat_settings.timezone).localize(datetime.today()), 
                time(hour, minute)
            )
        ).replace(day=20)  # middle of the month so that the next calculation won't end with negative day
        run_date = run_date.replace(day=run_date.day - (run_date.isoweekday() % 7 - day_of_week))
        assert run_date.isoweekday() % 7 == day_of_week
        run_date = run_date.astimezone(pytz.utc)
        if datetime.now().astimezone(pytz.utc) > run_date:
            run_date = run_date + timedelta(days=7)
        return run_date.strftime("%Y-%m-%dT%H:%M:00")
    elif REMINDER_MONTHLY in reminder.frequency:
        day = int(reminder.frequency.split('-')[1])
        run_date = pytz.timezone(chat_settings.timezone).localize(
            datetime.combine(
                pytz.timezone(chat_settings.timezone).localize(datetime.today()), 
                time(hour, minute)
            )
        ).replace(day=day)
        run_date = run_date.astimezone(pytz.utc)
        if datetime.now().astimezone(pytz.utc) > run_date:
            run_date = run_date + relativedelta(months=1)
        return run_date.strftime("%Y-%m-%dT%H:%M:00")
    elif REMINDER_YEARLY in reminder.frequency:
        time_str = f"{'-'.join(reminder.frequency.split('-')[1:])}-{hour}-{minute}"
        run_date = datetime.strptime(time_str, "%Y/%m/%d-%H-%M").astimezone(pytz.timezone(chat_settings.timezone))
        run_date = run_date.astimezone(pytz.utc)   
        while datetime.now().astimezone(pytz.utc) > run_date:
            run_date = run_date + relativedelta(years=1)
        return run_date.strftime("%Y-%m-%dT%H:%M:00")

    raise AssertionError("INVALID FREQUENCY")

def main():
    with open(Path(__file__).parent / 'chat_collection.json', "r") as f:
        chat_collection = json.load(f)
    for row in tqdm(chat_collection):
        chat_settings = ChatSettings(
            chat_id=row["chat_id"],
            timezone=row["timezone"],
            updating=False,
        )
        insert_directus_chat_settings(chat_settings)
        for _reminder in row["reminders"]:
            reminder = Reminder(
                chat_id=chat_settings.chat_id,
                from_user_id=_reminder["user_id"],
                file_id=_reminder["file_id"] if "file_id" in _reminder is not None else "",
                reminder_text=_reminder["reminder_text"],
                frequency=_reminder["frequency"],
                time=_reminder["time"],
                in_construction=False,
                next_trigger_time="",
            )
            reminder.time = convert_time(chat_settings, reminder)
            reminder.frequency = convert_frequency(reminder)
            reminder.next_trigger_time = calculate_next_trigger_time(chat_settings, reminder)
            insert_reminder(reminder)


if __name__ == '__main__': 
    main()