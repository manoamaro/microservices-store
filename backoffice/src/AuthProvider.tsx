import {createContext, Dispatch, PropsWithChildren, useReducer} from "react";


export enum AuthActionType {
    SET_AUTHENTICATED,
    LOGOUT
}

export interface AuthAction {
    type: AuthActionType
}

export interface AuthState {
    isAuthenticated: boolean;
}

function authReducer(state: AuthState, action: AuthAction): AuthState {
    const {type} = action;
    switch (type) {
        case AuthActionType.SET_AUTHENTICATED:
            return {
                ...state,
                isAuthenticated: true
            };
        case AuthActionType.LOGOUT:
            return {
                ...state,
                isAuthenticated: false
            };
        default:
            return state
    }
}

const initialState: AuthState = {isAuthenticated: false}

export const AuthContext = createContext<AuthState>(initialState);
export const AuthDispatchContext = createContext<Dispatch<AuthAction>>(() => {
});

export const AuthProvider = (props: PropsWithChildren): JSX.Element => {
    const [auth, dispatch] = useReducer(authReducer, initialState);
    return (
        <AuthContext.Provider value={auth}>
            <AuthDispatchContext.Provider value={dispatch}>
                {props.children}
            </AuthDispatchContext.Provider>
        </AuthContext.Provider>
    )

}