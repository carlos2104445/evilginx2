import { Outlet, Link, useNavigate, useLocation } from "react-router-dom";
import { getToken, setToken } from "./api";

export default function App() {
  const nav = useNavigate();
  const loc = useLocation();
  const authed = !!getToken();

  if (!authed && loc.pathname !== "/login") {
    nav("/login", { replace: true });
  }

  const logout = () => {
    setToken("");
    nav("/login", { replace: true });
  };

  return (
    <div style={{ display: "flex", minHeight: "100vh" }}>
      <nav style={{ width: 220, padding: 16, borderRight: "1px solid #ddd" }}>
        <h3>Evilginx2</h3>
        <ul style={{ listStyle: "none", padding: 0 }}>
          <li><Link to="/dashboard">Dashboard</Link></li>
          <li><Link to="/phishlets">Phishlets</Link></li>
          <li><Link to="/sessions">Sessions</Link></li>
          <li><Link to="/lures">Lures</Link></li>
          <li><Link to="/config">Config</Link></li>
        </ul>
        {authed && <button onClick={logout}>Logout</button>}
      </nav>
      <main style={{ flex: 1, padding: 16 }}>
        <Outlet />
      </main>
    </div>
  );
}
