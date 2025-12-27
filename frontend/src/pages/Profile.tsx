import { useState } from "react"
import { useUser } from "../components/user/userContext";

enum flagSwitch {
    password,
    personalInformation,
    contactInformation
};

function ChangePassword(prop: { flag: number; setFlag: (arg0: number) => void; }) {
    const { changePassword }  = useUser();
    const [ password, setPassword ] = useState<string>("");
    const [ valid, setValid ] = useState<boolean>(false);
    const { error, setError } = useState<string | null>(null);
    return (
        prop.flag & (1 << flagSwitch.password)
        ?
        <>
            { error !== null ? <div className="span2 error">{error}</div> : <></> }
            <div>Password</div>
            <input type="password"
                onChange={ e => {
                    setPassword(e.target.value);
                } }
            />
            <div>Re-enter password</div>
            <input type="password"
                onChange={ e => {
                    setValid(e.target.value === password)
                } }
                />
            <button
                className="span2"
                onClick={
                    async () => {
                        const rgPassword = /\w{8,}/
                        if (!valid) {
                            setError("Password don't match");
                        }
                        if (!rgPassword.test(password)) {
                            setError("Password must ve at least 8 characters");
                        }
                        const ok = await changePassword(password);
                        if (ok) prop.setFlag(prop.flag ^ (1 << flagSwitch.password));
                        else setError("Could not change Password");
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
