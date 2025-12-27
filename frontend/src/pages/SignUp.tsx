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
    const [validPassword, setValidPassword] = useState<boolean>(false);
    const [phone,    setPhone]    = useState<string>("");
    const [address,  setAddress]  = useState<string>("");
    const [bdate,    setBdate]    = useState<Date>(new Date(Date.now()));
    return (
        <>
            <title>Saloon PODA - Sign Up Page</title>
            <div className="form">
                <h1>Sign Up</h1>
                {
                    error ? <div className="span2 error">{error}</div> : <></>
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
                <div>Password again</div>
                <input type="password"
                    onChange={e => setValidPassword(e.target.value === password)}    
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
                            if (fname === "") {
                                setError("Missing First Name");
                                return;
                            }
                            if (lname === "") {
                                setError("Missing Last Name");
                                return;
                            }
                            if (email === "") {
                                setError("Missing Email");
                                return;
                            }
                            const rgEmail = /\w+@\w+\.\w{2,3}?/;
                            if (!rgEmail.test(email)) {
                                setError("Give valid Email");
                                return;
                            }
                            if (bdate === null) {
                                setError("Missing birthday");
                                return;
                            }
                            if (phone === "") {
                                setError("Missing Phone");
                                return;
                            }
                            const rgPhone = /(\+[0-9]{2})?[0-9]{10}/
                            if (!rgPhone.test(phone)) {
                                setError("Give valid phone");
                                return;
                            }
                            if (!validPassword) {
                                setError("Passwords don't match");
                                return;
                            }
                            const rgPassword = /\w{8,}/
                            if (!rgPassword.test(password)) {
                                setError("Password must contain 8 characters");
                                return;
                            }
                            const addr = address !== "" ? address : null;
                            const ok = await userCtx.createAccount(
                                fname, lname,
                                email, password,
                                bdate, phone,
                                addr, UserKind.client
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
