import React from "react";

import { useEffect, useMemo, useState } from "react";
import { api, clearToken } from "../api";
import { useToast } from "../components/Toast";
import { useNavigate } from "react-router-dom";
import ConfirmDialog from "../components/ConfirmDialog";
import Pagination from "../components/Pagination";
import TableSort from "../components/TableSort";

function fmt(t?: string) {
  if (!t) return "";
  try {
    const d = new Date(t);
    return isNaN(d.getTime()) ? t : d.toISOString().replace(".000Z", "Z");
  } catch {
    return t;
  }
}

export default function Certificates() {
  const [list, setList] = useState<any[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form, setForm] = useState({
    domain: "",
    issuer: "",
    not_before: "",
    not_after: "",
    is_wildcard: false,
  });

  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);
  const [sortKey, setSortKey] = useState<"domain" | "issuer" | "not_after">("domain");
  const [sortDir, setSortDir] = useState<"asc" | "desc" | null>(null);

  const canSubmit = useMemo(() => !!form.domain && !submitting, [form.domain, submitting]);

  const { show } = useToast();
  const nav = useNavigate();

  const load = () => {
    setLoading(true);
    setError("");
    return api
      .listCertificates()
      .then((d: any) => setList(d.certificates || []))
      .catch((e: any) => {
        if (e?.message === "UNAUTHORIZED") {
          clearToken();
          nav("/login");
          return;
        }
        setError(String(e));
        show("Failed to load certificates", "error");
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load();
  }, []);

  const toggleSort = (key: "domain" | "issuer" | "not_after") => {
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

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.domain) {
      setError("Domain is required");
      return;
    }
    setSubmitting(true);
    setError("");
    try {
      await api.createCertificate(form);
      setForm({ domain: "", issuer: "", not_before: "", not_after: "", is_wildcard: false });
      await load();
      show("Certificate created", "success");
    } catch (e: any) {
      if (e?.message === "UNAUTHORIZED") {
        clearToken();
        nav("/login");
        return;
      }
      setError(String(e));
      show("Create failed", "error");
    } finally {
      setSubmitting(false);
    }
  };

  const [confirmOpen, setConfirmOpen] = useState(false);
  const [toDelete, setToDelete] = useState<string | null>(null);

  const askDelete = (domain: string) => {
    setToDelete(domain);
    setConfirmOpen(true);
  };

  const doDelete = async () => {
    if (!toDelete) return;
    setError("");
    setConfirmOpen(false);
    try {
      await api.deleteCertificate(toDelete);
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
      <h2>Certificates</h2>
      {error && <div style={{ color: "red", marginBottom: 8 }}>{error}</div>}

      <section>
        <h3>Create</h3>
        <form onSubmit={submit} style={{ display: "grid", gap: 8, maxWidth: 520 }}>
          <label>
            Domain
            <input
              value={form.domain}
              onChange={(e) => setForm({ ...form, domain: e.target.value })}
              placeholder="example.com"
            />
          </label>
          <label>
            Issuer
            <input
              value={form.issuer}
              onChange={(e) => setForm({ ...form, issuer: e.target.value })}
              placeholder="Let's Encrypt"
            />
          </label>
          <label>
            Not Before (RFC3339)
            <input
              value={form.not_before}
              onChange={(e) => setForm({ ...form, not_before: e.target.value })}
              placeholder="2025-01-01T00:00:00Z"
            />
          </label>
          <label>
            Not After (RFC3339)
            <input
              value={form.not_after}
              onChange={(e) => setForm({ ...form, not_after: e.target.value })}
              placeholder="2026-01-01T00:00:00Z"
            />
          </label>
          <label>
            <input
              type="checkbox"
              checked={form.is_wildcard}
              onChange={(e) => setForm({ ...form, is_wildcard: e.target.checked })}
            />{" "}
            Wildcard
          </label>
          <button type="submit" disabled={!canSubmit}>
            {submitting ? "Creating..." : "Create"}
          </button>
        </form>
      </section>

      <section>
        <h3 style={{ display: "flex", alignItems: "center", gap: 8 }}>
          Existing
          <button onClick={load} disabled={loading}>
            {loading ? "Refreshing..." : "Refresh"}
          </button>
        </h3>
        <table border={1} cellPadding={6} style={{ borderCollapse: "collapse", minWidth: 720 }}>
          <thead>
            <tr>
              <TableSort
                label="Domain"
                active={sortKey === "domain"}
                direction={sortKey === "domain" ? sortDir : null}
                onToggle={() => toggleSort("domain")}
              />
              <TableSort
                label="Issuer"
                active={sortKey === "issuer"}
                direction={sortKey === "issuer" ? sortDir : null}
                onToggle={() => toggleSort("issuer")}
              />
              <TableSort
                label="Not After"
                active={sortKey === "not_after"}
                direction={sortKey === "not_after" ? sortDir : null}
                onToggle={() => toggleSort("not_after")}
              />
              <th>Valid</th>
              <th>Wildcard</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {items.map((c) => (
              <tr key={c.domain}>
                <td>{c.domain}</td>
                <td>{c.issuer}</td>
                <td>{fmt(c.not_after)}</td>
                <td>{String(c.is_valid)}</td>
                <td>{String(c.is_wildcard)}</td>
                <td>
                  <button onClick={() => askDelete(c.domain)} disabled={loading || submitting}>
                    Delete
                  </button>
                </td>
              </tr>
            ))}
            {!loading && !items.length && (
              <tr>
                <td colSpan={6} style={{ textAlign: "center" }}>
                  No certificates
                </td>
              </tr>
            )}
            {loading && (
              <tr>
                <td colSpan={6} style={{ textAlign: "center" }}>
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
      </section>
      <ConfirmDialog
        open={confirmOpen}
        message={`Delete certificate for ${toDelete}?`}
        onCancel={() => setConfirmOpen(false)}
        onConfirm={doDelete}
      />
    </div>
  );
}
