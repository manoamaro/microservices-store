import {AppBar, Box, Button, IconButton, Menu, Toolbar} from "@mui/material";
import Typography from "@mui/material/Typography";
import Link from "./Link";
import {paths} from "./paths";
import * as React from "react";
import {useContext, useEffect, useState} from "react";
import {AuthActionType, AuthContext, AuthDispatchContext} from "./AuthProvider";
import {logout} from "./services/AuthService";

export default function Bar() {

    const authDispatch = useContext(AuthDispatchContext);
    const auth = useContext(AuthContext);
    const [isAuthenticated, setAuthenticated] = useState(false);

    useEffect(() => {
        setAuthenticated(auth.isAuthenticated);
    }, [auth]);

    const logoutAction = () => {
        logout().then(() => {
            authDispatch({type: AuthActionType.LOGOUT});
        });
    }

    if (isAuthenticated) {
        return (
            <AppBar position="static">
                <Toolbar>
                    <IconButton
                        size="large"
                        edge="start"
                        color="inherit"
                        aria-label="menu"
                        sx={{mr: 2}}
                    >
                    </IconButton>
                    <Typography variant="h6" component="div" sx={{flexGrow: 1}}>
                        Store
                    </Typography>

                    <Box sx={{flexGrow: 1, display: {xs: 'none', md: 'flex'}}}>
                        <Button color="inherit" component={Link} noLinkStyle href={paths.products}>Products</Button>
                    </Box>

                    {!isAuthenticated && (
                        <Button color="inherit" component={Link} noLinkStyle href={paths.signIn}>Login</Button>
                    )}
                    {isAuthenticated && (
                        <Button color="inherit" onClick={logoutAction}>Logout</Button>
                    )}
                </Toolbar>
            </AppBar>
        )
    } else {
        return <React.Fragment></React.Fragment>
    }
}