
import { useState, type ReactNode } from "react";
import {
    AppointmentKind,
    UserContext,
    UserKind,
    // useUser,
    type Appointment,
    type AppointmentView,
    type AuthUser,
    type User,
    type Uuid
} from "../userContext";


export function UserProvider(prop: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    // const userCtx = useUser();
    const [jwt, setJwt] = useState<string | null>();

    async function connect(email: string, password: string): Promise<boolean> {

        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");

        const res = await fetch(
             location.origin + "/api/v1/users/login",
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
        setJwt(auth);

        setUser({
            email    : email as string,
            kind     : data.kind as UserKind,
            uuid     : data.id as number,
            fname    : data.fname as string,
            lname    : data.lname as string,
            phone    : data.phone as string,
            birthday : new Date(data.birthday),
            address  : data.address as string,
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
        kind: UserKind
    ): Promise<boolean> {
        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");

        const res = await fetch(
            location.origin + "/api/v1/users/createaccount", {
                method: "POST",
                body: JSON.stringify({
                    uuid: user!.uuid,
                    fname:fname,
                    lname:lname,
                    email: email,
                    password: password,
                    address: address,
                    kind: kind,
                    birthday: birthday.toString().split("T")[0],
                }),
                headers: reqHeaders,
            }
        );
        if (res.status >= 300) return false;

        const resp = await res.json() as AuthUser;

        setUser(resp as AuthUser);

        return true;
    }

    async function createAppointment(
uuid: Uuid, kind: AppointmentKind, datetime: Date
    ): Promise<boolean> {
        for (let i = 0; i < users.length; ++i) {
            if (user && users[i].uuid === uuid) {
                for (let j = 0; j < users[i].appointments.length; ++j) {
                    if (users[i].appointments[j].datetime === datetime) {
                        return false;
                    }
                }

                const ap: Appointment = { datetime: datetime, id: idc++, kind: kind }

                users[i].appointments.push(ap);
                return true;
            }
        }
        return false;
    }

    async function deleteAppointment(uuid: Uuid, id: number): Promise<boolean> {
        for (let i = 0; i < users.length; ++i) {
            if (user && users[i].uuid === uuid) {
                for (let j = 0; j < users[i].appointments.length; ++j) {
                    if (users[i].appointments[j].id === id) {
                            users[i].appointments.splice(j, 1);
                            return true;
                    }
                }
            }
        }
    return false;
    }

    async function viewAppointment(uuid: Uuid, datetime: Date) {
        const app: AppointmentView[] = [];
        if (    user && (user.kind === UserKind.client
            ||  user.kind === UserKind.admin)) {
            for (let i = 0; i < users.length; ++i) {
                for (let j = 0; j < users[i].appointments.length; ++j) {
                    if (users[i].appointments[j].datetime >= datetime) {
                        app.push({
                            fname: users[i].fname,
                            lname: users[i].lname,
                            userUuid: users[i].uuid,
                            appointment: users[i].appointments[j]
                        });
                    }
                }
            }
            return app;
        }
        for (let i = 0; i < users.length; ++i) {
            if (user && users[i].uuid === uuid) {
                for (let j = 0; j < users[i].appointments.length; ++j) {
                    if (users[i].appointments[j].datetime >= datetime) {
                        app.push({
                            fname: users[i].fname,
                            lname: users[i].lname,
                            userUuid: uuid,
                            appointment: users[i].appointments[j]
                        });
                    }
                }
                break;
            }
        }
        return app;
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
