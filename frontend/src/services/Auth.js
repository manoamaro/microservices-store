const URL = 'http://localhost:8080';

export async function signUp(email, password) {

    const response = await fetch(`${URL}/public/sign_up`);
    if (response.ok) {
        return await response.json();
    } else {
        return null;
    }
}
