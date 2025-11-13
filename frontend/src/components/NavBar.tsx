import { Link } from "react-router"

function NavBar() {
    return (
        <>
            <h1>Nav Bar</h1>
            <Link to="/">Home Page</Link>
            <Link to="/login">Login</Link>
            <Link to="/sign_up">Sign up</Link>
        </>
    )
}

export default NavBar

