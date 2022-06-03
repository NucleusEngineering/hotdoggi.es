#!/usr/bin/env python3

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

endpoint = "https://api.hotdoggies.stamer.demo.altostrat.com"
token = os.environ["TOKEN"]
headers = {"Authorization": f"Bearer {token}"}
source = "python-loader"
pack_size = 8
thread_executor = concurrent.futures.ThreadPoolExecutor(max_workers=pack_size)
terminate = False

class colors:
    BLUE = '\033[94m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    WHITE = '\033[0m'

def addRandomDog():
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
    
    r = requests.post(f"{endpoint}/events/{event_type}/{source}", data=json.dumps(data), headers=headers)
    if r.status_code != 201:
        print("error publishing event")

def getAllDogs():
    r = requests.get(f"{endpoint}/dogs/", headers=headers)
    print(f"{colors.GREEN} SYNC[dogs/*]\t\t\t\t{colors.WHITE} listing all dogs")
    if r.status_code > 299:
        print("error getting dogs")
    return json.loads(r.text)

def getDog(dog):
    r = requests.get(f"{endpoint}/dogs/{dog['id']}", headers=headers)
    print(f"{colors.GREEN} SYNC[dogs/{dog['id']}]\t{colors.WHITE} update {colors.RED}{dog['dog']['name']}{colors.WHITE}")
    if r.status_code > 299:
        print("error getting dog")
    return json.loads(r.text)

def simulateDogMovement(dog):
    while not terminate:
        update = getDog(dog)
        data = {
            "id": dog['id'],
            "dog": {
                "location": {
                    "longitude": update['dog']['location']['longitude'] + (0.001 * random.choice((-1, 1))),
                    "latitude": update['dog']['location']['latitude'] + (0.001 * random.choice((-1, 1)))
                }
            } 
        }
        event_type = "es.hotdoggi.events.dog_moved"
        print(f"{colors.BLUE}ASYNC[{event_type}]\t{colors.WHITE} moving {colors.RED}{dog['dog']['name']}{colors.WHITE} to {colors.YELLOW}({data['dog']['location']['latitude']},{data['dog']['location']['latitude']}){colors.WHITE}")

        r = requests.post(f"{endpoint}/events/{event_type}/{source}", data=json.dumps(data), headers=headers)
        if r.status_code != 201:
            print("error publishing event")

        time.sleep(random.uniform(8.0, 12.0))

def removeDog(dog):
    data = {
        "id": dog['id']
    }

    event_type = "es.hotdoggi.events.dog_removed"
    print(f"{colors.BLUE}ASYNC[{event_type}]\t{colors.WHITE} removing {colors.RED}{dog['dog']['name']}{colors.WHITE} (id {dog['id']})")

    r = requests.post(f"{endpoint}/events/{event_type}/{source}", data=json.dumps(data), headers=headers)
    if r.status_code != 201:
        print("error publishing event")

def randomName():
    names = ["Max","Kobe","Oscar","Cooper","Oakley","Mac","Charlie","Rex","Rudy","Teddy","Ailey","Chip","Bear","Cash","Walter","Milo","Jasper","Blaze","Bentley","Bo","Ozzy","Ollie","Boomer","Odin","Buddy","Lucky","Axel","Rocky","Ruger","Bruce","Leo","Beau","Odie","Zeus","Baxter","Arlo","Duke","Oreo","Echo","Finn","Gunner","Tank","Apollo","Henry","Romeo","Murphy","Simba","Porter","Diesel","George","Harley","Toby","Coco","Otis","Louie","Rocket","Rocco","Tucker","Ziggy","Remi","Jax","Prince","Whiskey","Ace","Shadow","Sam","Jack","Riley","Buster","Koda","Copper","Bubba","Winston","Luke","Jake","Oliver","Marley","Benny","Gus","Zeke","Bowie","Loki","Levi","Dozer","Moose","Benji","Rusty","Archie","Ranger","Joey","Bandit","Remy","Kylo","Scout","Dexter","Ryder","Thor","Gizmo","Tyson","Bruno","Chase","Samson","King","Cody","Rambo","Blue","Sarge","Harry","Atlas","Chester","Gucci","Theo","Maverick","Miles","Jackson","Lincoln","Watson","Hank","Wally","Peanut","Titan"]
    return random.choice(names)

def randomBreed():
    breeds = ["Affenpinscher","Afghan Hound","Afghan Shepherd","Aidi","Airedale Terrier","Akhttps://console.developers.google.com/apis/api/firestorekeyvisualizer.googleapis.com/overview?project=640843850686bash","Akita","Alano Español","Alaskan husky","Alaskan Klee Kai","Alaskan Malamute","Alaunt","Alopekis","Alpine Dachsbracke","Alpine Mastiff","Alpine Spaniel","American Akita","American Bulldog","American Cocker Spaniel","American English Coonhound","American Eskimo Dog","American Foxhound","American Hairless Terrier","American Pit Bull Terrier","American Staffordshire Terrier","American Water Spaniel","Anatolian Shepherd Dog","Andalusian Hound","Anglo-Français de Petite Vénerie","Appenzeller Sennenhund","Braque de l'Ariege","Ariegeois","Armant","Armenian Gampr dog","Artois Hound","Australian Cattle Dog","Australian Kelpie","Australian Shepherd","Australian Silky Terrier","Australian Stumpy Tail Cattle Dog[10]","Australian Terrier","Austrian Black and Tan Hound","Austrian Pinscher","Azawakh","Bakharwal Dog","Barbet","Basenji","Basque Ratter","Basque Shepherd Dog","Basset Artésien Normand","Basset Bleu de Gascogne","Basset Fauve de Bretagne","Basset Griffon Vendéen, Grand","Basset Griffon Vendéen, Petit","Basset Hound","Bavarian Mountain Hound","Beagle","Beagle-Harrier","Bearded Collie","Beauceron","Bedlington Terrier","Belgian Shepherd Dog (Groenendael)","Belgian Shepherd Dog (Laekenois)","Belgian Shepherd Dog (Malinois)","Belgian Shepherd Dog (Tervuren)","Bergamasco Shepherd","Berger Blanc Suisse","Berger Picard","Bernese Mountain Dog","Bichon Frisé","Billy","Black and Tan Coonhound","Black and Tan Virginia Foxhound","Black Norwegian Elkhound","Black Russian Terrier","Black Mouth Cur","Bleu de Gascogne, Grand","Bleu de Gascogne, Petit","Bloodhound","Blue Heeler","Blue Lacy","Blue Paul Terrier","Blue Picardy Spaniel","Bluetick Coonhound","Boerboel","Bohemian Shepherd","Bolognese","Border Collie","Border Terrier","Borzoi","Bosnian Coarse-haired Hound","Boston Terrier","Bouvier des Ardennes","Bouvier des Flandres","Boxer","Boykin Spaniel","Bracco Italiano","Braque d'Auvergne","Braque du Bourbonnais","Braque du Puy","Braque Francais","Braque Saint-Germain","Brazilian Dogo","Brazilian Terrier","Briard","Briquet Griffon Vendéen","Brittany","Broholmer","Bruno Jura Hound","Bucovina Shepherd Dog","Bull and Terrier","Bull Terrier","Bull Terrier (Miniature)","Bulldog","Bullenbeisser","Bullmastiff","Bully Kutta","Burgos Pointer","Cairn Terrier","Canaan Dog","Canadian Eskimo Dog","Cane Corso","Cantabrian Water Dog","Cão da Serra de Aires","Cão de Castro Laboreiro","Cão de Gado Transmontano","Cão Fila de São Miguel","Carolina Dog","Carpathian Shepherd Dog","Catahoula Leopard Dog","Catalan Sheepdog","Caucasian Shepherd Dog","Cavalier King Charles Spaniel","Central Asian Shepherd Dog","Cesky Fousek","Cesky Terrier","Chesapeake Bay Retriever","Chien Français Blanc et Noir","Chien Français Blanc et Orange","Chien Français Tricolore","Chien-gris","Chihuahua","Chilean Fox Terrier","Chinese Chongqing Dog","Chinese Crested Dog","Chinese Imperial Dog","Chinook","Chippiparai","Chow Chow","Cierny Sery","Cirneco dell'Etna","Clumber Spaniel","Collie, Rough","Collie, Smooth","Combai","Cordoba Fighting Dog","Coton de Tulear","Cretan Hound","Croatian Sheepdog","Cumberland Sheepdog","Curly-Coated Retriever","Cursinu","Czechoslovakian Wolfdog","Dachshund","Dalmatian","Dandie Dinmont Terrier","Danish-Swedish Farmdog","Deutsche Bracke","Doberman Pinscher","Dogo Argentino","Dogo Cubano","Dogue de Bordeaux","Drentse Patrijshond","Drever","Dunker","Dutch Shepherd","Dutch Smoushond","East Siberian Laika","East European Shepherd","Elo","English Cocker Spaniel","English Foxhound","Mastiff","English Setter","English Shepherd","English Springer Spaniel","English Toy Terrier (Black & Tan)","English Water Spaniel","English White Terrier","Entlebucher Mountain Dog","Estonian Hound","Estrela Mountain Dog","Eurasier","Eurohound","Field Spaniel","Fila Brasileiro","Finnish Hound","Finnish Lapphund","Finnish Spitz","Flat-Coated Retriever","Fox Terrier, Smooth","Fox Terrier, Wire","French Brittany","French Bulldog","French Spaniel","Gaddi Dog","Galgo Español","Galician Cattle Dog","Garafian Shepherd","Gascon Saintongeois","Georgian Shepherd Dog","German Longhaired Pointer","German Pinscher","German Roughhaired Pointer","German Shepherd Dog","German Shorthaired Pointer","German Spaniel","German Spitz","German Wirehaired Pointer","Giant Schnauzer","Glen of Imaal Terrier","Golden Retriever","Gordon Setter","Gran Mastín de Borínquen","Grand Anglo-Français Blanc et Noir","Grand Anglo-Français Blanc et Orange","Grand Anglo-Français Tricolore","Grand Griffon Vendéen","Great Dane","Great Pyrenees","Greater Swiss Mountain Dog","Hellenic Hound","Greenland Dog","Greyhound","Griffon Bleu de Gascogne","Griffon Bruxellois","Griffon Fauve de Bretagne","Griffon Nivernais","Guatemalan Dogo","Hamiltonstövare","Hanover Hound","Hare Indian Dog","Harrier","Havanese","Hawaiian Poi Dog","Himalayan Sheepdog","Hokkaido","Hortaya Borzaya","Hovawart","Huntaway","Hygenhund","Ibizan Hound","Icelandic Sheepdog","Indian pariah dog","Indian Spitz","Irish Red and White Setter","Irish Setter","Irish Terrier","Irish Water Spaniel","Irish Wolfhound","Istrian Coarse-haired Hound","Istrian Shorthaired Hound","Italian Greyhound","Jack Russell Terrier","Jagdterrier","Jämthund","Japanese Chin","Japanese Spitz","Japanese Terrier","Kaikadi","Kai Ken","Kangal Dog","Kanni","Karakachan Dog","Karelian Bear Dog","Karst Shepherd","Keeshond","Kerry Beagle","Kerry Blue Terrier","King Charles Spaniel","King Shepherd","Kintamani","Kishu Ken","Komondor","Kooikerhondje","Koolie","Korean Jindo","Kromfohrländer","Kumaon Mastiff","Kunming Wolfdog","Kurī","Kuvasz","Kyi-Leo","Labrador Husky","Labrador Retriever","Lagotto Romagnolo","Lakeland Terrier","Lancashire Heeler","Landseer","Lapponian Herder","Leonberger","Lhasa Apso","Lithuanian Hound","Löwchen","Mackenzie River Husky","Magyar agár","Mahratta Greyhound","Majorca Ratter","Majorca Shepherd Dog","Maltese","Manchester Terrier","Maremma Sheepdog","McNab","Mexican Hairless Dog","Miniature Australian Shepherd","Miniature American Shepherd","Miniature Fox Terrier","Miniature Pinscher","Miniature Schnauzer","Miniature Shar Pei","Romanian Mioritic Shepherd Dog","Molossus","Molossus of Epirus","Montenegrin Mountain Hound","Moscow Watchdog","Moscow Water Dog","Mountain Cur","Mucuchies","Mudhol Hound","Mudi","Münsterländer, Large","Münsterländer, Small","Neapolitan Mastiff","Newfoundland","New Zealand Heading Dog","Norfolk Spaniel","Norfolk Terrier","Norrbottenspets","North Country Beagle","Northern Inuit Dog","Norwegian Buhund","Norwegian Elkhound","Norwegian Lundehund","Norwich Terrier","Nova Scotia Duck Tolling Retriever","Old Croatian Sighthound","Old Danish Pointer","Old English Sheepdog","Old English Terrier","Old German Shepherd Dog","Old Time Farm Shepherd","Olde English Bulldogge","Otterhound","Pachon Navarro","Pandikona Hunting Dog","Paisley Terrier","Papillon","Parson Russell Terrier","Patterdale Terrier","Pekingese","Perro de Presa Canario","Perro de Presa Mallorquin","Peruvian Hairless Dog","Phalène","Pharaoh Hound","Phu Quoc Ridgeback","Picardy Spaniel","Plummer Terrier","Plott Hound","Podenco Canario","Pointer","Poitevin","Polish Greyhound","Polish Hound","Polish Hunting Dog","Polish Lowland Sheepdog","Polish Tatra Sheepdog","Pomeranian","Pont-Audemer Spaniel","Poodle","Porcelaine","Portuguese Podengo","Portuguese Pointer","Portuguese Water Dog","Posavac Hound","Pražský Krysařík","Pudelpointer","Pug","Puli","Pumi","Pungsan Dog","Pyrenean Mastiff","Pyrenean Shepherd","Rafeiro do Alentejo","Rajapalayam","Rampur Greyhound","Rastreador Brasileiro","Ratonero Bodeguero Andaluz","Ratonero Murciano de Huerta","Ratonero Valenciano","Rat Terrier","Redbone Coonhound","Rhodesian Ridgeback","Rottweiler","Russian Spaniel","Russian Toy","Russian Tracker","Russo-European Laika","Russell Terrier","Saarloos Wolfdog","Sabueso Español","Sabueso fino Colombiano","Saint-Usuge Spaniel","Sakhalin Husky","Saluki","Samoyed","Sapsali","Šarplaninac","Schapendoes","Schillerstövare","Schipperke","Standard Schnauzer","Schweizer Laufhund","Schweizerischer Niederlaufhund","Scotch Collie","Scottish Deerhound","Scottish Terrier","Sealyham Terrier","Segugio Italiano","Seppala Siberian Sleddog","Serbian Hound","Serbian Tricolour Hound","Seskar Seal Dog","Shar Pei","Shetland Sheepdog","Shiba Inu","Shih Tzu","Shikoku","Shiloh Shepherd","Siberian Husky","Silken Windhound","Sinhala Hound","Skye Terrier","Sloughi","Slovak Cuvac","Slovakian Rough-haired Pointer","Slovenský Kopov","Smålandsstövare","Small Greek Domestic Dog","Soft-Coated Wheaten Terrier","South Russian Ovcharka","Southern Hound","Spanish Mastiff","Spanish Water Dog","Spinone Italiano","Sporting Lucas Terrier","St. Bernard","St. John's water dog","Stabyhoun","Staffordshire Bull Terrier","Stephens Cur","Styrian Coarse-haired Hound","Sussex Spaniel","Swedish Lapphund","Swedish Vallhund","Tahltan Bear Dog","Taigan","Taiwan Dog","Talbot","Tamaskan Dog","Teddy Roosevelt Terrier","Telomian","Tenterfield Terrier","Terceira Mastiff","Thai Bangkaew Dog","Thai Ridgeback","Tibetan Mastiff","Tibetan Spaniel","Tibetan Terrier","Tornjak","Tosa","Toy Bulldog","Toy Fox Terrier","Toy Manchester Terrier","Toy Trawler Spaniel","Transylvanian Hound","Treeing Cur","Treeing Tennessee Brindle","Treeing Walker Coonhound","Trigg Hound","Tweed Water Spaniel","Tyrolean Hound","Cimarrón Uruguayo","Valencian Ratter","Vanjari Hound","Villano de Las Encartaciones","Vizsla","Volpino Italiano","Weimaraner","Welsh Corgi, Cardigan","Welsh Corgi, Pembroke","Welsh Sheepdog","Welsh Springer Spaniel","Welsh Terrier","West Highland White Terrier","West Siberian Laika","Westphalian Dachsbracke","Wetterhoun","Whippet","White Shepherd","Wirehaired Pointing Griffon","Wirehaired Vizsla","Xiasi Dog","Yorkshire Terrier"]
    return random.choice(breeds)

def randomCoordinate():
    return random.uniform(8.0, 52.0)

def randomBirthday():
    latest = datetime.today() - timedelta(days=60) # 60 days ago
    oldest = datetime.today() - timedelta(days=(18*365)) # 18 years ago
    delta = latest - oldest
    random_day = randrange(delta.days)
    birthday = oldest + timedelta(days=random_day)
    return birthday.strftime("%Y-%m-%d")

def randomColor():
    colors = ["Brown","Dark Chocolate","Red","Black","White","Gold","Yellow","Cream","Blue","Grey"]
    return random.choice(colors)

def abortHandler(signum, frame):
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