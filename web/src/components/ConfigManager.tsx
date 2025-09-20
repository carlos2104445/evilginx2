import React, { useState, useEffect } from 'react';

interface Config {
  server: {
    port: number;
    hostname: string;
    httpsPort: number;
    enableTLS: boolean;
  };
  stealth: {
    enableTrafficFiltering: boolean;
    enableDomainFronting: boolean;
    enableObfuscation: boolean;
    enableSandboxDetection: boolean;
    rotationInterval: number;
    maxFrontDomains: number;
  };
  geolocation: {
    enableGeoFiltering: boolean;
    allowedCountries: string[];
    blockVPN: boolean;
    blockTor: boolean;
    blockCloudProviders: boolean;
  };
  logging: {
    level: string;
    enableFileLogging: boolean;
    logPath: string;
    enableMetrics: boolean;
  };
}

const ConfigManager: React.FC = () => {
  const [config, setConfig] = useState<Config>({
    server: {
      port: 8080,
      hostname: 'localhost',
      httpsPort: 8443,
      enableTLS: false,
    },
    stealth: {
      enableTrafficFiltering: true,
      enableDomainFronting: false,
      enableObfuscation: true,
      enableSandboxDetection: true,
      rotationInterval: 3600,
      maxFrontDomains: 10,
    },
    geolocation: {
      enableGeoFiltering: false,
      allowedCountries: [],
      blockVPN: true,
      blockTor: true,
      blockCloudProviders: false,
    },
    logging: {
      level: 'info',
      enableFileLogging: true,
      logPath: '/var/log/evilginx2.log',
      enableMetrics: true,
    },
  });
  
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [saveMessage, setSaveMessage] = useState<string | null>(null);

  useEffect(() => {
    const fetchConfig = async () => {
      try {
        setLoading(true);
        const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
        const response = await fetch(`${apiUrl}/api/config`);
        
        if (response.ok) {
          const data = await response.json();
          setConfig(data);
        } else {
          setError('Failed to fetch configuration');
        }
      } catch (err) {
        setError('Failed to fetch configuration');
        console.error('Config fetch error:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchConfig();
  }, []);

  const saveConfig = async () => {
    try {
      setSaving(true);
      setSaveMessage(null);
      
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/config`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(config)
      });

      if (response.ok) {
        setSaveMessage('Configuration saved successfully');
      } else {
        setSaveMessage('Failed to save configuration');
      }
    } catch (err) {
      setSaveMessage('Failed to save configuration');
      console.error('Config save error:', err);
    } finally {
      setSaving(false);
      setTimeout(() => setSaveMessage(null), 3000);
    }
  };

  const updateConfig = (section: keyof Config, field: string, value: any) => {
    setConfig(prev => ({
      ...prev,
      [section]: {
        ...prev[section],
        [field]: value
      }
    }));
  };

  if (loading) {
    return <div className="loading">Loading configuration...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  return (
    <div className="config-manager">
      <div className="card">
        <div className="card-header">
          <h2 className="card-title">Configuration Management</h2>
          <button 
            className="btn btn-primary" 
            onClick={saveConfig}
            disabled={saving}
          >
            {saving ? 'Saving...' : 'Save Configuration'}
          </button>
        </div>

        {saveMessage && (
          <div className={`alert ${saveMessage.includes('success') ? 'alert-success' : 'alert-danger'}`}>
            {saveMessage}
          </div>
        )}

        <div className="config-sections">
          <div className="config-section">
            <h3>Server Configuration</h3>
            <div className="form-group">
              <label className="form-label">Hostname</label>
              <input
                type="text"
                className="form-control"
                value={config.server.hostname}
                onChange={(e) => updateConfig('server', 'hostname', e.target.value)}
              />
            </div>
            <div className="form-group">
              <label className="form-label">HTTP Port</label>
              <input
                type="number"
                className="form-control"
                value={config.server.port}
                onChange={(e) => updateConfig('server', 'port', parseInt(e.target.value))}
              />
            </div>
            <div className="form-group">
              <label className="form-label">HTTPS Port</label>
              <input
                type="number"
                className="form-control"
                value={config.server.httpsPort}
                onChange={(e) => updateConfig('server', 'httpsPort', parseInt(e.target.value))}
              />
            </div>
            <div className="form-group">
              <label className="form-label">
                <input
                  type="checkbox"
                  checked={config.server.enableTLS}
                  onChange={(e) => updateConfig('server', 'enableTLS', e.target.checked)}
                />
                Enable TLS
              </label>
            </div>
          </div>

          <div className="config-section">
            <h3>Stealth & Evasion</h3>
            <div className="form-group">
              <label className="form-label">
                <input
                  type="checkbox"
                  checked={config.stealth.enableTrafficFiltering}
                  onChange={(e) => updateConfig('stealth', 'enableTrafficFiltering', e.target.checked)}
                />
                Enable Traffic Filtering
              </label>
            </div>
            <div className="form-group">
              <label className="form-label">
                <input
                  type="checkbox"
                  checked={config.stealth.enableDomainFronting}
                  onChange={(e) => updateConfig('stealth', 'enableDomainFronting', e.target.checked)}
                />
                Enable Domain Fronting
              </label>
            </div>
            <div className="form-group">
              <label className="form-label">
                <input
                  type="checkbox"
                  checked={config.stealth.enableObfuscation}
                  onChange={(e) => updateConfig('stealth', 'enableObfuscation', e.target.checked)}
                />
                Enable Content Obfuscation
              </label>
            </div>
            <div className="form-group">
              <label className="form-label">
                <input
                  type="checkbox"
                  checked={config.stealth.enableSandboxDetection}
                  onChange={(e) => updateConfig('stealth', 'enableSandboxDetection', e.target.checked)}
                />
                Enable Sandbox Detection
              </label>
            </div>
            <div className="form-group">
              <label className="form-label">Domain Rotation Interval (seconds)</label>
              <input
                type="number"
                className="form-control"
                value={config.stealth.rotationInterval}
                onChange={(e) => updateConfig('stealth', 'rotationInterval', parseInt(e.target.value))}
              />
            </div>
            <div className="form-group">
              <label className="form-label">Max Front Domains</label>
              <input
                type="number"
                className="form-control"
                value={config.stealth.maxFrontDomains}
                onChange={(e) => updateConfig('stealth', 'maxFrontDomains', parseInt(e.target.value))}
              />
            </div>
          </div>

          <div className="config-section">
            <h3>Geolocation Filtering</h3>
            <div className="form-group">
              <label className="form-label">
                <input
                  type="checkbox"
                  checked={config.geolocation.enableGeoFiltering}
                  onChange={(e) => updateConfig('geolocation', 'enableGeoFiltering', e.target.checked)}
                />
                Enable Geo-filtering
              </label>
            </div>
            <div className="form-group">
              <label className="form-label">Allowed Countries (comma-separated)</label>
              <input
                type="text"
                className="form-control"
                value={config.geolocation.allowedCountries.join(', ')}
                onChange={(e) => updateConfig('geolocation', 'allowedCountries', 
                  e.target.value.split(',').map(c => c.trim()).filter(c => c)
                )}
                placeholder="US, GB, CA"
              />
            </div>
            <div className="form-group">
              <label className="form-label">
                <input
                  type="checkbox"
                  checked={config.geolocation.blockVPN}
                  onChange={(e) => updateConfig('geolocation', 'blockVPN', e.target.checked)}
                />
                Block VPN Traffic
              </label>
            </div>
            <div className="form-group">
              <label className="form-label">
                <input
                  type="checkbox"
                  checked={config.geolocation.blockTor}
                  onChange={(e) => updateConfig('geolocation', 'blockTor', e.target.checked)}
                />
                Block Tor Traffic
              </label>
            </div>
            <div className="form-group">
              <label className="form-label">
                <input
                  type="checkbox"
                  checked={config.geolocation.blockCloudProviders}
                  onChange={(e) => updateConfig('geolocation', 'blockCloudProviders', e.target.checked)}
                />
                Block Cloud Provider IPs
              </label>
            </div>
          </div>

          <div className="config-section">
            <h3>Logging & Monitoring</h3>
            <div className="form-group">
              <label className="form-label">Log Level</label>
              <select
                className="form-control"
                value={config.logging.level}
                onChange={(e) => updateConfig('logging', 'level', e.target.value)}
              >
                <option value="debug">Debug</option>
                <option value="info">Info</option>
                <option value="warn">Warning</option>
                <option value="error">Error</option>
              </select>
            </div>
            <div className="form-group">
              <label className="form-label">
                <input
                  type="checkbox"
                  checked={config.logging.enableFileLogging}
                  onChange={(e) => updateConfig('logging', 'enableFileLogging', e.target.checked)}
                />
                Enable File Logging
              </label>
            </div>
            <div className="form-group">
              <label className="form-label">Log File Path</label>
              <input
                type="text"
                className="form-control"
                value={config.logging.logPath}
                onChange={(e) => updateConfig('logging', 'logPath', e.target.value)}
              />
            </div>
            <div className="form-group">
              <label className="form-label">
                <input
                  type="checkbox"
                  checked={config.logging.enableMetrics}
                  onChange={(e) => updateConfig('logging', 'enableMetrics', e.target.checked)}
                />
                Enable Prometheus Metrics
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ConfigManager;
