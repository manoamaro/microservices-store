type HttpMethod = "GET" | "POST" | "PUT" | "DELETE";

interface RequestOptions extends RequestInit {
    body?: any;
}

interface RequestConfig {
    baseUrl: string;
    getAccessToken?: () => Promise<string>;
    onRefreshToken?: () => Promise<string>;
}

export class ApiClient {
    private readonly baseUrl: string;
    private readonly getAccessToken?: () => Promise<string>;
    private readonly onRefreshToken?: () => Promise<string>;

    constructor(config: RequestConfig) {
        this.baseUrl = config.baseUrl;
        this.getAccessToken = config.getAccessToken;
        this.onRefreshToken = config.onRefreshToken;
    }

    async request<T>(
        method: HttpMethod,
        path: string,
        options?: RequestOptions,
        refreshAttempts = 1
    ): Promise<T> {
        const url = `${this.baseUrl}${path}`;
        let headers: HeadersInit = {
            "Accept": "application/json",
        };

        if (this.getAccessToken) {
            const accessToken = await this.getAccessToken();
            headers = {
                ...headers,
                Authorization: `Bearer ${accessToken}`
            }
        }

        if (options && options.body) {
            if (options.body instanceof FormData) {
                console.log("FormData");
            } else if (Array.isArray(options.body) || options.body instanceof Object) {
                headers = {
                    ...headers,
                    'Content-Type': 'application/json'
                }
                options.body = JSON.stringify(options.body);
            }
        }

        try {
            const response = await fetch(url, {
                method,
                headers,
                ...options,
            });

            if (response.status === 401 && this.onRefreshToken && refreshAttempts > 0) {
                await this.onRefreshToken();
                return await this.request<T>(method, path, options, refreshAttempts - 1);
            } else if (response.ok) {
                return await response.json();
            } else {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
        } catch (e) {
            throw e;
        }
    }
}
