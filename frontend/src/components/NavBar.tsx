function NavBar(prop: { setPath: (arg0: string) => void, path: string }) {
    const setUrl = (s: string) => {
        location.assign(s);
        prop.setPath(s);
    }
    return (
        <div className="nav">
            <div className="nav-left">
                <a onClick={ () => setUrl("/") }>Home Page</a>
                <a onClick={ () => setUrl("/login") }>Login</a>
                <a onClick={ () => setUrl("/sign_up") }>Sign up</a>
            </div>
            <div className="nav-right">
                <a onClick={ () => setUrl("/profile") }>Profile</a>
            </div>
        </div>
    )
}

export default NavBar

