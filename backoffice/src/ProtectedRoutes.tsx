import {useRouter} from "next/router";
import {PropsWithChildren, useEffect, useState} from "react";
import {isAuthenticated} from "./services/authentication";


export default function ProtectedRoutes(props: PropsWithChildren): JSX.Element {
    const {children} = props
    const router = useRouter();

    const [isLoading, setLoading] = useState(false);

    useEffect(() => {
        isAuthenticated().then(value => {
            if (!value) {
                router.push("/sign_in");
            }
        }).finally(() => setLoading(false))
    }, []);

    if (isLoading) {
        return <div/>
    } else {
        return children as JSX.Element;
    }
}
