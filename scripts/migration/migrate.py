import json
import requests
import os
from datetime import datetime
import pytz
from pydantic import BaseModel
from pathlib import Path
from tqdm import tqdm

DIRECTUS_URL = os.environ.get("DIRECTUS_URL", "http://localhost:8055")
REMINDER_ONCE = 'Once'
REMINDER_DAILY = 'Daily'
REMINDER_WEEKLY = 'Weekly'
REMINDER_MONTHLY = 'Monthly'
REMINDER_YEARLY = 'Yearly'

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
        f"{DIRECTUS_URL}/items/chat_settings",
        json=chat_settings.model_dump()
    )
    if resp.status_code != 200 and not (resp.status_code == 400 and "unique" in resp.text):
        raise AssertionError(resp.text)

def insert_reminder(reminder: Reminder):
    resp = requests.post(
        f"{DIRECTUS_URL}/items/reminder",
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


def convert_frequency(reminder: Reminder):
    if REMINDER_ONCE in reminder.frequency:
        time_str = f"{reminder.frequency.split()[1]}"
        return f'Once-{pytz.utc.localize(datetime.strptime(time_str, "%Y-%m-%d")).strftime("%Y/%m/%d")}'
    elif REMINDER_YEARLY in reminder.frequency:
        time_str = '-'.join(reminder.frequency.split('-')[1:])
        return f'Yearly-{pytz.utc.localize(datetime.strptime(time_str, "%Y-%m-%d")).strftime("%Y/%m/%d")}'
    
    return reminder.frequency

# TODO
# const DIRECTUS_DATETIME_FORMAT = "2006-01-02T15:04:05"
def calculate_next_trigger_time(reminder: Reminder):
    hour, minute = [int(t) for t in reminder.time.split(":")]
    if REMINDER_ONCE in reminder.frequency:
        time_str = f"{reminder.frequency.split()[1]}-{hour}-{minute}"
        run_date = pytz.utc.localize(datetime.strptime(time_str, "%Y-%m-%d-%H-%M")).strftime("%Y-%m-%d")
        # scheduler.add_job(reminder_trigger,
        #                   'date',
        #                   run_date=run_date,
        #                   args=[chat_id, reminder_id],
        #                   id=job_id)
    elif REMINDER_DAILY in reminder.frequency:
        # extract hour and minute
        run_date = datetime.combine(datetime.today(),
                                    time(hour, minute)).replace(day=10)
        run_date = pytz.utc.localize(run_date)
        # scheduler.add_job(reminder_trigger,
        #                   'cron',
        #                   day="*",
        #                   hour=run_date.hour,
        #                   minute=run_date.minute,
        #                   args=[chat_id, reminder_id],
        #                   id=job_id)
    elif REMINDER_WEEKLY in reminder.frequency:
        day = int(reminder.frequency.split('-')[1]) - 1
        hour, minute = [
            int(t)
            for t in convert_time_str(f"{hour}:{minute}", timezone).split(":")
        ]
        run_date = datetime.combine(datetime.today(
        ), time(hour, minute)).replace(
            day=20
        )  # middle of the month so that the next calculation won't end with negative day
        run_date = run_date.replace(day=run_date.day -
                                    (run_date.weekday() - day))
        run_date = pytz.timezone(timezone).localize(run_date).astimezone(
            pytz.utc)
        # scheduler.add_job(
        #     reminder_trigger,
        #     'cron',
        #     week="*",
        #     day_of_week=run_date.weekday(),  # day of week goes from 0-6
        #     hour=run_date.hour,
        #     minute=run_date.minute,
        #     args=[chat_id, reminder_id],
        #     id=job_id)
    elif REMINDER_MONTHLY in reminder.frequency:
        day = int(reminder.frequency.split('-')[1])
        hour, minute = [
            int(t)
            for t in convert_time_str(f"{hour}:{minute}", timezone).split(":")
        ]
        run_date = datetime.combine(datetime.today(),
                                    time(hour, minute)).replace(day=day)
        run_date = pytz.timezone(timezone).localize(run_date).astimezone(
            pytz.utc)
        # scheduler.add_job(reminder_trigger,
        #                   'cron',
        #                   month="*",
        #                   day=run_date.day,
        #                   hour=run_date.hour,
        #                   minute=run_date.minute,
        #                   args=[chat_id, reminder_id],
        #                   id=job_id)
    elif REMINDER_YEARLY in reminder.frequency:
        _, month, day = [int(num) for num in reminder.frequency.split('-')[1:]]
        hour, minute = [
            int(t)
            for t in convert_time_str(f"{hour}:{minute}", timezone).split(":")
        ]
        run_date = datetime.combine(datetime.today(),
                                    time(hour, minute)).replace(month=month, day=day)
        run_date = pytz.timezone(timezone).localize(run_date).astimezone(
            pytz.utc)
        # scheduler.add_job(reminder_trigger,
        #                   'cron',
        #                   year="*",
        #                   month=run_date.month,
        #                   day=run_date.day,
        #                   hour=run_date.hour,
        #                   minute=run_date.minute,
        #                   args=[chat_id, reminder_id],
        #                   id=job_id)


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
            insert_reminder(reminder)


if __name__ == '__main__': 
    main()