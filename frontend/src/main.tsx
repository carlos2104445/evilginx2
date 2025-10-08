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

const root = createRoot(document.getElementById("root")!);

root.render(
  <React.StrictMode>
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
        </Route>
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);
