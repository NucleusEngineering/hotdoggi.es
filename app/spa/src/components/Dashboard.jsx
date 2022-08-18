//  Copyright 2022 Google

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at

//      http://www.apache.org/licenses/LICENSE-2.0

//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

import React, { useContext } from "react";
import { UserContext } from "../providers/UserProvider";
import { auth } from "../firebase";
import DogPark from "./DogPark/DogPark";

const Dashboard = () => {
  const user = useContext(UserContext);
  const { photoURL, displayName, email } = user;

  console.log(user);

  return (
    <div>
      <div className="flex flex-row items-center justify-end">
        <div
          style={{
            background: `url(${
              photoURL ||
              "https://res.cloudinary.com/dqcsk8rsc/image/upload/v1577268053/avatar-1-bitmoji_upgwhc.png"
            })  no-repeat center center`,
            backgroundSize: "cover",
            height: "50px",
            width: "50px",
          }}
          className="border border-blue-300 rounded-full"
        ></div>
        <div className="flex-col px-4">
          <h2 className="text-2xl font-semibold">{displayName}</h2>
          <h3 className="italic">{email}</h3>
        </div>
        <button
          className="w-24 py-3 bg-red-600 hover:bg-red-700 text-white rounded-full"
          onClick={() => {
            auth.signOut();
          }}
        >
          Sign out
        </button>
      </div>
      <DogPark />
    </div>
  );
};

export default Dashboard;
