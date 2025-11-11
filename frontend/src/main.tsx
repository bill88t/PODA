import { hydrateRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import NotFound from './NotFound.tsx'
import { BrowserRouter, Route, Routes } from 'react-router'

hydrateRoot(document,
    <BrowserRouter>
        <Routes>
            <Route path="/" element={<App />} />
            <Route path="*" element={<NotFound />} />
        </Routes>
    </BrowserRouter>
)
