import { useEffect } from "react";

function Login() {

    useEffect(() => {
        document.title = "Saloon PADA - Home Page";
    }, []);

    return (
        <>
            <div className="form">
                <div className="form-label">Username</div>
                <div className="form-input"><input
                    name="username"
                    type="text"
                /></div>
                <div className="form-label">Password</div>
                <div className="form-input"><input
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
