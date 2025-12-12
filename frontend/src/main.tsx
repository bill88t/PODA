import { createRoot } from 'react-dom/client'
import NavBar from './components/NavBar.tsx'
import "./index.sass"
import React, { lazy } from 'react'
import PathProvider from './components/PathProvider.tsx'
import { usePath } from './components/pathContext.tsx'

const Home = lazy(() => import("./pages/Home.tsx"))
const SignUp = lazy(() => import("./pages/SignUp.tsx"))
const Login = lazy(() => import("./pages/Login.tsx"))
const Profile = lazy(() => import("./pages/Profile.tsx"))
const NotFound = lazy(() => import("./pages/NotFound.tsx"))

function Main() {
    const { path } = usePath();
    return (
        <>
            {
                path === "/" ? <Home /> :
                path === "/login" ? <Login /> :
                path === "/sign_up" ? <SignUp /> :
                path === "/profile" ? <Profile /> :
                <NotFound />
            }
        </>
    )
}

export function Root() {
    return (
        <React.StrictMode>
        <PathProvider>
            <NavBar />
            <Main />
            </PathProvider>
        </React.StrictMode>
    )
}

createRoot(document.getElementById("root")!).render(<Root />
)
