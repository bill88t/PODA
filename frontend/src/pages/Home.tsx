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
                We would like to have you here with us
                and help us leverage your hairstyle,
                because you deserve that.
                You are amazing and we are gonna make it
                appear, so everyone is gonna see it.
            </p>
        </>
    )
}

export default App
