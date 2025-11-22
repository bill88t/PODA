import { useEffect } from "react";
import { NavLink, useNavigate } from "react-router"

function NotFound() {

    useEffect(() => {
        document.title = "Saloon PADA - Page Not found";
    }, []);

    const navigate = useNavigate();
    setTimeout(navigate, 5000, "/", { replace: true });

    return (
        <>
            <h1>Page not found</h1>
            <p>
                Unfortunally this link is broken.
                We redirect you to home page
                If that doesn't work, press <NavLink to="/">here</NavLink>
            </p>
        </>
    )
}

export default NotFound
