import { Link } from "react-router"

function NavBar() {
    return (
        <div className="nav">
            <div className="nav-left">
                <Link to="/">Home Page</Link>
                <Link to="/login">Login</Link>
                <Link to="/sign_up">Sign up</Link>
            </div>
            <div className="nav-right">
                <Link to="/profile">Profile</Link>
            </div>
        </div>
    )
}

export default NavBar

