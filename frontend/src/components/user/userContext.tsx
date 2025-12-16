import { createContext, useContext } from "react";

export enum UserKind {
    client = "client",
    barber = "barber",
    admin = "admin",
    guess = "guess"
}

export type User = {
    kind: UserKind;
    fname?: string;
    lname?: string;
    email?: string;
    phone?: string;
    password?: string;
    birthday?: Date;
}

export type UserContextType = {
    connect: (email: string, password: string) => User,
    disconnect: () => void,
    changePassword: (password: string) => User,
    changeInfo: (fname: string, lname: string, birthdate: Date) => User,
    changeEmail: (email: string) => User
    createAccount: (
            fname: string, lname: string,
            email: string, password: string,
            birthday: Date, phone: string,
            kind: UserKind ) => User | null;
    getUser(): User;
}

export const UserContext = createContext<UserContextType | undefined>(undefined);

export function useUser(): UserContextType {
    const ctx = useContext(UserContext);
    if (ctx == undefined) throw("User Context is not setup");
    return ctx;
}
