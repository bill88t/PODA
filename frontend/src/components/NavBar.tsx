import { Link } from "react-router"

function NavBar() {
    return (
        <div className="nav">
            <Link to="/">Home Page</Link>
            <Link to="/login">Login</Link>
            <Link to="/sign_up">Sign up</Link>
            <Link to="/profile">Profile</Link>
        </div>
    )
}

export default NavBar

