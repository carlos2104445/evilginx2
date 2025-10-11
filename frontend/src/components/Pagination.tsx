import React from "react";

type Props = {
  page: number;
  perPage: number;
  total: number;
  onChange: (p: { page: number; perPage: number }) => void;
};

export default function Pagination({ page, perPage, total, onChange }: Props) {
  const pages = Math.max(1, Math.ceil(total / perPage));
  const prev = () => onChange({ page: Math.max(1, page - 1), perPage });
  const next = () => onChange({ page: Math.min(pages, page + 1), perPage });
  const setPP = (e: React.ChangeEvent<HTMLSelectElement>) =>
    onChange({ page: 1, perPage: Number(e.target.value) });

  return (
    <div style={{ display: "flex", alignItems: "center", gap: 8, marginTop: 8 }}>
      <button onClick={prev} disabled={page <= 1}>
        Prev
      </button>
      <span>
        Page {page} / {pages}
      </span>
      <button onClick={next} disabled={page >= pages}>
        Next
      </button>
      <label style={{ marginLeft: 8 }}>
        Per page{" "}
        <select value={perPage} onChange={setPP}>
          <option value={10}>10</option>
          <option value={25}>25</option>
          <option value={50}>50</option>
        </select>
      </label>
      <span style={{ marginLeft: "auto" }}>{total} total</span>
    </div>
  );
}
