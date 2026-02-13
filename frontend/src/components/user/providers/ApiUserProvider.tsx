
import { useState, type ReactNode } from "react";
import {
    AppointmentKind,
    UserContext,
    UserKind,
    useUser,
    type Appointment,
    type AppointmentView,
    type User,
    type Uuid
} from "../userContext";


export function UserProvider(prop: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const userCtx = useUser();
    const [jwt, setJwt] = useState<string | null>();

    async function connect(email: string, password: string): Promise<boolean> {

        const reqHeaders = new Headers();
        reqHeaders.append("Content-Type", "application/json");

        const res = await fetch(
            "https://localhost:5173/api/v1/users/login",
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
        const auth = headers.get("authorization");
        setJwt(auth);

        user!.email    = email;
        user!.kind     = data.kind;
        user!.uuid     = data.id;
        user!.fname    = data.fname;
        user!.lname    = data.lname;
        user!.phone    = data.phone;
        user!.birthday = data.birthday;
        user!.address  = data.address;

        return true;

    }

    async function disconnect() {
        setUser(null);
    }

    async function changeContact(email: string, phone? : string): Promise<boolean> {
        for (let i = 0; i < users.length; ++i) {
            if (user && users[i].email === user.email) {
                users[i].email = email;
                break;
            }
        }
        if (user) {
            let u = {...user, email: email};
            if (phone !== undefined) {
                u = {...user, phone: phone};
            }
            setUser(u);
            return true;
        }
        return false;
    }

    async function changeInfo(
        fname: string, lname: string, birthday: Date
    ): Promise<boolean> {
        for (let i = 0; i < users.length; ++i) {
            if (user && users[i].email === user.email) {
                users[i].fname = fname;
                users[i].lname = lname;
                users[i].birthday = birthday;
            }
        }

        if (user) {
            const u = {...user, fname: fname,
                        lname: lname, birthday: birthday};
            setUser(u);
            return true;
        }
        return false;
    }

    async function changePassword(password: string): Promise<boolean> {
        for (let i = 0; i < users.length; ++i) {
            if (user && users[i].email === user.email) {
                users[i].password = password;
            }
        }

        if (user) {
            const u = {...user, password: password};
            setUser(u);
            return true;
        }
        return false;

    }

    async function createAccount (
        fname: string, lname: string,
        email: string, password: string,
        birthday: Date, phone: string, address: string | null,
        kind: UserKind
    ): Promise<boolean> {
        const u = {fname: fname, lname: lname, email: email,
        password: password, birthday: birthday, phone: phone,
        address: address, kind: kind, uuid: uui++,
        appointments: []}
        if( users.filter(us => us.email == email).length == 0 ) {
            users = [...users, u];
            setUser(u);
            return true;
        }
        return false;
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
