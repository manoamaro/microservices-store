const URL = 'http://localhost:8082';

export async function getProducts() {
    const response = await fetch(`${URL}/public/`);
    if (response.ok) {
        return await response.json();
    } else {
        return null;
    }
}

export async function getProduct(id) {
    const response = await fetch(`${URL}/public/${id}`);
    if (response.ok) {
        return await response.json();
    } else {
        return null;
    }
}