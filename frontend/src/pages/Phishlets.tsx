import { useEffect, useState } from "react";
import { api } from "../api";

export default function Phishlets() {
  const [list, setList] = useState<any[]>([]);
  const [form, setForm] = useState<any>({ name: "", display_name: "" });
  const [error, setError] = useState("");

  const load = () =>
    api
      .listPhishlets()
      .then((r: any) => setList(r.phishlets || r))
      .catch((e) => setError(String(e)));

  useEffect(() => {
    load();
  }, []);

  const create = async () => {
    try {
      await api.createPhishlet(form);
      setForm({ name: "", display_name: "" });
      load();
    } catch (e: any) {
      setError(String(e));
    }
  };

  const remove = async (name: string) => {
    try {
      await api.deletePhishlet(name);
      load();
    } catch (e: any) {
      setError(String(e));
    }
  };

  return (
    <div>
      <h2>Phishlets</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      <div style={{ margin: "8px 0" }}>
        <input
          placeholder="name"
          value={form.name}
          onChange={(e) => setForm({ ...form, name: e.target.value })}
        />
        <input
          placeholder="display_name"
          value={form.display_name}
          onChange={(e) => setForm({ ...form, display_name: e.target.value })}
        />
        <button onClick={create}>Create</button>
      </div>
      <table>
        <thead>
          <tr><th>Name</th><th>Display</th><th>Actions</th></tr>
        </thead>
        <tbody>
          {list.map((p: any) => (
            <tr key={p.name || p.id}>
              <td>{p.name}</td>
              <td>{p.display_name}</td>
              <td>
                <button onClick={() => remove(p.name)}>Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
