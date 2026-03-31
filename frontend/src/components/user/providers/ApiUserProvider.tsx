
import { useEffect, useState, type ReactNode } from "react";
import {
    AppointmentKind,
    UserContext,
    UserKind,
    type AppointmentView,
    type AuthUser,
    type User,
    type Uuid
} from "../userContext";


export function UserProvider(prop: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [jwt, setJwt] = useState<string | null>(null);

    useEffect(() => {
        const storedJwt = sessionStorage.getItem("jwt");
        const storedUser = sessionStorage.getItem("user");

        if (storedJwt && storedUser) {
            try {
                setJwt(storedJwt);
                setUser(JSON.parse(storedUser));
            } catch (e) {
                sessionStorage.clear();
            }
        }
    }, []);

    async function connect(email: string, password: string): Promise<boolean> {

        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");

        const res = await fetch("/api/v1/users/login",
            {
                method: "POST",
                body: JSON.stringify({
                    "email": email,
                    "password": password,
                }),
                headers: reqHeaders,
            }
        );

        if (res.status >= 300) return false;

        const data = await res.json();
        const headers = res.headers;
        const auth = headers.get("Authorization");

        if (auth) sessionStorage.setItem("jwt", auth);
        setJwt(auth);

        const token = "Bearer " + data.token;
        setJwt(token);
        sessionStorage.setItem("jwt", token);

        // Fetch full profile since login only returns the token
        const profileHeaders = new Headers();
        profileHeaders.append("Authorization", token);
        const profileRes = await fetch(
            location.origin + "/api/v1/profile/",
            { method: "GET", headers: profileHeaders }
        );

        if (profileRes.status >= 300) return false;

        const profile = await profileRes.json();
        const authUser = {
            email    : profile.email as string,
            kind     : profile.kind as UserKind,
            uuid     : profile.id as string,
            fname    : profile.fname as string,
            lname    : profile.lname as string,
            phone    : profile.phone as string,
            birthday : new Date(profile.birthday),
            address  : profile.address as string,
            appointments: profile.appointments ?? [],
        } as AuthUser;

        setUser(authUser);
        sessionStorage.setItem("user", JSON.stringify(authUser));

        return true;

    }

    async function disconnect() {
        setUser(null);
        sessionStorage.clean();
        setJwt(null);
    }

    async function changeContact(email: string, phone? : string): Promise<boolean> {
        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");
        reqHeaders.append("Authorization", jwt as string);

        const res = await fetch("/api/v1/users/changecontact", {
                method: "POST",
                body: JSON.stringify({
                    uuid: user!.uuid,
                    email: email,
                    phone: phone,
                }),
                headers: reqHeaders,
            }
        );
        if (res.status >= 300) return false;

        setUser({...user, phone: phone} as AuthUser);

        return true;
    }

    async function changeInfo(
        fname: string, lname: string, birthday: Date
    ): Promise<boolean> {
        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");
        reqHeaders.append("Authorization", jwt as string);

        const res = await fetch("/api/v1/users/changeinfo", {
                method: "POST",
                body: JSON.stringify({
                    uuid: user!.uuid,
                    fname:fname,
                    lname:lname,
                    birthday: birthday.toString().split("T")[0],
                }),
                headers: reqHeaders,
            }
        );
        if (res.status >= 300) return false;

        setUser({
            ...user,
            fname: fname,
            lname: lname,
            birthday: birthday
        } as AuthUser);

        return true;
    }

    async function changePassword(password: string): Promise<boolean> {
        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");
        reqHeaders.append("Authorization", jwt as string);

        const res = await fetch(
            location.origin + "/api/v1/users/changepassword", {
                method: "POST",
                body: JSON.stringify({
                    uuid: user!.uuid,
                    password: password
                }),
                headers: reqHeaders,
            }
        );
        if (res.status >= 300) return false;

        return true;
    }

    async function createAccount (
        fname: string, lname: string,
        email: string, password: string,
        birthday: Date, phone: string, address: string | null,
    ): Promise<boolean> {
        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");

        const res = await fetch("/api/v1/users/signup", {
                method: "POST",
                body: JSON.stringify({
                    fname:fname,
                    lname:lname,
                    email: email,
                    phone: phone,
                    password: password,
                    address: address,
                    birthday: birthday.toString().split("T")[0],
                }),
                headers: reqHeaders,
            }
        );
        if (res.status >= 300) return false;

        const headers = res.headers;
        const auth = headers.get("Authorization");

        if (auth) sessionStorage.setItem("jwt", auth);
        setJwt(auth);

        const resp = await res.json() as AuthUser;

        setUser(resp as AuthUser);
        sessionStorage.setItem("user", JSON.stringify(resp));
        // Auto-login after signup
        const loginHeaders = new Headers();
        loginHeaders.append("Content-Type", "application/json");
        const loginRes = await fetch(
            location.origin + "/api/v1/users/login",
            {
                method: "POST",
                body: JSON.stringify({ "email": email, "password": password }),
                headers: loginHeaders,
            }
        );

        if (loginRes.status < 300) {
            const loginData = await loginRes.json();
            const token = "Bearer " + loginData.token;
            setJwt(token);
            sessionStorage.setItem("jwt", token);

            // Fetch profile
            const profileHeaders = new Headers();
            profileHeaders.append("Authorization", token);
            const profileRes = await fetch(
                location.origin + "/api/v1/profile/",
                { method: "GET", headers: profileHeaders }
            );

            if (profileRes.status < 300) {
                const profile = await profileRes.json();
                const authUser = {
                    email    : profile.email as string,
                    kind     : profile.kind as UserKind,
                    uuid     : profile.id as string,
                    fname    : profile.fname as string,
                    lname    : profile.lname as string,
                    phone    : profile.phone as string,
                    birthday : new Date(profile.birthday),
                    address  : profile.address as string,
                    appointments: profile.appointments ?? [],
                } as AuthUser;
                setUser(authUser);
                sessionStorage.setItem("user", JSON.stringify(authUser));
            }
        }

        return true;
    }

    async function createAppointment(
uuid: Uuid, kind: AppointmentKind, datetime: Date
    ): Promise<boolean> {
        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");
        reqHeaders.append("Authorization", jwt as string);

        const res = await fetch(
            "/api/v1/profile/appointments/", {
                method: "POST",
                body: JSON.stringify({
                    uuid: uuid,
                    kind: kind,
                    datetime: datetime,
                }),
                headers: reqHeaders,
            }
        );
        if (res.status >= 300) return false;

        return true;
    }

    async function deleteAppointment(uuid: Uuid, id: number): Promise<boolean> {
        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");
        reqHeaders.append("Authorization", jwt as string);

        const res = await fetch("/api/v1/users/changepassword", {
                method: "DELETE",
                body: JSON.stringify({
                    uuid: uuid,
                    id: id,
                }),
                headers: reqHeaders,
            }
        );
        if (res.status >= 300) return false;

        return true;
    }

    async function viewAppointment(uuid: Uuid, datetime: Date): Promise<AppointmentView[]> {
        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");
        reqHeaders.append("Authorization", jwt as string);

        const res = await fetch("/api/v1/users/viewappointment", {
                method: "POST",
                body: JSON.stringify({
                    uuid: uuid,
                    datetime: datetime,
                }),
                headers: reqHeaders,
            }
        );
        if (res.status >= 300) return [];

        const resp = res.json() as Promise<AppointmentView[]>;

        return resp;
    }

    return(
        <UserContext.Provider value={
            {
                connect: connect,
                disconnect: disconnect,
                changeContact: changeContact,
                changeInfo: changeInfo,
                changePassword: changePassword,
                createAccount: createAccount,
                createAppointment: createAppointment,
                deleteAppointment: deleteAppointment,
                viewAppointment: viewAppointment,
                user: user,
            }
        }>
            { prop.children }
        </UserContext.Provider>
    )
}
