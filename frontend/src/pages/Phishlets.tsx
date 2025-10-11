import { useEffect, useMemo, useState } from "react";
import { api, clearToken } from "../api";
import { useNavigate } from "react-router-dom";
import { useToast } from "../components/Toast";
import ConfirmDialog from "../components/ConfirmDialog";
import Pagination from "../components/Pagination";
import TableSort from "../components/TableSort";

export default function Phishlets() {
  const [list, setList] = useState<any[]>([]);
  const [form, setForm] = useState<any>({ name: "", display_name: "" });
  const [error, setError] = useState("");
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [toDelete, setToDelete] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);
  const [sortKey, setSortKey] = useState<"name" | "display_name">("name");
  const [sortDir, setSortDir] = useState<"asc" | "desc" | null>(null);

  const nav = useNavigate();
  const { show } = useToast();

  const load = async () => {
    setLoading(true);
    setError("");
    try {
      const r: any = await api.listPhishlets();
      setList(r.phishlets || r);
    } catch (e: any) {
      if (e?.message === "UNAUTHORIZED") {
        clearToken();
        nav("/login");
        return;
      }
      setError(String(e));
      show("Failed to load phishlets", "error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const toggleSort = (key: "name" | "display_name") => {
    if (sortKey !== key) {
      setSortKey(key);
      setSortDir("asc");
    } else {
      setSortDir(sortDir === "asc" ? "desc" : sortDir === "desc" ? null : "asc");
    }
  };

  const sorted = useMemo(() => {
    const arr = [...list];
    if (!sortDir) return arr;
    arr.sort((a, b) => {
      const va = (a[sortKey] || "").toString().toLowerCase();
      const vb = (b[sortKey] || "").toString().toLowerCase();
      return sortDir === "asc" ? va.localeCompare(vb) : vb.localeCompare(va);
    });
    return arr;
  }, [list, sortKey, sortDir]);

  const total = sorted.length;
  const start = (page - 1) * perPage;
  const items = sorted.slice(start, start + perPage);

  const create = async () => {
    setError("");
    if (!form.name.trim() || !form.display_name.trim()) {
      setError("name and display_name are required");
      return;
    }
    try {
      await api.createPhishlet(form);
      setForm({ name: "", display_name: "" });
      await load();
      show("Phishlet created", "success");
    } catch (e: any) {
      if (e?.message === "UNAUTHORIZED") {
        clearToken();
        nav("/login");
        return;
      }
      setError(String(e));
      show("Create failed", "error");
    }
  };

  const askDelete = (name: string) => {
    setToDelete(name);
    setConfirmOpen(true);
  };

  const doDelete = async () => {
    if (!toDelete) return;
    setConfirmOpen(false);
    try {
      await api.deletePhishlet(toDelete);
      await load();
      show("Deleted", "success");
    } catch (e: any) {
      if (e?.message === "UNAUTHORIZED") {
        clearToken();
        nav("/login");
        return;
      }
      setError(String(e));
      show("Delete failed", "error");
    } finally {
      setToDelete(null);
    }
  };

  return (
    <div>
      <h2>Phishlets</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      <div style={{ margin: "8px 0", display: "flex", gap: 8 }}>
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
        <button onClick={create} disabled={loading}>Create</button>
      </div>
      <table>
        <thead>
          <tr>
            <TableSort
              label="Name"
              active={sortKey === "name"}
              direction={sortKey === "name" ? sortDir : null}
              onToggle={() => toggleSort("name")}
            />
            <TableSort
              label="Display"
              active={sortKey === "display_name"}
              direction={sortKey === "display_name" ? sortDir : null}
              onToggle={() => toggleSort("display_name")}
            />
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {items.map((p: any) => (
            <tr key={p.name || p.id}>
              <td>{p.name}</td>
              <td>{p.display_name}</td>
              <td>
                <button onClick={() => askDelete(p.name)} disabled={loading}>Delete</button>
              </td>
            </tr>
          ))}
          {!loading && !items.length && (
            <tr>
              <td colSpan={3} style={{ textAlign: "center" }}>
                No phishlets
              </td>
            </tr>
          )}
          {loading && (
            <tr>
              <td colSpan={3} style={{ textAlign: "center" }}>
                Loading...
              </td>
            </tr>
          )}
        </tbody>
      </table>
      <Pagination
        page={page}
        perPage={perPage}
        total={total}
        onChange={({ page, perPage }) => {
          setPage(page);
          setPerPage(perPage);
        }}
      />
      <ConfirmDialog
        open={confirmOpen}
        message={`Delete ${toDelete}?`}
        onCancel={() => setConfirmOpen(false)}
        onConfirm={doDelete}
      />
    </div>
  );
}
