import React from "react";

export default function TableSort({
  label,
  active,
  direction,
  onToggle,
}: {
  label: string;
  active: boolean;
  direction: "asc" | "desc" | null;
  onToggle: () => void;
}) {
  return (
    <th onClick={onToggle} style={{ cursor: "pointer", userSelect: "none" }}>
      {label} {active ? (direction === "asc" ? "▲" : direction === "desc" ? "▼" : "") : ""}
    </th>
  );
}
