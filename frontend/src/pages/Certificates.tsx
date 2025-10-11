import React from "react";

import { useEffect, useMemo, useState } from "react";
import { api } from "../api";

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

  const canSubmit = useMemo(() => !!form.domain && !submitting, [form.domain, submitting]);

  const load = () => {
    setLoading(true);
    setError("");
    return api
      .listCertificates()
      .then((d: any) => setList(d.certificates || []))
      .catch((e: any) => setError(String(e)))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load();
  }, []);

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
    } catch (e: any) {
      setError(String(e));
    } finally {
      setSubmitting(false);
    }
  };

  const del = async (domain: string) => {
    if (!confirm(`Delete certificate for ${domain}?`)) return;
    setError("");
    try {
      await api.deleteCertificate(domain);
      await load();
    } catch (e: any) {
      setError(String(e));
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
              <th>Domain</th>
              <th>Issuer</th>
              <th>Not Before</th>
              <th>Not After</th>
              <th>Valid</th>
              <th>Wildcard</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {list.map((c) => (
              <tr key={c.domain}>
                <td>{c.domain}</td>
                <td>{c.issuer}</td>
                <td>{fmt(c.not_before)}</td>
                <td>{fmt(c.not_after)}</td>
                <td>{String(c.is_valid)}</td>
                <td>{String(c.is_wildcard)}</td>
                <td>
                  <button onClick={() => del(c.domain)} disabled={loading || submitting}>
                    Delete
                  </button>
                </td>
              </tr>
            ))}
            {!loading && !list.length && (
              <tr>
                <td colSpan={7} style={{ textAlign: "center" }}>
                  No certificates
                </td>
              </tr>
            )}
            {loading && (
              <tr>
                <td colSpan={7} style={{ textAlign: "center" }}>
                  Loading...
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </section>
    </div>
  );
}
