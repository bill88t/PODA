import { useEffect } from "react";

function Login() {

    useEffect(() => {
        document.title = "Saloon PADA - Home Page";
    }, []);

    return (
        <>
            <div className="form">
                <h1 className="grid-span-2">Login</h1>
                <div>Username</div>
                <div><input
                    name="username"
                    type="text"
                /></div>
                <div>Password</div>
                <div><input
                    name="password"
                    type="password"
                /></div>
                <button>Login</button>
                <button>Sign up</button>
            </div>
        </>
    )
}

export default Login;
