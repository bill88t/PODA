import { createRoot } from 'react-dom/client'
import React, { lazy } from 'react'
import "./index.sass"
import PathProvider from './components/path/PathProvider.tsx'
import NavBar from './components/NavBar.tsx'
import { usePath } from './components/path/pathContext.tsx'
import { UserProvider } from './components/user/UserProvider.tsx'

const Home = lazy(() => import("./pages/Home.tsx"))
const SignUp = lazy(() => import("./pages/SignUp.tsx"))
const Login = lazy(() => import("./pages/Login.tsx"))
const Profile = lazy(() => import("./pages/Profile.tsx"))
const NotFound = lazy(() => import("./pages/NotFound.tsx"))

function Main() {
    const { path } = usePath();
    switch (path) {
        case "/"       : return <Home />;
        case "/login"  : return <Login />;
        case "/sign_up": return <SignUp />;
        case "/profile": return <Profile />;
        default        : return <NotFound />;
    }
}

export function Root() {
    return (
        <React.StrictMode>
            <UserProvider>
                <PathProvider>
                    <NavBar />
                    <Main />
                </PathProvider>
            </UserProvider>
        </React.StrictMode>
    )
}

createRoot(document.getElementById("root")!).render(<Root />
)
