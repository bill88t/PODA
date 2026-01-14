import { createContext, useContext } from "react";

export enum UserKind {
    client = "client",
    barber = "barber",
    admin = "admin",
};

export enum AppointmentKind {
    haircut = "haircut",
    hairdye = "hairdye",
    other = "other",
}

// one to many
export type Appointment = {
    id: number;
    datetime: Date;
    kind: AppointmentKind;
};

export type Uuid = string | number;

export type AuthUser = {
    uuid: Uuid;
    kind: UserKind;
    fname: string;
    lname: string;
    email: string;
    phone: string;
    password: string;
    birthday: Date;
    address: string | null;
    appointments: Appointment [];
};

export type User = null | AuthUser;

export type Tuple<T, V> = [T, V];

export type UserContextType = {
    connect: (email: string, password: string) => Promise<boolean>,
    disconnect: () => void,
    changePassword: (password: string) => Promise<boolean>,
    changeInfo: (fname: string, lname: string, birthdate: Date) => Promise<boolean>,
    changeContact: (email: string, phone? : string) => Promise<boolean>;
    createAccount: (
            fname: string, lname: string,
            email: string, password: string,
            birthday: Date, phone: string,
            address:string | null, kind: UserKind ) => Promise<boolean>;
    createAppointment: (uuid: Uuid, kind: AppointmentKind, datetime: Date) => Promise<boolean>;
    deleteAppointment: (uuid: Uuid, id: number) => Promise<boolean>;
    viewAppointment: (uuid: Uuid, datetime: Date) => Promise<Tuple<Uuid, Appointment[]> | null>
    user: User;
}

export const UserContext = createContext<UserContextType | undefined>(undefined);

export function useUser(): UserContextType {
    const ctx = useContext(UserContext);
    if (ctx == undefined) throw new Error("User Context is not setup");
    return ctx;
}
