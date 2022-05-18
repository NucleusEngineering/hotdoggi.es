import firebase from "firebase/app";
import "firebase/auth";

const firebaseConfig = {
  apiKey: "AIzaSyCvf5TQolHXBDGF_28tNEgZATn0LvHi6bQ",
  authDomain: "hwsh-api.firebaseapp.com",
  projectId: "hwsh-api",
  storageBucket: "hwsh-api.appspot.com",
  messagingSenderId: "546978254761",
  appId: "1:546978254761:web:37eda206bbe04ad2d77eb8"
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
