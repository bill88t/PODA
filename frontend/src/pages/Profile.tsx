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
    const [ error, setError ] = useState<string | null>(null);
    return (
        prop.flag & (1 << flagSwitch.password)
        ?
        <>
            { error !== null
                ? <div className="span2 error">{error}</div>
                : <></> }
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
                onClick={
                    async () => {
                        const rgPassword = /\w{8,}/
                        if (!valid) {
                            setError("Password don't match");
                        }
                        if (!rgPassword.test(password)) {
                            setError("Password must be at least 8 characters");
                        }
                        const ok = await changePassword(password);
                        if (ok) prop.setFlag(prop.flag ^ (1 << flagSwitch.password));
                        else setError("Could not change Password");
                    }
                }
            >Confirm Password Change</button>
            <button
                onClick={ () => {
                    prop.setFlag(prop.flag ^ (1 << flagSwitch.password))
                }}>Abort Changes</button>

        </>
        :
        <>
            <button
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
    const [ error, setError ] = useState<string | null>(null);
    const { user, changeInfo } = useUser();
    if (user === null) throw Error("ChangePersonalData user is null");
    const [ fname, setFname ] = useState<string>(user.fname);
    const [ lname, setLname ] = useState<string>(user.lname);
    const [ birthday, setBirthdate ] = useState<Date>(user.birthday);
    return (
        prop.flag & (1 << flagSwitch.personalInformation)
        ?
        <>
            { error !== null
                ? <div className="span2 error">{error}</div>
                : <></> }
            <div>First Name</div>
            <input type="text"
            value={fname}
            onChange={ e => { setFname(e.target.value) }}/>
            <div>Last Name</div>
            <input type="text" value={lname}
            onChange={ e => { setLname(e.target.value) }}/>
            <div>Birthday</div>
            <input type="date"
            value={birthday.toISOString().split("T")[0]}
            onChange={ e => { setBirthdate(new Date(e.target.value)) } }
            />
            <button
                onClick={
                    async () => {
                        if (fname === "") {
                            setError("First name cannot be empty");
                            return;
                        }
                        if (lname === "") {
                            setError("Last name cannot be empty");
                            return;
                        }
                        const ok = await changeInfo(fname, lname, birthday)
                        if (!ok) {
                            setError("Changes are not complete");
                            return;
                        }
                        prop.setFlag(prop.flag ^ (1 << flagSwitch.personalInformation));
                    }
                }
            >Confirm Information Changes</button>
            <button
                onClick={ () => {
                    prop.setFlag(prop.flag ^ (1 << flagSwitch.personalInformation))
                }}>Abort Changes</button>
        </>
        :
        <>
            <div>First Name</div>
            <div className="field">{fname}</div>
            <div>Last Name</div>
            <div className="field">{lname}</div>
            <div>Birthday</div>
            <div className="field">{birthday.toISOString().split("T")[0]}</div>
            <button
                className="span2"
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
    const { user, changeContact } = useUser();
    const [ error, setError ] = useState<string | null>(null);
    if (user === null) throw Error("ChangeContactInformation user is null");
    const [ email, setEmail ] = useState<string>(user.email);
    const [ phone, setPhone ] = useState<string>(
        user.phone !== null ? user.phone : ""
    );
    return (
        prop.flag & (1 << flagSwitch.contactInformation)
        ?
        <>
            { error !== null
                ? <div className="span2 error">{error}</div>
                : <></> }
            <div>Email</div>
            <input type="text" className="Email"
            value={email}
            onChange={ e => {
                    setEmail(e.target.value);
            } }
            />
            <div>Phone Number</div>
            <input type="number"
            value={phone}
            onChange={ e => {
                setPhone(e.target.value);
            } }
            />
            <button
                onClick={
                    async () => {
                        const rgEmail = /\w+@\w+\.\w{2,3}/
                        if (!rgEmail.test(email)) {
                            setError("Give valid email");
                            return;
                        }
                        let ok: boolean = false;
                        if (phone !== "" || phone !== user.phone) {
                            ok = await changeContact(email, phone);
                        } else {
                            ok = await changeContact(email);
                        }
                        if (!ok) {
                            setError("Could not change Email");
                            return;
                        }
                        prop.setFlag(prop.flag ^ (1 << flagSwitch.contactInformation));
                } }
            >Confirm Contact changes</button>
            <button
                onClick={ () => {
                    prop.setFlag(prop.flag ^ (1 << flagSwitch.contactInformation))
                }}>Abort Changes</button>
        </>
        :
        <>
            <div>Email</div>
            <div className="field">{user.email}</div>
            <div>Phone Number</div>
            <div className="field">{user.phone}</div>
            <button
                className="span2"
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
