import {Col, Container, Row} from "react-bootstrap";
import Header from "../components/Header";
import {useEffect, useState} from "react";
import {getProducts} from "../services/Product";
import Product from "../components/Product";

function Home() {
    const [products, setProducts] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            const products = await getProducts()
            setProducts(products)
        };

        fetchData().catch(console.error);
    }, []);

    return (
        <Row>
            {products.map(value => <Col md="4" key={value.id}><Product product={value}/></Col>)}
        </Row>
    );
}

export default Home;
