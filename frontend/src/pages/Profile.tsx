import { useEffect, useState } from "react"

function ChangePassword() {

const [changingPassword, setChangingPassword] = useState<boolean>(false);
    useEffect(() => {
        return () => {
        }
    }, [changingPassword])
    return (
        changingPassword
            ?
        <>
            <button
                id="change-password"
                className="span2"
                onClick={
                    () => setChangingPassword(false)
                }
            >Change Password</button>
        </>
        :
        <>
            <div>Password</div>
            <input className="password1" type="password" />
            <div>Re-enter password</div>
            <input className="password2" type="password" />
            <button
                id="change-password"
                className="span2"
                onClick={
                    () => setChangingPassword(true)
                }
            >Confirm Change</button>
        </>
    );
}

export default function Profile() {
    return (
        <div className="form">
            <h1>Profile</h1>

            <h2>Personal Information</h2>
            <div>First Name</div>
            <input type="text" className="fname" disabled />
            <div>Last Name</div>
            <input type="text" className="lname" disabled />
            <a className="span2">Edit Personal Information</a>

            <h2>Contact Information</h2>
            <div>Email</div>
            <input type="text" className="Email" disabled />
            <a className="span2">Edit Personal Information</a>

            <h2>Security Information</h2>
            <ChangePassword />
        </div>
    );
}
