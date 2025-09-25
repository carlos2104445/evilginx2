import React, { useState, useEffect } from 'react';
import Dashboard from './components/Dashboard';
import PhishletManager from './components/PhishletManager';
import SessionMonitor from './components/SessionMonitor';
import ConfigManager from './components/ConfigManager';
import './App.css';

interface NavItem {
  id: string;
  label: string;
  component: React.ComponentType;
}

const App: React.FC = () => {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [isConnected, setIsConnected] = useState(false);

  const navItems: NavItem[] = [
    { id: 'dashboard', label: 'Dashboard', component: Dashboard },
    { id: 'phishlets', label: 'Phishlets', component: PhishletManager },
    { id: 'sessions', label: 'Sessions', component: SessionMonitor },
    { id: 'config', label: 'Configuration', component: ConfigManager },
  ];

  useEffect(() => {
    const checkConnection = async () => {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/health`);
        setIsConnected(response.ok);
      } catch (error) {
        setIsConnected(false);
      }
    };

    checkConnection();
    const interval = setInterval(checkConnection, 30000);
    return () => clearInterval(interval);
  }, []);

  const ActiveComponent = navItems.find(item => item.id === activeTab)?.component || Dashboard;

  return (
    <div className="app">
      <header className="app-header">
        <div className="header-content">
          <h1 className="app-title">Evilginx2 Management Console</h1>
          <div className="connection-status">
            <div className={`status-indicator ${isConnected ? 'connected' : 'disconnected'}`}></div>
            <span>{isConnected ? 'Connected' : 'Disconnected'}</span>
          </div>
        </div>
      </header>

      <nav className="app-nav">
        <div className="nav-content">
          {navItems.map(item => (
            <button
              key={item.id}
              className={`nav-item ${activeTab === item.id ? 'active' : ''}`}
              onClick={() => setActiveTab(item.id)}
            >
              {item.label}
            </button>
          ))}
        </div>
      </nav>

      <main className="app-main">
        <div className="main-content">
          <ActiveComponent />
        </div>
      </main>
    </div>
  );
};

export default App;
