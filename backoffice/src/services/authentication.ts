import axios from "axios";

type Token = {
    accessToken: string;
    refreshToken: string;
}

type SignInResponse = {
    token: string;
    refresh_token: string;
}


const AUTH_URL = "http://localhost:8080"
const AUTH_SIGN_IN = `${AUTH_URL}/public/sign_in`

export async function signIn(email: string, password: string): Promise<Boolean> {
    const response = await axios.post<SignInResponse>(AUTH_SIGN_IN, {
        email: email,
        password: password
    })

    if (response.status == 200) {
        saveToken({accessToken: response.data.token, refreshToken: response.data.refresh_token})
        return true
    } else {
        return false
    }
}

export function refreshToken() {

}

function getToken(): Token | null {
    const tokenString = sessionStorage.getItem("token")
    if (tokenString != null) {
        return JSON.parse(tokenString)
    } else {
        return null
    }
}

function saveToken(token: Token) {
    const tokenString = JSON.stringify(token)
    sessionStorage.setItem("token", tokenString)
    axios.defaults.headers.common["Authorization"] = `Bearer ${token.accessToken}`;
}