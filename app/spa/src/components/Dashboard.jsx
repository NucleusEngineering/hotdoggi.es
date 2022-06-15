//  Copyright 2022 Daniel Stamer

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
import {auth} from "../firebase";
const Dashboard = () => {
  const user = useContext(UserContext);
  const {photoURL, displayName, email, uid, _lat} = user;
  console.log(user);
  

  return (
    <div className = "mx-auto w-11/12 md:w-2/4 py-8 px-4 md:px-8">
      <div className="flex border flex-col items-center md:flex-row md:items-start border-blue-400 px-3 py-4">
        <div
          style={{
            background: `url(${photoURL || 'https://res.cloudinary.com/dqcsk8rsc/image/upload/v1577268053/avatar-1-bitmoji_upgwhc.png'})  no-repeat center center`,
            backgroundSize: "cover",
            height: "200px",
            width: "200px"
          }}
          className="border border-blue-300"
        ></div>
        <div className = "md:pl-4">
        <h2 className = "text-2xl font-semibold">{displayName}</h2>
        <h3 className = "italic">{email}</h3>
        </div>
      </div>
      <button className = "w-full py-3 bg-red-600 mt-4 text-white" onClick = {() => {auth.signOut()}}>Sign out</button>
      <div className = "md:pl-4">Debug Info</div>
      <div>
        <table style={{
            tableLayout: "fixed"
            }}>
        <tbody style={{
            border: "1px solid black"
            }}>
          <tr style={{
            border: "1px solid black"
            }}>
            <td><code>{uid}</code></td>
          </tr >
          <tr style={{
            border: "1px solid black"
            }}>
            <td style={{
            wordBreak: "break-all"
            }}><code>{_lat}</code></td>
          </tr>
        </tbody>
        </table>
      </div>
    </div>
  ) 
};

export default Dashboard;