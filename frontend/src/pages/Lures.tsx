import { useEffect, useMemo, useState } from "react";
import { api, clearToken } from "../api";
import { useNavigate } from "react-router-dom";
import { useToast } from "../components/Toast";
import ConfirmDialog from "../components/ConfirmDialog";
import Pagination from "../components/Pagination";
import TableSort from "../components/TableSort";

export default function Lures() {
  const [list, setList] = useState<any[]>([]);
  const [form, setForm] = useState<any>({ id: "", hostname: "", path: "" });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [toDelete, setToDelete] = useState<string | null>(null);

  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);
  const [sortKey, setSortKey] = useState<"id" | "hostname" | "path">("id");
  const [sortDir, setSortDir] = useState<"asc" | "desc" | null>(null);

  const nav = useNavigate();
  const { show } = useToast();

  const load = async () => {
    setLoading(true);
    setError("");
    try {
      const r: any = await api.listLures();
      setList(r.lures || r);
    } catch (e: any) {
      if (e?.message === "UNAUTHORIZED") {
        clearToken();
        nav("/login");
        return;
      }
      setError(String(e));
      show("Failed to load lures", "error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const toggleSort = (key: "id" | "hostname" | "path") => {
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
    if (!form.id.trim() || !form.hostname.trim() || !form.path.trim()) {
      setError("id, hostname and path are required");
      return;
    }
    try {
      await api.createLure(form);
      setForm({ id: "", hostname: "", path: "" });
      await load();
      show("Lure created", "success");
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

  const askDelete = (id: string) => {
    setToDelete(id);
    setConfirmOpen(true);
  };

  const doDelete = async () => {
    if (!toDelete) return;
    setConfirmOpen(false);
    try {
      await api.deleteLure(toDelete);
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
      <h2>Lures</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      <div style={{ margin: "8px 0", display: "flex", gap: 8 }}>
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
        <button onClick={create} disabled={loading}>Create</button>
      </div>
      <table>
        <thead>
          <tr>
            <TableSort
              label="ID"
              active={sortKey === "id"}
              direction={sortKey === "id" ? sortDir : null}
              onToggle={() => toggleSort("id")}
            />
            <TableSort
              label="Hostname"
              active={sortKey === "hostname"}
              direction={sortKey === "hostname" ? sortDir : null}
              onToggle={() => toggleSort("hostname")}
            />
            <TableSort
              label="Path"
              active={sortKey === "path"}
              direction={sortKey === "path" ? sortDir : null}
              onToggle={() => toggleSort("path")}
            />
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {items.map((l: any) => (
            <tr key={l.id}>
              <td>{l.id}</td>
              <td>{l.hostname}</td>
              <td>{l.path}</td>
              <td><button onClick={() => askDelete(l.id)} disabled={loading}>Delete</button></td>
            </tr>
          ))}
          {!loading && !items.length && (
            <tr>
              <td colSpan={4} style={{ textAlign: "center" }}>
                No lures
              </td>
            </tr>
          )}
          {loading && (
            <tr>
              <td colSpan={4} style={{ textAlign: "center" }}>
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
        message={`Delete lure ${toDelete}?`}
        onCancel={() => setConfirmOpen(false)}
        onConfirm={doDelete}
      />
    </div>
  );
}
