import { useEffect } from "react";
import { Link, useNavigate } from "react-router"

function NotFound() {

    useEffect(() => {
        document.title = "Saloon PADA - Page Not found";
    }, []);

    const navigate = useNavigate();

    useEffect(() => {

    const timeout = setTimeout(navigate, 5000, "/", { replace: true });

    return () => clearTimeout(timeout);
    }, [navigate]);

    return (
        <>
            <h1>Page not found</h1>
            <p>
                Unfortunally this link is broken.
                We redirect you to home page
                If that doesn't work, press <Link to="/" >here</Link>
            </p>
        </>
    )
}

export default NotFound
