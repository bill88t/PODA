import { useState, type ReactNode } from "react";
import { PathContext } from "./pathContext";

export default function PathProvider(prop: {children: ReactNode}) {
    const [path, setPath] = useState<string>(location.pathname);
    function setUrl(p: string) {
        setPath(p);
        history.pushState({}, "", p);
    }
    return (
        <PathContext.Provider value={{path: path, setPath: setUrl}}>
            {prop.children}
        </PathContext.Provider>
    )
}
