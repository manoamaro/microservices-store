import {ApiClient} from "./ApiClient";
import AuthService, {getToken} from "./AuthService";

export const PRODUCTS_HOST = "http://localhost:8082";

export interface Product {
    id: string | null,
    name: string,
    description: string,
    prices: ProductPrice[],
    images: ProductImage[],
}

export interface ProductPrice {
    currency: string,
    price: number,
}

export interface ProductImage {
    id: string,
    url: string,
    description: string
}

export const EMPTY_PRODUCT: Product = {
    id: "",
    name: "",
    description: "",
    prices: [],
    images: []
};

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

    async postProductImages(id: string, files: FileList): Promise<Product> {
        const formData = new FormData();
        for (let i = 0; i < files.length; i++) {
            formData.append("images", files[i]);
        }
        return await this.request<Product>("POST", `/admin/${id}/upload`, {body: formData});
    }

    async deleteProductImage(id: string, imageId: string): Promise<Product> {
        return await this.request<Product>("DELETE", `/admin/${id}/image/${imageId}`);
    }

}

export default new ProductService();