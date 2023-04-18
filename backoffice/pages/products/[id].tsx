import Container from "@mui/material/Container";
import Box from "@mui/material/Box";
import * as React from "react";
import {Grid, TextField} from "@mui/material";
import {useEffect, useState} from "react";
import ProductService, {Product} from "../../src/services/ProductService";
import {useRouter} from "next/router";
import Typography from "@mui/material/Typography";


export default function Product() {

    const router = useRouter()
    const id = router.query['id'] as string | null

    const [product, setProduct] = useState<Product | null>(null);

    useEffect(() => {
        if (id != null) {
            ProductService.getProduct(id).then(value => setProduct(value));
        }
    }, [id]);


    return <Container maxWidth="lg">
        <Box
            sx={{
                my: 8,
                display: 'flex',
                flexDirection: 'column',
            }}
        >
            {product && <React.Fragment>
                <Typography variant="h6" gutterBottom>
                    Product {product.id}
                </Typography>
                <Grid container spacing={2}>
                    <Grid item xs={8}>
                        <TextField
                            required
                            id="productName"
                            name="productName"
                            label="Name"
                            defaultValue={product.name}
                        />
                    </Grid>
                    <Grid item xs={8}>
                        <TextField
                            required
                            id="productDescription"
                            name="productDescription"
                            label="Description"
                            defaultValue={product.description}
                        />
                    </Grid>
                </Grid>
            </React.Fragment>}
        </Box>
    </Container>
}