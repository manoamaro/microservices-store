import {Carousel, Col, Container, Row} from "react-bootstrap";
import Header from "../components/Header";
import {useEffect, useState} from "react";
import {getProduct, getProducts} from "../services/Product";
import Product from "../components/Product";
import {useParams} from "react-router-dom";
import popover from "bootstrap/js/src/popover";

function ProductPage() {
    const [product, setProduct] = useState(null);
    let {id} = useParams();

    useEffect(() => {
        const fetchData = async () => {
            const product = await getProduct(id)
            setProduct(product)
        };

        fetchData().catch(console.error);
    }, [id]);

    if (product == null) {
        return <p>Loading</p>
    } else {
        return (
            <Row>
                <Carousel>
                    {product.images.map(image =>
                        <Carousel.Item>
                            <img src={image} alt=""/>
                        </Carousel.Item>
                    )}
                </Carousel>
            </Row>
        );
    }
}

export default ProductPage;
