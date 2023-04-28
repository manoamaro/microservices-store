import Container from "@mui/material/Container";
import * as React from "react";
import {useEffect, useState} from "react";
import ProductService, {Product} from "../../src/services/ProductService";
import {
    IconButton,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead, TablePagination,
    TableRow,
    Tooltip
} from "@mui/material";
import InventoryService from "../../src/services/InventoryService";
import ProductDialog from "../../src/products/ProductDialog";
import InventoryDialog from "../../src/products/InventoryDialog";
import {Delete, Edit, ImageTwoTone, Inventory} from "@mui/icons-material";
import ImageDialog from "../../src/products/ImagesDialog";


interface ProductWithInventory extends Product {
    inventory: number
}

const EMPTY_PRODUCT: ProductWithInventory = {
    id: "",
    name: "",
    description: "",
    prices: [],
    inventory: 0,
    images: []
};

export default function Index() {

    const [products, setProducts] = useState<ProductWithInventory[]>([]);
    const [productDialog, setProductDialog] = useState<boolean>(false);
    const [inventoryDialog, setInventoryDialog] = useState<boolean>(false);
    const [imageDialog, setImageDialog] = useState<boolean>(false);
    const [selectedProduct, setSelectedProduct] = useState<ProductWithInventory>(EMPTY_PRODUCT);

    const loadProducts = async () => {
        try {
            const products = await ProductService.getProducts();
            const productsWithInventory = await Promise.all(products.map(async product => {
                const inventory = await InventoryService.getInventory(product.id || "");
                return {...product, inventory: inventory.amount}
            }));
            setProducts(productsWithInventory);
        } catch (error) {
            console.error(error);
        }
    };

    const createOrUpdateProduct = async (product: Product) => {
        try {
            if (product.id) {
                await ProductService.putProduct(product);
            } else {
                await ProductService.postProduct(product);
            }
            await loadProducts();
            setSelectedProduct(EMPTY_PRODUCT);
            setProductDialog(false);
        } catch (error) {
            console.error(error);
        }
    };

    const deleteProduct = async (product: ProductWithInventory) => {
        try {
            void loadProducts();
        } catch (error) {
            console.error(error);
        }
    };

    const handleCreateClick = () => {
        setSelectedProduct(EMPTY_PRODUCT);
        setProductDialog(true);
    }

    const handleEditClick = (product: ProductWithInventory) => {
        if (product.id) {
            setSelectedProduct(product);
            setProductDialog(true);
        }
    }

    const handleInventoryClick = (product: ProductWithInventory) => {
        if (product.id) {
            setSelectedProduct(product);
            setInventoryDialog(true);
        }
    }

    const handleImagesClick = (product: ProductWithInventory) => {
        if (product.id) {
            setSelectedProduct(product);
            setImageDialog(true);
        }
    }

    const handleDeleteClick = (product: ProductWithInventory) => deleteProduct(product);
    const handleClose = () => setProductDialog(false);
    const handleInventoryClose = () => setInventoryDialog(false);

    useEffect(() => {
        void loadProducts();
    }, []);


    return <Container maxWidth="lg">
        <TableContainer component={Paper}>
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell>ID</TableCell>
                        <TableCell>Name</TableCell>
                        <TableCell>Description</TableCell>
                        <TableCell>Inventory</TableCell>
                        <TableCell>Actions</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {products.map((product, idx) => <TableRow key={idx}>
                        <TableCell>{product.id}</TableCell>
                        <TableCell>{product.name}</TableCell>
                        <TableCell>{product.description}</TableCell>
                        <TableCell>{product.inventory}</TableCell>
                        <TableCell>
                            <Tooltip title="Edit">
                                <IconButton onClick={() => handleEditClick(product)}>
                                    <Edit/>
                                </IconButton>
                            </Tooltip>
                            <Tooltip title="Inventory">
                                <IconButton onClick={() => handleInventoryClick(product)}>
                                    <Inventory/>
                                </IconButton>
                            </Tooltip>
                            <Tooltip title="Images">
                                <IconButton onClick={() => handleImagesClick(product)}>
                                    <ImageTwoTone/>
                                </IconButton>
                            </Tooltip>
                            <Tooltip title="Delete">
                                <IconButton onClick={() => handleDeleteClick(product)}>
                                    <Delete/>
                                </IconButton>
                            </Tooltip>
                        </TableCell>
                    </TableRow>)}
                </TableBody>
            </Table>
        </TableContainer>

        <ProductDialog
            key={`product-dialog-${selectedProduct.id}`}
            isOpen={productDialog}
            initialProduct={selectedProduct}
            onClose={handleClose}
            onSave={createOrUpdateProduct}/>

        <InventoryDialog
            key={`inventory-dialog-${selectedProduct.id}`}
            isOpen={inventoryDialog}
            initialInventory={selectedProduct.inventory}
            onClose={handleInventoryClose}
            onSave={async inventory => {
                await InventoryService.putInventory(selectedProduct.id || "", inventory);
                await loadProducts();
                handleInventoryClose();
            }}/>

        <ImageDialog
            key={`image-dialog-${selectedProduct.id}`}
            isOpen={imageDialog}
            onClose={() => setImageDialog(false)}
            initialImages={selectedProduct.images}
            onSave={async images => {
            }}/>

    </Container>
}