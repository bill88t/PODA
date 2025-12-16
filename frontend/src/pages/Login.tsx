import { useState } from "react";
import { usePath } from "../components/path/pathContext";
import { UserKind, useUser } from "../components/user/userContext"

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
                    pattern="[a-zA-Z0-9\.]+@[a-zA-Z0-9\.]+(\.[a-z]{2,3})?"
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
                        () => { 
                            const u = userCtx.connect(email, password)
                            if (u.kind !== UserKind.guess) { 
                                setPath("/")
                            } else setError("Invalid Input"); 
                        }
                    }
                >Login</button>
                <button onClick={
                    () => setPath("/sign_up")
                }>Sign up</button>
            </div>
        </>
    )
}
