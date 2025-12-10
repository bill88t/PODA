import { createRoot } from 'react-dom/client'
import NavBar from './components/NavBar.tsx'
import "./index.sass"
import React, { lazy, useState } from 'react'

const Home = lazy(() => import("./pages/Home.tsx"))
const SignUp = lazy(() => import("./pages/SignUp.tsx"))
const Login = lazy(() => import("./pages/Login.tsx"))
const Profile = lazy(() => import("./pages/Profile.tsx"))
const NotFound = lazy(() => import("./pages/NotFound.tsx"))

export function Root() {
    const [path, setPath] = useState<string>(window.location.pathname);
    return (
        <React.StrictMode>
            <NavBar path={path} setPath={setPath} />
            {
                path === "/" ? <Home /> :
                path === "/login" ? <Login /> :
                path === "/sign_up" ? <SignUp /> :
                path === "/profile" ? <Profile /> :
                <NotFound path={path} setPath={setPath}/>
            }
        </React.StrictMode>
    )
}

createRoot(document.getElementById("root")!).render(<Root />
)
