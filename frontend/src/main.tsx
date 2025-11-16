import { hydrateRoot } from 'react-dom/client'
import Home from './pages/Home.tsx'
import NotFound from './pages/NotFound.tsx'
import { BrowserRouter, Route, Routes } from 'react-router'
import NavBar from './components/NavBar.tsx'

hydrateRoot(document,
    <BrowserRouter>
        <NavBar />
        <Routes>
            <Route path="/" element={<Home />} />
            <Route path="*" element={<NotFound />} />
        </Routes>
    </BrowserRouter>
)
