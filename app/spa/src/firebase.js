import firebase from "firebase/app";
import "firebase/auth";


const firebaseConfig = {
  apiKey: "AIzaSyAZL8XJG7781vO-IaiN3Ej_F0Y5HUphwyQ",
  authDomain: "hotdoggi-es.firebaseapp.com",
  projectId: "hotdoggi-es",
  storageBucket: "hotdoggi-es.appspot.com",
  messagingSenderId: "640843850686",
  appId: "1:640843850686:web:da0016555b4396b8870fef"
};
firebase.initializeApp(firebaseConfig);

export const auth = firebase.auth();

const googleProvider = new firebase.auth.GoogleAuthProvider();
export const signInWithGoogle = () => {
  auth.signInWithPopup(googleProvider);
};

const githubProvider = new firebase.auth.GithubAuthProvider();
export const signInWithGithub = () => {
  auth.signInWithPopup(githubProvider);
};
