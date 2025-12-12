import { usePath } from "./pathContext"

export default function NavBar() {
    const { setPath } = usePath();
    return (
        <div className="nav">
            <div className="nav-left">
                <a onClick={ () => setPath("/") }>Home Page</a>
                <a onClick={ () => setPath("/login") }>Login</a>
                <a onClick={ () => setPath("/sign_up") }>Sign up</a>
            </div>
            <div className="nav-right">
                <a onClick={ () => setPath("/profile") }>Profile</a>
            </div>
        </div>
    )
}
