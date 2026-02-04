import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Route, Routes, Navigate } from "react-router-dom";
import App from "./pages/App";
import "./styles.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/store/:slug" element={<App />} />
        <Route path="*" element={<Navigate to="/store/demo" replace />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);
