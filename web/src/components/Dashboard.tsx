import React, { useState, useEffect } from 'react';

interface DashboardStats {
  activeCampaigns: number;
  totalSessions: number;
  capturedCredentials: number;
  blockedRequests: number;
  successRate: number;
}

interface RecentSession {
  id: string;
  phishlet: string;
  ip: string;
  userAgent: string;
  timestamp: string;
  status: 'active' | 'completed' | 'blocked';
}

const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<DashboardStats>({
    activeCampaigns: 0,
    totalSessions: 0,
    capturedCredentials: 0,
    blockedRequests: 0,
    successRate: 0,
  });
  const [recentSessions, setRecentSessions] = useState<RecentSession[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        setLoading(true);
        const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
        
        const [statsResponse, sessionsResponse] = await Promise.all([
          fetch(`${apiUrl}/api/stats`),
          fetch(`${apiUrl}/api/sessions/recent`)
        ]);

        if (statsResponse.ok) {
          const statsData = await statsResponse.json();
          setStats(statsData);
        }

        if (sessionsResponse.ok) {
          const sessionsData = await sessionsResponse.json();
          setRecentSessions(sessionsData);
        }

        setError(null);
      } catch (err) {
        setError('Failed to fetch dashboard data');
        console.error('Dashboard fetch error:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchDashboardData();
    const interval = setInterval(fetchDashboardData, 30000);
    return () => clearInterval(interval);
  }, []);

  if (loading) {
    return <div className="loading">Loading dashboard...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  return (
    <div className="dashboard">
      <div className="card">
        <div className="card-header">
          <h2 className="card-title">Campaign Overview</h2>
        </div>
        
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-value">{stats.activeCampaigns}</div>
            <div className="stat-label">Active Campaigns</div>
          </div>
          <div className="stat-card">
            <div className="stat-value">{stats.totalSessions}</div>
            <div className="stat-label">Total Sessions</div>
          </div>
          <div className="stat-card">
            <div className="stat-value">{stats.capturedCredentials}</div>
            <div className="stat-label">Captured Credentials</div>
          </div>
          <div className="stat-card">
            <div className="stat-value">{stats.blockedRequests}</div>
            <div className="stat-label">Blocked Requests</div>
          </div>
          <div className="stat-card">
            <div className="stat-value">{stats.successRate.toFixed(1)}%</div>
            <div className="stat-label">Success Rate</div>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <h3 className="card-title">Recent Sessions</h3>
        </div>
        
        {recentSessions.length === 0 ? (
          <p>No recent sessions found.</p>
        ) : (
          <table className="table">
            <thead>
              <tr>
                <th>Session ID</th>
                <th>Phishlet</th>
                <th>IP Address</th>
                <th>User Agent</th>
                <th>Timestamp</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {recentSessions.map(session => (
                <tr key={session.id}>
                  <td>{session.id.substring(0, 8)}...</td>
                  <td>{session.phishlet}</td>
                  <td>{session.ip}</td>
                  <td title={session.userAgent}>
                    {session.userAgent.length > 50 
                      ? `${session.userAgent.substring(0, 50)}...` 
                      : session.userAgent}
                  </td>
                  <td>{new Date(session.timestamp).toLocaleString()}</td>
                  <td>
                    <span className={`status-badge status-${session.status}`}>
                      {session.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      <div className="card">
        <div className="card-header">
          <h3 className="card-title">System Health</h3>
        </div>
        
        <div className="health-indicators">
          <div className="health-item">
            <span className="health-label">Proxy Service:</span>
            <span className="health-status health-ok">Running</span>
          </div>
          <div className="health-item">
            <span className="health-label">Traffic Filter:</span>
            <span className="health-status health-ok">Active</span>
          </div>
          <div className="health-item">
            <span className="health-label">Domain Fronting:</span>
            <span className="health-status health-ok">Enabled</span>
          </div>
          <div className="health-item">
            <span className="health-label">Evasion Engine:</span>
            <span className="health-status health-ok">Active</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
