import { hydrateRoot } from 'react-dom/client'
import Home from './pages/Home.tsx'
import NotFound from './pages/NotFound.tsx'
import { BrowserRouter, Route, Routes } from 'react-router'

hydrateRoot(document,
    <BrowserRouter>
        <Routes>
            <Route path="/" element={<Home />} />
            <Route path="*" element={<NotFound />} />
        </Routes>
    </BrowserRouter>
)
