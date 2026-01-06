import { useState, type ReactNode } from "react";
import { UserContext, UserKind, type AuthUser, type User } from "../userContext";

let uui = 1;

let users: AuthUser[] = [
    {
        uuid: uui++,
        fname: "John",
        lname: "Bar",
        email: "john@j.j",
        password: "john",
        phone: "123",
        address: null,
        birthday: new Date("2010-01-02"),
        kind: UserKind.client,
        appointments: []
    }
];

export function UserProvider(prop: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);

    async function connect(email: string, password: string): Promise<boolean> {
        for (let i = 0; i<users.length; ++i) {
            if (users[i].email === email
                    && users[i].password === password) {
                setUser(users[i]);
                return true;
            }
        }
        return false;
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

    return(
        <UserContext.Provider value={
            {
                connect: connect,
                disconnect: disconnect,
                changeContact: changeContact,
                changeInfo: changeInfo,
                changePassword: changePassword,
                createAccount: createAccount,
                user: user,
            }
        }>
            { prop.children }
        </UserContext.Provider>
    )
}
