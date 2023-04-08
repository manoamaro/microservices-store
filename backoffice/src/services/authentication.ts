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
const AUTH_REFRESH_TOKEN = `${AUTH_URL}/public/refresh`

export async function signIn(email: string, password: string) {
    try {
        const response = await axios.post<SignInResponse>(AUTH_SIGN_IN, {
            email: email,
            password: password
        })
        saveToken({accessToken: response.data.token, refreshToken: response.data.refresh_token})
    } catch (e) {
        console.log(e)
        saveToken(null)
        throw e
    }
}

export async function refreshToken() {
    const token = getToken();
    if (token != null) {
        try {
            const response = await axios.post<SignInResponse>(AUTH_REFRESH_TOKEN, {
                refresh_token: token.refreshToken
            })
            if (response.status == 200) {
                saveToken({accessToken: response.data.token, refreshToken: response.data.refresh_token})
            } else {
                saveToken(null)
            }
        } catch (e) {
            saveToken(null)
            throw e
        }
    } else {
        throw 'no token';
    }
}

function getToken(): Token | null {
    const tokenString = sessionStorage.getItem("token")
    if (tokenString != null) {
        return JSON.parse(tokenString)
    } else {
        return null
    }
}

export async function isAuthenticated(): Promise<boolean> {
    return getToken() != null
}

function saveToken(token: Token | null) {
    if (token != null) {
        const tokenString = JSON.stringify(token)
        sessionStorage.setItem("token", tokenString)
        axios.defaults.headers.common["Authorization"] = `Bearer ${token.accessToken}`;
    } else {
        sessionStorage.removeItem("token")
        delete axios.defaults.headers.common["Authorization"];
    }
}