import {ApiClient} from "./ApiClient";
import AuthService, {getToken} from "./AuthService";
import {number} from "prop-types";

// interface Inventory with amount number field
export interface Inventory {
    amount: number
}

export interface InventoryRequest {
    product_id: string,
    amount: number
}

class InventoryService extends ApiClient {

    constructor() {
        super({
            baseUrl: "http://localhost:8081",
            getAccessToken: () => getToken().then(value => value.accessToken),
            onRefreshToken: () => AuthService.refresh()
        });
    }

    async getInventory(productId: string): Promise<Inventory> {
        return await this.request<Inventory>("GET", `/public/inventory/${productId}`);
    }

    async putInventory(productId: string, amount: number): Promise<Inventory> {
        return await this.request<Inventory>("POST", `/internal/inventory/set`, {
            body: {
                product_id: productId,
                amount: +amount,
            }
        });
    }
}

export default new InventoryService();