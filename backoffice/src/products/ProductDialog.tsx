import {Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField} from "@mui/material";
import * as React from "react";
import {Product} from "../services/ProductService";
import Box from "@mui/material/Box";
import {useState} from "react";

interface EditProductDialogProps {
    isOpen: boolean;
    initialProduct?: Product;
    onClose: () => void;
    onSave: (product: Product) => void;
}

const ProductDialog: React.FC<EditProductDialogProps> = (props: EditProductDialogProps) => {
    const [product, setProduct] = useState<Product>(props.initialProduct || {
        id: "",
        name: "",
        description: "",
        prices: []
    });

    const handleChange = (name: string, value: any) => setProduct({
        ...product,
        [name]: value
    });

    const handlePriceCurrencyChange = (index: number, value: string) => {
        const updatedPrices = [...(product.prices || [])];
        updatedPrices[index].currency = value;
        handleChange("prices", updatedPrices);
    };

    const handlePricePriceChange = (index: number, value: number) => {
        const updatedPrices = [...(product.prices || [])];
        updatedPrices[index].price = value;
        handleChange("prices", updatedPrices);
    };

    const handleAddPrice = () => {
        handleChange("prices", [...(product.prices || []), {currency: "", price: 0}]);
    };

    const handleRemovePrice = (index: number) => {
        const updatedPrices = [...(product.prices || [])];
        updatedPrices.splice(index, 1);
        handleChange("prices", updatedPrices);
    };

    return (
        <Dialog open={props.isOpen} onClose={props.onClose}>
            <DialogTitle>{product.id ? "Edit Product" : "Create Product"}</DialogTitle>
            <DialogContent>
                <TextField
                    autoFocus
                    margin="dense"
                    id="name"
                    label="Name"
                    name="name"
                    type="text"
                    fullWidth
                    value={product.name ?? ""}
                    onChange={({target: {name, value}}) => handleChange(name, value)}
                />
                <TextField
                    margin="dense"
                    id="description"
                    label="Description"
                    type="text"
                    name="description"
                    fullWidth
                    value={product.description ?? ""}
                    onChange={({target: {name, value}}) => handleChange(name, value)}
                />

                {(product.prices || []).map((price, index) => (
                    <Box key={index}>
                        <TextField
                            label="Currency"
                            name="currency"
                            value={price.currency}
                            onChange={({target: {name, value}}) => handlePriceCurrencyChange(index, value)}
                        />
                        <TextField
                            label="Price"
                            name="price"
                            value={price.price}
                            onChange={({target: {name, value}}) => handlePricePriceChange(index, +value)}
                        />
                        <Button
                            variant="contained"
                            color="secondary"
                            onClick={() => handleRemovePrice(index)}
                        >
                            Remove Price
                        </Button>
                    </Box>
                ))}
                <Button
                    variant="contained"
                    color="primary"
                    onClick={handleAddPrice}
                >
                    Add Price
                </Button>
            </DialogContent>
            <DialogActions>
                <Button onClick={props.onClose}>Cancel</Button>
                <Button onClick={() => props.onSave(product)}>Save</Button>
            </DialogActions>
        </Dialog>
    )
};

export default ProductDialog;