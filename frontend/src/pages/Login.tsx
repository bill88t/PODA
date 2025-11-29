import { useEffect } from "react";

function Login() {

    useEffect(() => {
        document.title = "Saloon PADA - Login";
    }, []);

    return (
        <>
            <div className="form">
                <h1>Login</h1>
                <div>Email</div>
                <div><input
                    name="email"
                    pattern="[a-zA-Z0-9\.]+@[a-zA-Z0-9\.]+(\.[a-z]{2,3})?"
                    type="text"
                /></div>
                <div>Password</div>
                <div><input
                    name="password"
                    type="password"
                    min="8"
                /></div>
                <button>Login</button>
                <button>Sign up</button>
            </div>
        </>
    )
}

export default Login;
