import {ApiClient} from "./ApiClient";
import AuthService, {getToken} from "./AuthService";

const PRODUCTS_HOST = "http://localhost:8082";

export interface Product {
    id: string | null,
    name: string,
    description: string,
    price: number,
}

class ProductService extends ApiClient {

    constructor() {
        super({
            baseUrl: PRODUCTS_HOST,
            getAccessToken: () => getToken().then(value => value.accessToken),
            onRefreshToken: () => AuthService.refresh()
        });
    }

    async getProducts(): Promise<Product[]> {
        return await this.request<Product[]>("GET", "/admin/")
    }

    async getProduct(id: string): Promise<Product> {
        return await this.request<Product>("GET", `/admin/${id}`)
    }

    async postProduct(product: Product): Promise<Product> {
        return await this.request<Product>("POST", `/admin/`, {body: product});
    }

    async putProduct(product: Product): Promise<Product> {
        return await this.request<Product>("PUT", `/admin/${product.id}`, {body: product});
    }
}

export default new ProductService();