export const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1';

type AuthPayload = {
  email?: string;
  phone?: string;
  username?: string;
  login?: string;
  password: string;
};

export type ApiError = {
  status: number;
  code: string;
  message: string;
  fields?: Record<string, string>;
};

async function parseJSONSafe(res: Response) {
  const text = await res.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return null;
  }
}

async function postJSON<T>(path: string, payload: unknown): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });

  const body = await parseJSONSafe(res);
  if (!res.ok) {
    const apiErr = (body as { error?: { code?: string; message?: string; fields?: Record<string, string> } } | null)?.error;
    const err: ApiError = {
      status: res.status,
      code: apiErr?.code || 'http_error',
      message: apiErr?.message || `Request failed with status ${res.status}`,
      fields: apiErr?.fields,
    };
    throw err;
  }
  return body as T;
}

export async function signup(payload: Required<Pick<AuthPayload, 'email' | 'phone' | 'username' | 'password'>>) {
  return postJSON('/auth/signup', payload);
}

export async function login(payload: Required<Pick<AuthPayload, 'login' | 'password'>>) {
  return postJSON('/auth/login', payload);
}
