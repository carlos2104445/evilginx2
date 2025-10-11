import { useEffect, useMemo, useState } from "react";
import { api, clearToken } from "../api";
import { useToast } from "../components/Toast";
import { useNavigate } from "react-router-dom";
import ConfirmDialog from "../components/ConfirmDialog";
import Pagination from "../components/Pagination";
import TableSort from "../components/TableSort";

export default function Sessions() {
  const [list, setList] = useState<any[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [toDelete, setToDelete] = useState<string | null>(null);

  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);
  const [sortKey, setSortKey] = useState<"id" | "phishlet_name" | "username" | "is_active">("id");
  const [sortDir, setSortDir] = useState<"asc" | "desc" | null>(null);

  const { show } = useToast();
  const nav = useNavigate();

  const load = async () => {
    setLoading(true);
    setError("");
    try {
      const r: any = await api.listSessions();
      setList(r.sessions || r);
    } catch (e: any) {
      if (e?.message === "UNAUTHORIZED") {
        clearToken();
        nav("/login");
        return;
      }
      setError(String(e));
      show("Failed to load sessions", "error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const toggleSort = (key: "id" | "phishlet_name" | "username" | "is_active") => {
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
      const va = (a[sortKey] ?? "").toString().toLowerCase();
      const vb = (b[sortKey] ?? "").toString().toLowerCase();
      return sortDir === "asc" ? va.localeCompare(vb) : vb.localeCompare(va);
    });
    return arr;
  }, [list, sortKey, sortDir]);

  const total = sorted.length;
  const start = (page - 1) * perPage;
  const items = sorted.slice(start, start + perPage);

  const askDelete = (id: string) => {
    setToDelete(id);
    setConfirmOpen(true);
  };

  const doDelete = async () => {
    if (!toDelete) return;
    setConfirmOpen(false);
    try {
      await api.deleteSession(toDelete);
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
      <h2>Sessions</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
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
              label="Phishlet"
              active={sortKey === "phishlet_name"}
              direction={sortKey === "phishlet_name" ? sortDir : null}
              onToggle={() => toggleSort("phishlet_name")}
            />
            <TableSort
              label="User"
              active={sortKey === "username"}
              direction={sortKey === "username" ? sortDir : null}
              onToggle={() => toggleSort("username")}
            />
            <TableSort
              label="Active"
              active={sortKey === "is_active"}
              direction={sortKey === "is_active" ? sortDir : null}
              onToggle={() => toggleSort("is_active")}
            />
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {items.map((s: any) => (
            <tr key={s.id}>
              <td>{s.id}</td>
              <td>{s.phishlet_name}</td>
              <td>{s.username}</td>
              <td>{String(!!s.is_active)}</td>
              <td>
                <button onClick={() => askDelete(s.id)} disabled={loading}>
                  Delete
                </button>
              </td>
            </tr>
          ))}
          {!loading && !items.length && (
            <tr>
              <td colSpan={5} style={{ textAlign: "center" }}>
                No sessions
              </td>
            </tr>
          )}
          {loading && (
            <tr>
              <td colSpan={5} style={{ textAlign: "center" }}>
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
        message={`Delete session ${toDelete}?`}
        onCancel={() => setConfirmOpen(false)}
        onConfirm={doDelete}
      />
    </div>
  );
}
