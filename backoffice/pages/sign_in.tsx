import * as React from 'react';
import {useContext, useState} from 'react';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import FormControlLabel from '@mui/material/FormControlLabel';
import Checkbox from '@mui/material/Checkbox';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import {useRouter} from "next/router";
import {Snackbar} from "@mui/material";
import {AuthActionType, AuthDispatchContext} from "../src/AuthProvider";
import AuthService from "../src/services/AuthService";

export default function SignIn() {

    const router = useRouter();
    const [isLoading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const authDispatch = useContext(AuthDispatchContext);

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        setLoading(true);

        const form = event.target as typeof event.target & {
            email: { value: string },
            password: { value: string }
        }

        AuthService.signIn(form.email.value, form.password.value).then(() => {
            authDispatch({type: AuthActionType.SET_AUTHENTICATED});
            void router.push("/");
        }).catch(e => {
            setError(e.message);
        }).finally(() => setLoading(false));
    };

    return (
        <Container component="main" maxWidth="xs">
            <Snackbar open={error != null} message={error} onClose={() => setError(null)}/>
            <Box
                sx={{
                    marginTop: 8,
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                }}
            >
                <Typography component="h1" variant="h5">
                    Sign in
                </Typography>
                <Box component="form" onSubmit={handleSubmit} noValidate sx={{mt: 1}}>
                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        id="email"
                        label="Email Address"
                        name="email"
                        autoComplete="email"
                        autoFocus
                        disabled={isLoading}
                    />
                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        name="password"
                        label="Password"
                        type="password"
                        id="password"
                        autoComplete="current-password"
                        disabled={isLoading}
                    />
                    <FormControlLabel
                        control={<Checkbox value="remember" color="primary"/>}
                        label="Remember me"
                        disabled={isLoading}
                    />
                    <Button
                        type="submit"
                        fullWidth
                        variant="contained"
                        sx={{mt: 3, mb: 2}}
                        disabled={isLoading}
                    >
                        Sign In
                    </Button>
                    <Grid container>
                        <Grid item xs>
                            <Link href="#" variant="body2">
                                Forgot password?
                            </Link>
                        </Grid>
                        <Grid item>
                            <Link href="#" variant="body2">
                                {"Don't have an account? Sign Up"}
                            </Link>
                        </Grid>
                    </Grid>
                </Box>
            </Box>
        </Container>
    );
}