# Copyright 2022 Google

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/usr/bin/env python3

# This loader script generates random events to be ingested into the API.
#
# 1.) Create a pack of dogs
# 2.) Wait for dog creation
# 3.) Loop until CTRL-C received and update dogs concurrently
#   a.) Refresh current dog my querying API
#   b.) Alter dog's current location by randomly incrementing or decrementing latitude and longitude
#   c.) Send updated dog location as event/command to API (dog_moved)

import os
import json
import concurrent.futures
import requests
import datetime
import random
import signal
import time
from random import randrange
from datetime import datetime, timedelta

# API Endpoint
endpoint = "https://api.hotdoggies.stamer.demo.altostrat.com"

# JWT access token for API access
token = os.environ["TOKEN"]
headers = {"Authorization": f"Bearer {token}"}

# Static event source name
source = "python-loader"

# Number of dogs to simulate
pack_size = 4
thread_executor = concurrent.futures.ThreadPoolExecutor(max_workers=pack_size)
terminate = False

# Lower and upper bounds for coordinates to be display in frontend dog enclosure grid
coord_lower_bound = 0
coord_upper_bound = 7

# Color set for pretty printing
class colors:
    BLUE = '\033[94m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    WHITE = '\033[0m'


def addRandomDog():
    """ Create a random dog and push them to the API as event (dog_added) """
    data = {
        "dog": {
            "name": randomName(),
            "breed": randomBreed(),
            "birthday": randomBirthday(),
            "color": randomColor(),
            "picture": "https://i1.sndcdn.com/artworks-UWJlJnEMrYzN2Vqx-1ImGMw-t500x500.jpg",
            "location": {
                "longitude": randomCoordinate(),
                "latitude": randomCoordinate()
            }
        }
    }

    event_type="es.hotdoggi.events.dog_added"
    print(f"{colors.BLUE}ASYNC[{event_type}]\t{colors.WHITE} creating {colors.RED}{data['dog']['name']}{colors.WHITE} ({data['dog']['color']} {data['dog']['breed']})")
    
    r = requests.post(f"{endpoint}/v1/events/{event_type}/{source}", data=json.dumps(data), headers=headers)
    if r.status_code != 201:
        print("error publishing event")


def getAllDogs():
    """ Query for all dogs of the current user """
    r = requests.get(f"{endpoint}/v1/dogs/", headers=headers)
    print(f"{colors.GREEN} SYNC[dogs/*]\t\t\t\t{colors.WHITE} listing all dogs")
    if r.status_code > 299:
        print("error getting dogs")
    return json.loads(r.text)


def getDog(dog):
    """ Query for a single specific dog of the current user """
    r = requests.get(f"{endpoint}/v1/dogs/{dog['id']}", headers=headers)
    print(f"{colors.GREEN} SYNC[dogs/{dog['id']}]\t{colors.WHITE} update {colors.RED}{dog['dog']['name']}{colors.WHITE}")
    if r.status_code > 299:
        print("error getting dog")
    return json.loads(r.text)


def simulateDogMovement(dog):
    """ Simulate dog movement by randomly updating coordinates and pushing dog_moved events to the API """
    while not terminate:
        update = getDog(dog)
        new_latitude = randomMovement(update['dog']['location']['latitude'])
        new_longitude = randomMovement(update['dog']['location']['longitude'])
        data = {
            "id": dog['id'],
            "dog": {
                "location": {
                    "latitude": new_latitude,
                    "longitude": new_longitude
                }
            } 
        }
        event_type = "es.hotdoggi.events.dog_moved"
        print(f"{colors.BLUE}ASYNC[{event_type}]\t{colors.WHITE} moving {colors.RED}{dog['dog']['name']}{colors.WHITE} to {colors.YELLOW}({new_latitude},{new_longitude}){colors.WHITE}")

        r = requests.post(f"{endpoint}/v1/events/{event_type}/{source}", data=json.dumps(data), headers=headers)
        if r.status_code != 201:
            print("error publishing event")

        time.sleep(random.uniform(4.0, 8.0))


def removeDog(dog):
    """ Remove a specific dog via the API (dog_removed) """
    data = {
        "id": dog['id']
    }

    event_type = "es.hotdoggi.events.dog_removed"
    print(f"{colors.BLUE}ASYNC[{event_type}]\t{colors.WHITE} removing {colors.RED}{dog['dog']['name']}{colors.WHITE} (id {dog['id']})")

    r = requests.post(f"{endpoint}/v1/events/{event_type}/{source}", data=json.dumps(data), headers=headers)
    if r.status_code != 201:
        print("error publishing event")


def randomMovement(coord):
    """ Generate dog movement by randomly incrementing or decrementing latitude and longitude """
    coord += random.choice((-1, 1))

    if coord > coord_upper_bound:
        coord = coord_upper_bound

    if coord < coord_lower_bound:
        coord = coord_lower_bound
    
    return coord


def randomName():
    """ Get a random dog name """
    names = ["Max","Kobe","Oscar","Cooper","Oakley","Mac","Charlie","Rex","Rudy","Teddy","Ailey","Chip","Bear","Cash","Walter","Milo","Jasper","Blaze","Bentley","Bo","Ozzy","Ollie","Boomer","Odin","Buddy","Lucky","Axel","Rocky","Ruger","Bruce","Leo","Beau","Odie","Zeus","Baxter","Arlo","Duke","Oreo","Echo","Finn","Gunner","Tank","Apollo","Henry","Romeo","Murphy","Simba","Porter","Diesel","George","Harley","Toby","Coco","Otis","Louie","Rocket","Rocco","Tucker","Ziggy","Remi","Jax","Prince","Whiskey","Ace","Shadow","Sam","Jack","Riley","Buster","Koda","Copper","Bubba","Winston","Luke","Jake","Oliver","Marley","Benny","Gus","Zeke","Bowie","Loki","Levi","Dozer","Moose","Benji","Rusty","Archie","Ranger","Joey","Bandit","Remy","Kylo","Scout","Dexter","Ryder","Thor","Gizmo","Tyson","Bruno","Chase","Samson","King","Cody","Rambo","Blue","Sarge","Harry","Atlas","Chester","Gucci","Theo","Maverick","Miles","Jackson","Lincoln","Watson","Hank","Wally","Peanut","Titan"]
    return random.choice(names)


def randomBreed():
    """ Get a random dog breed """
    breeds = ['Boxer', 'Bulldog', 'Chiuahua', 'Golden', 'Husky', 'Pincer', 'Pomeranian', 'Rottweiler']
    return random.choice(breeds)


def randomCoordinate():
    """ Get a single random coordinate """

    return random.choice(range(coord_lower_bound, coord_upper_bound+1))


def randomBirthday():
    """ Get a random dog birthday """
    latest = datetime.today() - timedelta(days=60) # 60 days ago
    oldest = datetime.today() - timedelta(days=(18*365)) # 18 years ago
    delta = latest - oldest
    random_day = randrange(delta.days)
    birthday = oldest + timedelta(days=random_day)
    return birthday.strftime("%Y-%m-%d")


def randomColor():
    """ Get a random dog fur color """
    colors = ["Brown","Dark Chocolate","Red","Black","White","Gold","Yellow","Cream","Blue","Grey"]
    return random.choice(colors)


def abortHandler(_, __):
    """ CTRL-C listener for graceful exits """
    print("\nCaught exit... Suspending movement simulation")
    global terminate
    terminate = True
    thread_executor.shutdown
    time.sleep(12)
    print("Removing dogs from the pack...")
    dogs = getAllDogs()
    for dog in dogs:
        removeDog(dog)
    
    print("Clean exit.")

signal.signal(signal.SIGINT, abortHandler)


def main():
    random.seed()
    print("Adding some dogs to the pack...")
    for _ in range(pack_size):
        addRandomDog()
    
    print("Waiting 10 seconds for dog registration...")
    time.sleep(10)

    dogs = getAllDogs()
    print(f"Found {len(dogs)} dogs in the pack.")

    print("Simulating movement ...")
    thread_executor.map(simulateDogMovement, dogs)


if __name__ == "__main__":
    main()
