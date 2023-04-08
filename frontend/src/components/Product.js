import {Card} from "react-bootstrap";
import './Product.css';
import {Link} from "react-router-dom";

function Product({product}) {
    return <Card className="mb-4 product-wap rounded-0">
        <Card className="rounded-0">
            <Card.Img src={product.images[0]}/>
        </Card>
        <Card.Body>
            <Card.Title><Link to={`/product/${product.id}`}>{product.name}</Link></Card.Title>
            <Card.Text>{product.description}</Card.Text>
            <p className="text-center">{product.price.currency} {product.price.value / 100}</p>
        </Card.Body>
    </Card>
}


export default Product;