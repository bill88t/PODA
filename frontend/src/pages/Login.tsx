import { useState } from "react";
import { usePath } from "../components/path/pathContext";
import { useUser } from "../components/user/userContext"

export default function Login() {
    const { setPath } = usePath();
    const userCtx = useUser();
    const [email, setEmail] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [error, setError] = useState<string | null>(null);
    return (
        <>
            <title>Saloon PODA - Login</title>
            <div className="form">
                <h1>Login</h1>
                {error ? <div className="span2">{error}</div> : <></>}
                <div>Email</div>
                <div><input
                    type="text"
                    onChange={e => setEmail(e.target.value)}
                /></div>
                <div>Password</div>
                <div><input
                    type="password"
                    min="8"
                    onChange={e => setPassword(e.target.value)}
                /></div>
                <button
                    onClick={
                        async () => { 
                            const ok = await userCtx.connect(email, password)
                        if (ok) {
                                setPath("/");
                            } else setError("Invalid Input"); 
                        }
                    }
                >Login</button>
                <button onClick={
                    () => { setPath("/sign_up"); }
                }>Sign up</button>
            </div>
        </>
    )
}
