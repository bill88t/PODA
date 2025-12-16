import { useState, type ReactNode } from "react";
import { UserContext, UserKind, type User } from "../userContext";

const users: User[] = [
    {
        kind: UserKind.client,
        email: "john@email.com",
        password: "john1",
        phone: "12345",
        birthday: new Date("2020-01-02")
    }
]

export function UserProvider(prop: { children: ReactNode }) {

    const [user, setUser] = useState<User>({kind:UserKind.guess});

    function connect(email: string, password: string) {
        const u = users.filter(
            u => u.email == email && u.password == password);
        if (u.length === 1) {
            setUser(u[0]);
            return u[0];
        }
        return {kind: UserKind.guess}
    }

    function disconnect() {
        setUser({kind:UserKind.guess})
        return;
    }

    function changeEmail(email: string) {
        for (let i = 0; i < users.length; ++i) {
            if (users[i].email == user.email) {
                users[i].email = email
            }
        }
        setUser({...user, email: email});
        return user;
    }

    function changeInfo(fname: string, lname: string) {
        for (let i = 0; i < users.length; ++i) {
            if (users[i].email == user.email) {
                users[i].fname = fname;
                users[i].lname = lname;
            }
        }
        setUser({...user, fname: fname, lname: lname});
        return user;

    }

    function changePassword(password: string) {
        for (let i = 0; i < users.length; ++i) {
            if (users[i].email == user.email) {
                users[i].password = password;
            }
        }
        setUser({...user, password: password});
        return user;

    }

    function getUser() { return user; }

    function createAccount(
        fname: string, lname: string,
        email: string, password: string,
        birthday: Date, phone: string, kind: UserKind
    ) {
        const u = {fname: fname, lname: lname, email: email,
        password: password, birthday: birthday, phone: phone,
        kind: kind}
        if( users.filter(us => us.email == email).length == 0 ) {
            users.push(u);
            setUser(u);
            return u;
        }

        return null;
    }

    return(
        <UserContext.Provider value={
            {
                connect: connect,
                disconnect: disconnect,
                changeEmail: changeEmail,
                changeInfo: changeInfo,
                changePassword: changePassword,
                createAccount: createAccount,
                getUser: getUser,
            }
        }>
            { prop.children }
        </UserContext.Provider>
    )
}
