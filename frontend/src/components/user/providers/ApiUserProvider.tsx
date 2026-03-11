
import { useState, type ReactNode } from "react";
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

    async function connect(email: string, password: string): Promise<boolean> {
        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");

        const res = await fetch(
            location.origin + "/api/v1/users/login",
            {
                method: "POST",
                body: JSON.stringify({ "email": email, "password": password }),
                headers: reqHeaders,
            }
        );

        if (res.status >= 300) return false;

        const data = await res.json();
        const token = "Bearer " + data.token;
        setJwt(token);

        // Fetch full profile since login only returns the token
        const profileHeaders = new Headers();
        profileHeaders.append("Authorization", token);
        const profileRes = await fetch(
            location.origin + "/api/v1/profile/",
            { method: "GET", headers: profileHeaders }
        );

        if (profileRes.status >= 300) return false;

        const profile = await profileRes.json();
        setUser({
            email    : profile.email as string,
            kind     : profile.kind as UserKind,
            uuid     : profile.id as string,
            fname    : profile.fname as string,
            lname    : profile.lname as string,
            phone    : profile.phone as string,
            birthday : new Date(profile.birthday),
            address  : profile.address as string,
            appointments: profile.appointments ?? [],
        } as AuthUser);

        return true;
    }

    async function disconnect() {
        setUser(null);
    }

    async function changeContact(email: string, phone? : string): Promise<boolean> {
        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");
        reqHeaders.append("Authorization", jwt as string);

        const res = await fetch(
            location.origin + "/api/v1/users/changecontact", {
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

        const res = await fetch(
            location.origin + "/api/v1/users/changeinfo", {
                method: "POST",
                body: JSON.stringify({
                    uuid: user!.uuid,
                    fname:fname,
                    lname:lname,
                    birthday: birthday.toISOString().split("T")[0],
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

    async function createAccount(
        fname: string, lname: string,
        email: string, password: string,
        birthday: Date, phone: string, address: string | null,
    ): Promise<boolean> {
        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");

        const res = await fetch(
            location.origin + "/api/v1/users/signup", {
                method: "POST",
                body: JSON.stringify({
                    fname, lname, email, phone, password, address,
                    birthday: birthday.toISOString().split("T")[0],
                }),
                headers: reqHeaders,
            }
        );
        return res.status < 300;
    }

    async function createAppointment(
        _uuid: Uuid, kind: AppointmentKind, datetime: Date
    ): Promise<boolean> {
        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");
        reqHeaders.append("Authorization", jwt as string);

        const res = await fetch(
            location.origin + "/api/v1/profile/appointments/", {
                method: "POST",
                body: JSON.stringify({
                    kind: kind,
                    datetime: datetime.toISOString(),
                }),
                headers: reqHeaders,
            }
        );
        return res.status < 300;
    }

    async function deleteAppointment(_uuid: Uuid, id: number): Promise<boolean> {
        const reqHeaders = new Headers();
        reqHeaders.append("Authorization", jwt as string);

        const res = await fetch(
            location.origin + `/api/v1/profile/appointments/${id}`, {
                method: "DELETE",
                headers: reqHeaders,
            }
        );
        return res.status < 300;
    }

    async function viewAppointment(_uuid: Uuid, _datetime: Date): Promise<AppointmentView[]> {
        const reqHeaders = new Headers();
        reqHeaders.append("Authorization", jwt as string);

        const res = await fetch(
            location.origin + "/api/v1/profile/appointments/", {
                method: "GET",
                headers: reqHeaders,
            }
        );
        if (res.status >= 300) return [];

        const appointments = await res.json() as {
            id: number;
            user_id?: string;
            datetime: string;
            kind: string;
        }[];

        return appointments.map(a => ({
            userUuid: a.user_id ?? user!.uuid,
            fname: user!.fname,
            lname: user!.lname,
            appointment: {
                id: a.id,
                datetime: new Date(a.datetime),
                kind: a.kind as AppointmentKind,
            },
        }));
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
