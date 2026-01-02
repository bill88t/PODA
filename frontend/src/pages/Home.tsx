import { UserKind, useUser } from "../components/user/userContext";

function BasicHome() {
    return (
        <>
            <title>Sallon PODA - Home Page</title>
            <h1>Welcome to Our Saloon</h1>
            <p>
                The PADA Saloon is open from 2025
                our history is new, but we would like to create
                stories for our community.
                We would like to have you here with us
                and help us leverage your hairstyle,
                because you deserve that.
                You are amazing and we are gonna make it
                appear, so everyone is gonna see it.
            </p>
        </>
    )
}

function UserHome() {
    const userCtx = useUser();
    if (userCtx.user == null) throw Error("UserHome user is null");
    return (
        <>
            <title>Sallon PODA - Home Page</title>
            <h1>Welcome back {userCtx.user.fname}!</h1>
        </>
    )
}

function ClientHome() {
    const userCtx = useUser();
    if (userCtx.user == null) throw Error("ClientHome user is null");
    return (
        <>
            <UserHome />
        </>
    )
}

function BarberHome() {
    const userCtx = useUser();
    if (userCtx.user == null) throw Error("BarberHome user is null");
    return (
        <>
            <UserHome />
        </>
    )
}

function AdminHome() {
    const userCtx = useUser();
    if (userCtx.user == null) throw Error("AdminHome user is null");
    return (
        <>
            <UserHome />
        </>
    )
}

function Home() {
    const userCtx = useUser();
    if (userCtx.user      == null)            return <BasicHome />;
    if (userCtx.user.kind == UserKind.client) return <ClientHome />;
    if (userCtx.user.kind == UserKind.barber) return <BarberHome />;
    if (userCtx.user.kind == UserKind.admin)  return <AdminHome />;
}

export default Home
