import { useEffect, useState } from "react";
import { api } from "../api";

export default function Config() {
  const [cfg, setCfg] = useState<any>(null);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);

  const load = () =>
    api.getConfig().then(setCfg).catch((e) => setError(String(e)));

  useEffect(() => {
    load();
  }, []);

  const save = async () => {
    setSaving(true);
    try {
      await api.updateConfig(cfg);
      load();
    } catch (e: any) {
      setError(String(e));
    } finally {
      setSaving(false);
    }
  };

  return (
    <div>
      <h2>Config</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      {cfg && (
        <>
          <textarea
            style={{ width: "100%", height: 300 }}
            value={JSON.stringify(cfg, null, 2)}
            onChange={(e) => setCfg(JSON.parse(e.target.value || "{}"))}
          />
          <div style={{ marginTop: 8 }}>
            <button onClick={save} disabled={saving}>
              {saving ? "Saving..." : "Save"}
            </button>
          </div>
        </>
      )}
    </div>
  );
}
