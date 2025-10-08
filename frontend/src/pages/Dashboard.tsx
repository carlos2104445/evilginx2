import { useEffect, useState } from "react";
import { api } from "../api";

export default function Dashboard() {
  const [health, setHealth] = useState<any>(null);
  const [stats, setStats] = useState<any>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    api.health().then(setHealth).catch((e) => setError(String(e)));
    api.sessionStats().then(setStats).catch((e) => setError(String(e)));
  }, []);

  return (
    <div>
      <h2>Dashboard</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      <section>
        <h3>Health</h3>
        <pre>{JSON.stringify(health, null, 2)}</pre>
      </section>
      <section>
        <h3>Session Stats</h3>
        <pre>{JSON.stringify(stats, null, 2)}</pre>
      </section>
    </div>
  );
}
