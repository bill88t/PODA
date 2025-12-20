import { useState } from "react";
import { UserKind, useUser } from "../components/user/userContext"
import { usePath } from "../components/path/pathContext";

export default function SignUp() {
    const userCtx = useUser();
    const { setPath } = usePath();

    const [error,    setError]    = useState<string | null>(null);
    const [fname,    setFname]    = useState<string>("");
    const [lname,    setLname]    = useState<string>("");
    const [email,    setEmail]    = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [phone,    setPhone]    = useState<string>("");
    const [address,  setAddress]  = useState<string>("");
    const [bdate,    setBdate]    = useState<Date>(new Date(Date.now()));
    return (
        <>
            <title>Saloon PODA - Sign Up Page</title>
            <div className="form">
                <h1>Sign Up</h1>
                {
                    error ? <div className="span2">{error}</div> : <></>
                }
                <div>First Name</div>
                <input type="text"
                    onChange={e => setFname(e.target.value)}    
                />
                <div>Last Name</div>
                <input type="text"
                    onChange={e => setLname(e.target.value)}    
                />
                <div>Email</div>
                <input
                    type="text"
                    onChange={e => setEmail(e.target.value)}    
                />
                <div>Password</div>
                <input type="password"
                    onChange={e => setPassword(e.target.value)}    
                />
                <div>Phone</div>
                <input type="number"
                    onChange={e => setPhone(e.target.value)}    
                />
                <div>Address</div>
                <input type="text" 
                    onChange={e => setAddress(e.target.value)}    
                />
                <div>Birthday</div>
                <input type="date"
                    onChange={e => setBdate(new Date(e.target.value))}    
                />
                <button
                    onClick={
                        async () => {
                            const ok = await userCtx.createAccount(
                                fname, lname,
                                email, password,
                                bdate, phone,
                                address, UserKind.client
                            );
                            if (ok) setPath("/");
                            else setError("Could not create account")
                        }
                    }
                >Sign Up</button>
                <button>Abort</button>
            </div>
        </>
    )
}
