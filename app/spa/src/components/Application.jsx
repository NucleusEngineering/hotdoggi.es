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
import { Router } from "@reach/router";
import SignIn from "./SignIn";
import SignUp from "./SignUp";
import Dashboard from "./Dashboard";
import { UserContext } from "../providers/UserProvider";
import PasswordReset from "./PasswordReset";

function Application() {
  const { user, loading } = useContext(UserContext);
  if (loading) {
    return null;
  } else if (user) {
    return <Dashboard />;
  } else {
    return (
      <Router>
        <SignUp path="signUp" />
        <SignIn path="/" />
        <PasswordReset path="passwordReset" />
      </Router>
    );
  }
}

export default Application;
