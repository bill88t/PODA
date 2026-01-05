import { createRoot } from 'react-dom/client'
import { lazy, StrictMode } from 'react'
import "./index.sass"
import PathProvider from './components/path/PathProvider.tsx'
import NavBar from './components/NavBar.tsx'
import { usePath } from './components/path/pathContext.tsx'
import { UserProvider } from './components/user/UserProvider.tsx'
import { useUser } from './components/user/userContext.tsx'

const Home = lazy(() => import("./pages/Home.tsx"))
const SignUp = lazy(() => import("./pages/SignUp.tsx"))
const Login = lazy(() => import("./pages/Login.tsx"))
const Profile = lazy(() => import("./pages/Profile.tsx"))
const NotFound = lazy(() => import("./pages/NotFound.tsx"))

function Main() {
    const { path } = usePath();
    const { user } = useUser();
    switch (path) {
        case "/"       : return <Home />;
        case "/login"  : {
            if (user === null) {
                return <Login />
            }
            else { 
                history.pushState({}, "", "/");
                return <Home />
            };
        }
        case "/sign_up": return <SignUp />;
        case "/profile": {
            if (user !== null) {
                return <Profile />
            } else {
                history.pushState({}, "", "/");
                return <Home />
            }
        }
        default        : return <NotFound />;
    }
}

export function Root() {
    return (
        <StrictMode>
            <PathProvider>
                <UserProvider>
                    <NavBar />
                    <Main />
                </UserProvider>
            </PathProvider>
        </StrictMode>
    )
}

createRoot(document.getElementById("root")!).render(<Root />
)
