import { useEffect } from "react";

function Login() {

    useEffect(() => {
        document.title = "Saloon PADA - Home Page";
    }, []);

    return (
        <>
            <div className="form">
                <div className="form-label">Username</div>
                <div className="form-input"><input name="Username" /></div>
            </div>
        </>
    )
}

export default Login;
