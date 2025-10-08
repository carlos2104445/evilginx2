import { useEffect, useState } from "react";
import { api } from "../api";

export default function Sessions() {
  const [list, setList] = useState<any[]>([]);
  const [error, setError] = useState("");

  const load = () =>
    api
      .listSessions()
      .then((r: any) => setList(r.sessions || r))
      .catch((e) => setError(String(e)));

  useEffect(() => {
    load();
  }, []);

  const remove = async (id: string) => {
    try {
      await api.deleteSession(id);
      load();
    } catch (e: any) {
      setError(String(e));
    }
  };

  return (
    <div>
      <h2>Sessions</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      <table>
        <thead>
          <tr><th>ID</th><th>Phishlet</th><th>User</th><th>Active</th><th>Actions</th></tr>
        </thead>
        <tbody>
          {list.map((s: any) => (
            <tr key={s.id}>
              <td>{s.id}</td>
              <td>{s.phishlet_name}</td>
              <td>{s.username}</td>
              <td>{String(!!s.is_active)}</td>
              <td><button onClick={() => remove(s.id)}>Delete</button></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
