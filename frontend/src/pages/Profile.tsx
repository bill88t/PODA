import { useState } from "react"

enum flagSwitch {
    password,
    personalInformation,
    contactInformation
};

function ChangePassword(prop: { flag: number; setFlag: (arg0: number) => void; }) {
    return (
        prop.flag & (1 << flagSwitch.password)
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
                        prop.setFlag(prop.flag ^ (1 << flagSwitch.password));
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
                    () => {
                        prop.setFlag(prop.flag ^ (1 << flagSwitch.password));
                    }
                }
            >Change Password</button>
        </>
    );
}

function ChangePersonalData(prop: { flag: number; setFlag: (arg0: number) => void; }) {
    return (
        prop.flag & (1 << flagSwitch.personalInformation)
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
                        prop.setFlag(prop.flag ^ (1 << flagSwitch.personalInformation));
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
                        prop.setFlag(prop.flag ^ (1 << flagSwitch.personalInformation));
                    }
                }
            >Edit Personal Information</button>
        </>

    )
}

function ChangeContactInformation(prop: { flag: number; setFlag: (arg0: number) => void; }) {
    return (
        prop.flag & (1 << flagSwitch.contactInformation)
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
                        prop.setFlag(prop.flag ^ (1 << flagSwitch.contactInformation));
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
                        prop.setFlag(prop.flag ^ (1 << flagSwitch.contactInformation));
                    }
                }
            >Edit Contact Information</button>
        </>
    )
}

export default function Profile() {
    const [
        flag,
        setFlag
    ] = useState<number>(0);
    return (
        <>
        <title>Saloon PODA - Profile</title>
            <div className="form">
                <h1>Profile</h1>
                <h2>Personal Information</h2>
                <ChangePersonalData flag={flag} setFlag={setFlag} />
                <h2>Contact Information</h2>
                <ChangeContactInformation flag={flag} setFlag={setFlag} />
                <h2>Security Information</h2>
                <ChangePassword flag={flag} setFlag={setFlag} />
            </div>
        </>
    );
}
