import {useRouter} from "next/router";
import {PropsWithChildren, useContext, useEffect, useState} from "react";
import {paths, publicPaths} from "./paths";
import {Box, CircularProgress} from "@mui/material";
import {AuthActionType, AuthContext, AuthDispatchContext} from "./AuthProvider";
import {isAuthenticated} from "./services/AuthService";


export default function RouteGuard(props: PropsWithChildren): JSX.Element {
    const {children} = props
    const router = useRouter();

    const [isAllowed, setAllowed] = useState(false);

    const auth = useContext(AuthContext);
    const authDispatch = useContext(AuthDispatchContext);

    useEffect(() => {
        const authCheck = () => {
            const isPublicPath = publicPaths.includes(router.asPath.split("?")[0]);
            if (isPublicPath) {
                setAllowed(true);
            } else {
                if (!auth.isAuthenticated) {
                    isAuthenticated().then(authenticated => {
                        if (authenticated) {
                            authDispatch({type: AuthActionType.SET_AUTHENTICATED})
                        } else {
                            void router.push({
                                pathname: paths.signIn
                            });
                        }
                        setAllowed(authenticated);
                    })
                } else {
                    setAllowed(true);
                }
            }
        }

        authCheck();

        router.events.on('routeChangeComplete', authCheck);

        return () => {
            router.events.off('routeChangeComplete', authCheck);
        };
    }, [router, router.events, auth]);

    if (!isAllowed) {
        return <Box sx={{display: 'flex'}}>
            <CircularProgress/>
        </Box>
    } else {
        return children as JSX.Element;
    }
}
