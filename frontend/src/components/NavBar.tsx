import { usePath } from "./path/pathContext"
import { UserKind, useUser } from "./user/userContext";

export default function NavBar() {
    const { setPath } = usePath();
    const { getUser } = useUser();
    return (
        <div className="nav">
            <a onClick={ () => setPath("/") }>Home Page</a>
            <div className="nav-right">
                {
                    (getUser().kind === UserKind.guess) ?
                    <>
                        <a onClick={ () => setPath("/login") }>Login</a>
                        <a onClick={ () => setPath("/sign_up") }>Sign up</a>
                    </>
                    : <></>
                }
                {
                    getUser().kind !== UserKind.guess ?
                    <a onClick={ () => setPath("/profile") }>Profile</a>
                    : <></>
                }
            </div>
        </div>
    )
}
