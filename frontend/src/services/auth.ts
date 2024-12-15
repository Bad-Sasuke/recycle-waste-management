import { getAuth,  GoogleAuthProvider, signInWithPopup, type UserCredential } from "firebase/auth";
import config from '../config';
import { initializeApp } from "firebase/app";
import { firebaseConfig } from './firebase';

const app = initializeApp(firebaseConfig);
const auth = getAuth(app);
const provider = new GoogleAuthProvider();
const authLoginGithub = async () => {
    try {
        const url = new URL('https://github.com/login/oauth/authorize');
        url.searchParams.append('client_id', config.githubClientId);
        url.searchParams.append('scope', 'user:email');
        window.location.href = url.toString();
    } catch (error) {
        console.log(error);

    }
}

const authLoginGoogle = async () => {
  try {
      // เริ่มกระบวนการ Sign In
      const result : UserCredential = await signInWithPopup(auth, provider);
      const user = result.user;
      const userID = user.uid;
      console.log(userID);
      // send api google/callback?uid=${userID}
      const url = new URL(config.webAPI + '/api/auth/google/callback' + '?uid=' + userID);
      await fetch(url);

  } catch (error) {
      console.log("Error during redirect login:", error);
  }
};

export { authLoginGithub,authLoginGoogle }
