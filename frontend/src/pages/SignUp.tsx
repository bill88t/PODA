import { useEffect } from "react";

function SignUp() {

    useEffect(() => {
        document.title = "Saloon PADA - Sign Up";
    }, []);

    return (
        <>
            <div className="form">
                <h1>Sign Up</h1>
                <div>First Name</div>
                <input type="text" name="fname"/>
                <div>Last Name</div>
                <input type="text" name="lname"/>
                <div>Email</div>
                <input
                    name="email"
                    pattern="[a-zA-Z0-9\.]+@[a-zA-Z0-9\.]+(\.[a-z]{2,3})?"
                    type="text"
                />
                <div>Password</div>
                <input
                    name="password"
                    type="password"
                    min="8"
                />
                <div>Phone</div>
                <input
                    name="phone"
                    type="number"
                />
                <div>Address</div>
                <input
                    name="address"
                    type="text"
                />
                <div>Birthday</div>
                <input
                    type="date"
                    name="bdate"
                />
                <button>Sign Up</button>
                <button>Abort</button>
            </div>
        </>
    )
}

export default SignUp;
