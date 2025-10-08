import { useState } from "react";
import { setToken } from "../api";
import { useNavigate } from "react-router-dom";

export default function Login() {
  const [token, setTok] = useState("");
  const nav = useNavigate();

  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;
    setToken(token);
    nav("/dashboard", { replace: true });
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
        <button type="submit">Enter</button>
      </form>
    </div>
  );
}
