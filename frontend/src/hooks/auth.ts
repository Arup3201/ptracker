import { useContext } from "react";

import { AuthContext } from "@/contexts/auth-context";

const useAuth = () => {
    const ctx = useContext(AuthContext)
    return ctx
}

export default useAuth