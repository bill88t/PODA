import { useEffect } from "react";

function NotFound(prop: { setPath: (arg0: string) => void; path: unknown; }) {

    useEffect(() => {
        document.title = "Saloon PADA - Page Not found";
    }, []);

    useEffect(() => {

    const timeout = setTimeout(() => prop.setPath("/"), 5000, "/", { replace: true });

    return () => clearTimeout(timeout);
    }, [prop.path]);

    return (
        <>
            <h1>Page not found</h1>
            <p>
                Unfortunally this link is broken.
                We redirect you to home page
                If that doesn't work, press <a onClick={ () => prop.setPath("/") } >here</a>
            </p>
        </>
    )
}

export default NotFound
