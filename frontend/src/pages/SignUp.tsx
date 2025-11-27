import { useEffect } from "react";

function SignUp() {

    useEffect(() => {
        document.title = "Saloon PADA - Sign Up";
    }, []);

    return (
        <>
            <div className="form">
                <h1>Login</h1>
                <div>First Name</div>
                <input type="text" name="fname"/>
                <div>Last Name</div>
                <input type="text" name="lname"/>
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
                <div>Phone</div>
                <div><input
                    name="phone"
                    type="number"
                /></div>
                <div>Address</div>
                <div><input
                    name="address"
                    type="text"
                /></div>
                <button>Sign Up</button>
                <button>Abort</button>
            </div>
        </>
    )
}

export default SignUp;
