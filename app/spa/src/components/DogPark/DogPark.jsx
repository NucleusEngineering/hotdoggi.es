import React, { useState, useContext, useEffect } from "react";

import "./DogPark.scss";
import Dog from "./Dog";
import { UserContext } from "../../providers/UserProvider";

export default function DogPark() {
  const idToken = useContext(UserContext).idToken;

  const [dogs, setDogs] = useState([]);

  useEffect(() => {
    async function fetchData() {
      if (idToken) {
        const response = await fetch(
          "https://api.hotdoggies.stamer.demo.altostrat.com/v1/dogs/",
          {
            headers: {
              Authorization: `Bearer ${idToken}`,
            },
          }
        );
        const dogsJson = await response.json();
        setDogs(dogsJson);
      }
    }

    function startWebsocket() {
      const dogSocket = new WebSocket(
        `wss://api.hotdoggies.stamer.demo.altostrat.com/v1/dogs/?access_token=${idToken}`
      );
      dogSocket.onmessage = (event) => {
        let dogsState = [...dogs];
        const parsedEventPayload = JSON.parse(event.data);
        console.log(parsedEventPayload)
        let dogToBeUpdated = dogsState.find(
          (doggo) => doggo.id === parsedEventPayload.id
        );
        if (dogToBeUpdated) {
          console.log("dog", dogToBeUpdated.dog);
          dogToBeUpdated.dog.location.latitude =
            parsedEventPayload.dog.location.latitude;
          dogToBeUpdated.dog.location.longitude =
            parsedEventPayload.dog.location.longitude;
          setDogs(dogsState);
        } else {
          setDogs([...dogs, parsedEventPayload  ])
        }
      };
    }

    if (idToken) {
      fetchData();
      startWebsocket()
    }
  }, [idToken]);

  return (
    <div>
      <div className="enclosure">
        <div className="tabs">
          <button className="tab active">
            <div>
              <span>Doggies Enclosure</span>
            </div>
          </button>
        </div>
        <div className="playground">
          <div className="playground__inner" id="tab-1">
            <div className="grid">
              {dogs.map((dogObject) => (
                <Dog
                  key={dogObject.dog.name}
                  name={dogObject.dog.name}
                  type={dogObject.dog.breed.toLowerCase()}
                  location={{
                    x: dogObject.dog.location.longitude,
                    y: dogObject.dog.location.latitude,
                  }}
                />
              ))}
            </div>
          </div>
        </div>
      </div>
      <h1 className="mt-16">Hot doggos are ready to roam in the park</h1>
    </div>
  );
}
