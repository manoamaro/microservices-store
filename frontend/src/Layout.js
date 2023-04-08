import {Container} from "react-bootstrap";
import {Outlet} from "react-router-dom";
import Header from "./components/Header";

function Layout() {
    return (
        <Container fluid>
            <Header/>
            <Container>
                <Outlet/>
            </Container>
        </Container>
    );
}

export default Layout;
