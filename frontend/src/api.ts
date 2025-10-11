const BASE = (import.meta as any)?.env?.VITE_API_BASE || "http://localhost:8081/api/v1";

const TOKEN_KEY = "evil_admin_token";

export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) || "";
}

export function setToken(token: string) {
  if (token) localStorage.setItem(TOKEN_KEY, token);
  else localStorage.removeItem(TOKEN_KEY);
}

async function request(path: string, init: RequestInit = {}) {
  const headers = new Headers(init.headers || {});
  headers.set("Content-Type", "application/json");
  const token = getToken();
  if (token) headers.set("Authorization", `Bearer ${token}`);
  const res = await fetch(`${BASE}${path}`, { ...init, headers });
  if (res.status === 401) {
    throw new Error("Unauthorized");
  }
  if (!res.ok) {
    const msg = await res.text();
    throw new Error(msg || res.statusText);
  }
  if (res.status === 204) return null;
  return res.json();
}

export const api = {
  health: () => request("/health"),
  getConfig: () => request("/config"),
  updateConfig: (data: any) =>
    request("/config", { method: "PUT", body: JSON.stringify(data) }),

  listPhishlets: () => request("/phishlets"),
  getPhishlet: (name: string) => request(`/phishlets/${encodeURIComponent(name)}`),
  createPhishlet: (data: any) =>
    request("/phishlets", { method: "POST", body: JSON.stringify(data) }),
  updatePhishlet: (name: string, data: any) =>
    request(`/phishlets/${encodeURIComponent(name)}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),
  deletePhishlet: (name: string) =>
    request(`/phishlets/${encodeURIComponent(name)}`, { method: "DELETE" }),

  listSessions: () => request("/sessions"),
  getSession: (id: string) => request(`/sessions/${id}`),
  updateSession: (id: string, data: any) =>
    request(`/sessions/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  deleteSession: (id: string) =>
    request(`/sessions/${id}`, { method: "DELETE" }),
  sessionStats: () => request("/sessions/stats"),

  listLures: () => request("/lures"),
  getLure: (id: string) => request(`/lures/${id}`),
  createLure: (data: any) =>
    request("/lures", { method: "POST", body: JSON.stringify(data) }),
  updateLure: (id: string, data: any) =>
    request(`/lures/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  deleteLure: (id: string) => request(`/lures/${id}`, { method: "DELETE" }),
  
  listCertificates: () => request("/certificates"),
  createCertificate: (data: any) =>
    request("/certificates", { method: "POST", body: JSON.stringify(data) }),
  deleteCertificate: (domain: string) =>
    request(`/certificates/${encodeURIComponent(domain)}`, { method: "DELETE" }),
};
