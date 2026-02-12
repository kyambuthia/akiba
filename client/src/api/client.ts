export const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1';

type AuthPayload = {
  email?: string;
  phone?: string;
  username?: string;
  login?: string;
  password: string;
};

export async function signup(payload: Required<Pick<AuthPayload, 'email' | 'phone' | 'username' | 'password'>>) {
  const res = await fetch(`${API_BASE_URL}/auth/signup`, {
    method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload)
  });
  return res.json();
}

export async function login(payload: Required<Pick<AuthPayload, 'login' | 'password'>>) {
  const res = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload)
  });
  return res.json();
}
