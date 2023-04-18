import Container from "@mui/material/Container";
import Box from "@mui/material/Box";
import * as React from "react";
import {useEffect, useState} from "react";
import {
    DataGrid,
    GridCellParams,
    GridColDef
} from "@mui/x-data-grid";
import ProductService, {Product} from "../../src/services/ProductService";
import {Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField} from "@mui/material";


interface ProductDialog {
    open: boolean;
    product: Product | null;
}

export default function Index() {

    const [products, setProducts] = useState<Product[]>([]);
    const [productDialog, setProductDialog] = useState<ProductDialog>({open: false, product: null});

    const loadProducts = async () => {
        try {
            const data = await ProductService.getProducts();
            setProducts(data);
        } catch (error) {
            console.error(error);
        }
    };

    const createProduct = async (product: Product) => {
        try {
            await ProductService.postProduct(product);
            await loadProducts();
            setProductDialog({open: false, product: null});
        } catch (error) {
            console.error(error);
        }
    };

    const updateProduct = async (product: Product) => {
        try {
            await ProductService.putProduct(product);
            await loadProducts();
            setProductDialog({open: false, product: null});
        } catch (error) {
            console.error(error);
        }
    };

    const deleteProduct = async (product: Product) => {
        try {
            void loadProducts();
        } catch (error) {
            console.error(error);
        }
    };

    const handleCreateClick = () => setProductDialog({open: true, product: null});
    const handleEditClick = (product: Product) => {
        if (product.id) {
            ProductService.getProduct(product.id).then(p => setProductDialog({open: true, product: p}));
        }
    }
    const handleDeleteClick = (params: GridCellParams) => deleteProduct(params.row as Product);
    const handleClose = () => setProductDialog({open: false, product: null});
    const handleChange = (name: string, value: any) => setProductDialog({
        ...productDialog,
        product: {
            ...productDialog.product || {id: null, name: "", description: "", price: 0},
            [name]: value
        }
    })

    useEffect(() => {
        void loadProducts();
    }, []);

    const columns: GridColDef<Product>[] = [
        {field: "id", headerName: "ID"},
        {field: "name", headerName: "Name"},
        {field: "description", headerName: "Description"},
        {
            field: "actions",
            headerName: "Actions",
            width: 200,
            renderCell: (params: GridCellParams<Product>) => (
                <>
                    <Button onClick={() => handleEditClick(params.row)}>Edit</Button>
                    <Button onClick={() => handleDeleteClick(params)}>Delete</Button>
                </>
            ),
        }
    ]

    return <Container maxWidth="lg">
        <Box sx={{height: 400, width: '100%'}}>
            <DataGrid columns={columns} rows={products}/>
            <Button variant="contained" onClick={handleCreateClick}>
                Create Product
            </Button>
        </Box>
        <Dialog open={productDialog.open} onClose={handleClose}>
            <DialogTitle>{productDialog.product?.id ? "Edit Product" : "Create Product"}</DialogTitle>
            <DialogContent>
                <TextField
                    autoFocus
                    margin="dense"
                    id="name"
                    label="Name"
                    name="name"
                    type="text"
                    fullWidth
                    value={productDialog.product?.name ?? ""}
                    onChange={({target: {name, value}}) => handleChange(name, value)}
                />
                <TextField
                    margin="dense"
                    id="description"
                    label="Description"
                    type="text"
                    name="description"
                    fullWidth
                    value={productDialog.product?.description ?? ""}
                    onChange={({target: {name, value}}) => handleChange(name, value)}
                />
                <TextField
                    margin="dense"
                    id="price"
                    label="Price"
                    type="number"
                    name="price"
                    fullWidth
                    value={productDialog.product?.price ?? ""}
                    onChange={({target: {name, value}}) => handleChange(name, value)}
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose}>Cancel</Button>
                <Button
                    onClick={() => {
                        if (productDialog.product) {
                            if (productDialog.product.id) {
                                void updateProduct(productDialog.product);
                            } else {
                                void createProduct(productDialog.product);
                            }
                        }
                    }}
                >
                    Save
                </Button>
            </DialogActions>
        </Dialog>
    </Container>
}