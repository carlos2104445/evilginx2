import { useState } from "react";
import { setToken, clearToken } from "../api";
import { api } from "../api";
import { useNavigate } from "react-router-dom";

export default function Login() {
  const [token, setTok] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const nav = useNavigate();

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (!token.trim()) {
      setError("Token required");
      return;
    }
    setLoading(true);
    try {
      setToken(token.trim());
      await api.getConfig();
      nav("/dashboard", { replace: true });
    } catch {
      clearToken();
      setError("Invalid token");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 400 }}>
      <h2>Login</h2>
      <form onSubmit={submit}>
        <label>Admin Token</label>
        <input
          type="password"
          value={token}
          onChange={(e) => setTok(e.target.value)}
          style={{ width: "100%", margin: "8px 0" }}
        />
        {error && <div style={{ color: "#d64545", marginBottom: 8 }}>{error}</div>}
        <button type="submit" disabled={loading}>{loading ? "Validating…" : "Enter"}</button>
      </form>
    </div>
  );
}
