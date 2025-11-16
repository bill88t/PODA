import { useEffect } from "react"

function App() {

    useEffect(() => {
        document.title = "Saloon PADA - Home Page";
    }, []);

    return (
        <>
            <h1>Welcome to Our Saloon</h1>
            <p>
                The PADA Saloon is open from 2025
                our history is new, but we would like to create
                stories for our community.
            </p>
        </>
    )
}

export default App
