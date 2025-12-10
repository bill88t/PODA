function NavBar(prop: { setPath: (arg0: string) => void }) {
    return (
        <div className="nav">
            <div className="nav-left">
                <a onClick={ () => prop.setPath("/") }>Home Page</a>
                <a onClick={ () => prop.setPath("/login") }>Login</a>
                <a onClick={ () => prop.setPath("/sign_up") }>Sign up</a>
            </div>
            <div className="nav-right">
                <a onClick={ () => prop.setPath("/profile") }>Profile</a>
            </div>
        </div>
    )
}

export default NavBar

