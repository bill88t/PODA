import { useState } from "react"

function ChangePassword() {

    const [
        changingPassword,
        setChangingPassword
    ] = useState<boolean>(false);
    return (
        changingPassword
        ?
        <>
            <div>Password</div>
            <input className="password1" type="password" />
            <div>Re-enter password</div>
            <input className="password2" type="password" />
            <button
                id="change-password"
                className="span2"
                onClick={
                    () => {
                        setChangingPassword(false);
                    }
                }
            >Confirm Password Change</button>
        </>
        :
        <>
            <button
                id="change-password"
                className="span2"
                onClick={
                    () => setChangingPassword(true)
                }
            >Change Password</button>
        </>
    );
}

function ChangePersonalData() {
    const [
        changingPersonalData,
        setChangingPersonalData
    ] = useState<boolean>(false);
    return (
        changingPersonalData
        ?
        <>
            <div>First Name</div>
            <input type="text" id="fname"/>
            <div>Last Name</div>
            <input type="text" id="lname"/>
            <button
                className="span2"
                id="change-personal-data"
                onClick={
                    () => {
                        setChangingPersonalData(false);
                    }
                }
            >Confirm Information Changes</button>
        </>
        :
        <>
            <div>First Name</div>
            <div id="fname">----------</div>
            <div>Last Name</div>
            <div id="lname">----------</div>
            <button
                className="span2"
                id="change-personal-data"
                onClick={
                    () => {
                        setChangingPersonalData(true);
                    }
                }
            >Edit Personal Information</button>
        </>

    )
}

function ChangeContactInformation() {
    const [
        changingContactInformation,
        setChangingContactInformation
    ] = useState<boolean>(false);
    return (
        changingContactInformation
        ?
        <>
            <div>Email</div>
            <input type="text" className="Email" disabled />
            <div>Phone Number</div>
            <input type="number" className="phone" disabled />
            <button
                className="span2"
                id="change-personal-data"
                onClick={
                    () => {
                        setChangingContactInformation(false);
                    }
                }
            >Confirm Contact changes</button>
        </>
        :
        <>
            <div>Email</div>
            <div id="email">-------@-----.--</div>
            <div>Phone Number</div>
            <div id="phone">-------------</div>
            <button
                className="span2"
                id="change-personal-data"
                onClick={
                    () => {
                        setChangingContactInformation(true);
                    }
                }
            >Edit Contact Information</button>
        </>
    )
}

export default function Profile() {
    return (
        <div className="form">
            <h1>Profile</h1>
            <h2>Personal Information</h2>
            <ChangePersonalData />
            <h2>Contact Information</h2>
            <ChangeContactInformation />
            <h2>Security Information</h2>
            <ChangePassword />
        </div>
    );
}
