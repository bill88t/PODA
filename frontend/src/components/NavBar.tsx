import { usePath } from "./path/pathContext"
import { useUser } from "./user/userContext";

export default function NavBar() {
    const { setPath } = usePath();
    const { user, disconnect } = useUser();
    const noLogged = user === null;
    return (
        <div className="nav">
            <a onClick={e => {
                e.preventDefault();
                setPath("/");
            }}>Home Page</a>
            {noLogged ?
                <>
                </>
                :
                <>
                    <a onClick={e => {
                        e.preventDefault();
                        setPath("/events");
                    }}>Events</a>
                </>}
            <div className="nav-right">
                {noLogged ?
                    <>
                        <a onClick={e => {
                            e.preventDefault();
                            setPath("/login");
                        }}>Login</a>
                        <a onClick={e => {
                            e.preventDefault();
                            setPath("/sign_up");
                        }}>Sign up</a>
                    </>
                    :
                    <>
                        <a onClick={e => {
                            e.preventDefault();
                            disconnect();
                            setPath("/");
                        }}>Log out</a>
                        <a onClick={e => {
                            e.preventDefault();
                            setPath("/profile");
                        }}>Profile</a>
                    </>}
            </div>
        </div>
    )
}
