import { Link } from "react-router"

function NavBar() {
    return (
        <>
            <Link to="/">Home Page</Link>
            <Link to="/login">Login</Link>
            <Link to="/sign_up">Sign up</Link>
            <Link to="/profile">Profile</Link>
        </>
    )
}

export default NavBar

