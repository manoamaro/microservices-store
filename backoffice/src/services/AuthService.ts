import {ApiClient} from "./ApiClient";

const AUTH_URL = "http://localhost:8080"

type SignInResponse = {
    token: string;
    refresh_token: string;
}

type Token = {
    accessToken: string;
    refreshToken: string;
}

const TOKEN_KEY = "token";

export function saveToken(token: Token | null) {
    if (token != null) {
        const tokenString = JSON.stringify(token)
        sessionStorage.setItem(TOKEN_KEY, tokenString)
    } else {
        sessionStorage.removeItem(TOKEN_KEY)
    }
}

export async function getToken(): Promise<Token> {
    const tokenString = sessionStorage.getItem(TOKEN_KEY)
    if (tokenString != null) {
        return JSON.parse(tokenString);
    } else {
        return Promise.reject("No token");
    }
}

export async function isAuthenticated(): Promise<boolean> {
    return getToken().then(() => true).catch(() => false)
}

export async function logout() {
    saveToken(null);
}


class AuthService extends ApiClient {
    constructor() {
        super({
            baseUrl: AUTH_URL
        });
    }

    async signIn(email: string, password: string): Promise<string> {
        const response = await this.request<SignInResponse>("POST", "/public/sign_in", {
            body: {
                email: email,
                password: password
            }
        });
        saveToken({
            accessToken: response.token,
            refreshToken: response.refresh_token
        });
        return response.token;
    }

    async refresh(): Promise<string> {
        const token = await getToken();
        const response = await this.request<SignInResponse>("POST", "/public/refresh", {
            body: {
                refresh_token: token.refreshToken
            }
        });
        saveToken({accessToken: response.token, refreshToken: response.refresh_token});
        return response.refresh_token;
    }
}

export default new AuthService();