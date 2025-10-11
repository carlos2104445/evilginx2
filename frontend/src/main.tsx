import React from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import App from "./App";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import Phishlets from "./pages/Phishlets";
import Sessions from "./pages/Sessions";
import Lures from "./pages/Lures";
import Config from "./pages/Config";
import Certificates from "./pages/Certificates";
import { ToastProvider } from "./components/Toast";

const root = createRoot(document.getElementById("root")!);

root.render(
  <React.StrictMode>
    <ToastProvider>
      <BrowserRouter>
        <Routes>
          <Route element={<App />}>
            <Route path="/login" element={<Login />} />
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/phishlets" element={<Phishlets />} />
            <Route path="/sessions" element={<Sessions />} />
            <Route path="/lures" element={<Lures />} />
            <Route path="/config" element={<Config />} />
            <Route path="/certificates" element={<Certificates />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ToastProvider>
  </React.StrictMode>
);
