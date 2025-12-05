import { createRoot } from 'react-dom/client'
import Home from './pages/Home.tsx'
import NotFound from './pages/NotFound.tsx'
import { BrowserRouter, Route, Routes } from 'react-router'
import NavBar from './components/NavBar.tsx'
import Login from './pages/Login.tsx'
import "./index.sass"
import React from 'react'
import SignUp from './pages/SignUp.tsx'
import Profile from './pages/Profile.tsx'

createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
    <BrowserRouter>
        <NavBar />
        <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route path="/sign_up" element={<SignUp />} />
            <Route path="/profile" element={<Profile />} />
            <Route path="*" element={<NotFound />} />
        </Routes>
    </BrowserRouter>
    </React.StrictMode>
)
