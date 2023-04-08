const URL = 'http://localhost:8081';

export async function getCart() {
    const response = await fetch(`${URL}/public/`);
    if (response.ok) {
        return await response.json();
    } else {
        return null;
    }
}
