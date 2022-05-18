#!/usr/bin/env python3

import os
import json
import concurrent.futures
import requests
import datetime
import random
import threading
import uuid
import time
import string

thread_local = threading.local()
endpoint = os.environ["HOTDOGGIES_ENDPOINT"]
token = os.environ["HOTDOGGIES_TOKEN"]
source = "python-loader"

def load(event):
    while True:
        print(f"Generating event: {event}")
        data = None
        if event == "es.hotdoggi.events.pickup_arrived":
            data = {
                "trip_id": str(uuid.uuid4()),
                "doggy_id": f"dog-{''.join(random.choices(string.ascii_uppercase + string.digits, k=14))}",
                "driver": {
                    "id": f"driver-{''.join(random.choices(string.ascii_uppercase + string.digits, k=5))}"
                },
                "location": {
                    "longitude": random.uniform(44.75000, 44.76999),
                    "latitude": random.uniform(9.12000, 9.13999)
                },
                "status": {
                    "signature_timestamp": str(datetime.datetime.utcnow()),
                    "reference": f"ref-{''.join(random.choices(string.ascii_lowercase + string.digits, k=10))}"
                },
                "timestamp": str(datetime.datetime.utcnow())
            }
        elif event == "es.hotdoggi.events.doggy_arrived":
            data = {
                "trip_id": str(uuid.uuid4()),
                "doggy_id": f"dog-{''.join(random.choices(string.ascii_uppercase + string.digits, k=14))}",
                "driver": {
                    "id": f"driver-{''.join(random.choices(string.ascii_uppercase + string.digits, k=5))}"
                },
                "location": {
                    "longitude": random.uniform(44.75000, 44.76999),
                    "latitude": random.uniform(9.12000, 9.13999)
                },
                "note": "He seems quite happy!",
                "device_timestamp": str(datetime.datetime.utcnow()),
                "timestamp": str(datetime.datetime.utcnow())
            }
        else:
            print("unknown event type!")
            break

        headers = {
            "Authorization": f"Bearer {token}"
        }

        r = requests.post(f"{endpoint}/{event}/{source}", data=json.dumps(data), headers=headers)
        if r.status_code != 201:
            print("error publishing event")

        time.sleep(0.2)

def main():
    events = [
        "es.hotdoggi.events.pickup_arrived",
	    "es.hotdoggi.events.doggy_arrived"
    ]

    with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
        executor.map(load, events)

if __name__ == "__main__":
    main()