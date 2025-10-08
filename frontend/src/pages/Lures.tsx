import { useEffect, useState } from "react";
import { api } from "../api";

export default function Lures() {
  const [list, setList] = useState<any[]>([]);
  const [form, setForm] = useState<any>({ id: "", hostname: "", path: "" });
  const [error, setError] = useState("");

  const load = () =>
    api
      .listLures()
      .then((r: any) => setList(r.lures || r))
      .catch((e) => setError(String(e)));

  useEffect(() => {
    load();
  }, []);

  const create = async () => {
    try {
      await api.createLure(form);
      setForm({ id: "", hostname: "", path: "" });
      load();
    } catch (e: any) {
      setError(String(e));
    }
  };

  const remove = async (id: string) => {
    try {
      await api.deleteLure(id);
      load();
    } catch (e: any) {
      setError(String(e));
    }
  };

  return (
    <div>
      <h2>Lures</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      <div style={{ margin: "8px 0" }}>
        <input
          placeholder="id"
          value={form.id}
          onChange={(e) => setForm({ ...form, id: e.target.value })}
        />
        <input
          placeholder="hostname"
          value={form.hostname}
          onChange={(e) => setForm({ ...form, hostname: e.target.value })}
        />
        <input
          placeholder="path"
          value={form.path}
          onChange={(e) => setForm({ ...form, path: e.target.value })}
        />
        <button onClick={create}>Create</button>
      </div>
      <table>
        <thead>
          <tr><th>ID</th><th>Hostname</th><th>Path</th><th>Actions</th></tr>
        </thead>
        <tbody>
          {list.map((l: any) => (
            <tr key={l.id}>
              <td>{l.id}</td>
              <td>{l.hostname}</td>
              <td>{l.path}</td>
              <td><button onClick={() => remove(l.id)}>Delete</button></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
