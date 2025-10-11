import React from "react";

import { useEffect, useState } from "react";
import { api } from "../api";

export default function Certificates() {
  const [list, setList] = useState<any[]>([]);
  const [error, setError] = useState("");
  const [form, setForm] = useState({
    domain: "",
    issuer: "",
    not_before: "",
    not_after: "",
    is_wildcard: false,
  });

  const load = () =>
    api
      .listCertificates()
      .then((d: any) => setList(d.certificates || []))
      .catch((e: any) => setError(String(e)));

  useEffect(() => {
    load();
  }, []);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.domain) {
      setError("Domain is required");
      return;
    }
    try {
      await api.createCertificate(form);
      setForm({ domain: "", issuer: "", not_before: "", not_after: "", is_wildcard: false });
      await load();
    } catch (e: any) {
      setError(String(e));
    }
  };

  const del = async (domain: string) => {
    if (!confirm(`Delete certificate for ${domain}?`)) return;
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
      {error && <div style={{ color: "red" }}>{error}</div>}
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
          <button type="submit">Create</button>
        </form>
      </section>

      <section>
        <h3>Existing</h3>
        <table border={1} cellPadding={6} style={{ borderCollapse: "collapse", minWidth: 600 }}>
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
                <td>{c.not_before}</td>
                <td>{c.not_after}</td>
                <td>{String(c.is_valid)}</td>
                <td>{String(c.is_wildcard)}</td>
                <td>
                  <button onClick={() => del(c.domain)}>Delete</button>
                </td>
              </tr>
            ))}
            {!list.length && (
              <tr>
                <td colSpan={7} style={{ textAlign: "center" }}>
                  No certificates
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </section>
    </div>
  );
}
