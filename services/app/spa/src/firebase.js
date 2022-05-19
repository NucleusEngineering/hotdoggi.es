import firebase from "firebase/app";
import "firebase/auth";


const firebaseConfig = {
  apiKey: "", // TODO replace
  authDomain: "hotdoggi-es.firebaseapp.com",
  projectId: "hotdoggi-es",
  storageBucket: "hotdoggi-es.appspot.com",
  messagingSenderId: "", // TODO replace
  appId: "" // TODO replace
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
