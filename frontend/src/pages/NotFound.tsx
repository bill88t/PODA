import { useEffect } from "react";
import { usePath } from "../components/pathContext";

export default function NotFound() {

    const { path, setPath } = usePath();

    useEffect(() => {

    const timeout = setTimeout(() => setPath("/"), 5000, "/", { replace: true });

    return () => clearTimeout(timeout);
    }, [path, setPath]);

    return (
        <>
            <title>Saloon PODA - Page Not Found</title>
            <h1>Page not found</h1>
            <p>
                Unfortunately this link is broken.
                We redirect you to home page
                If that doesn't work, press <a onClick={ () => setPath("/") } >here</a>
            </p>
        </>
    )
}
