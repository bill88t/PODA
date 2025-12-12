import { createContext, type Dispatch, useContext } from "react";

export type PathContextType = {path: string; setPath: Dispatch<string>};

export const PathContext = createContext<PathContextType | undefined>(undefined);

export function usePath(): PathContextType {
    const ctx = useContext(PathContext);
    if (ctx == undefined) throw("Not implemented path ctx");
    return ctx;
}
