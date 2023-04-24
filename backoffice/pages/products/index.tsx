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
import {Button} from "@mui/material";
import InventoryService from "../../src/services/InventoryService";
import ProductDialog from "../../src/products/ProductDialog";
import InventoryDialog from "../../src/products/InventoryDialog";


interface ProductWithInventory extends Product {
    inventory: number
}

export default function Index() {

    const [products, setProducts] = useState<ProductWithInventory[]>([]);
    const [productDialog, setProductDialog] = useState<boolean>(false);
    const [inventoryDialog, setInventoryDialog] = useState<boolean>(false);
    const EMPTY_PRODUCT = {
        id: "",
        name: "",
        description: "",
        prices: [],
        inventory: 0
    };
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

    const deleteProduct = async () => {
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

    const handleDeleteClick = () => deleteProduct();
    const handleClose = () => setProductDialog(false);
    const handleInventoryClose = () => setInventoryDialog(false);

    useEffect(() => {
        void loadProducts();
    }, []);

    const columns: GridColDef<ProductWithInventory>[] = [
        {field: "id", headerName: "ID"},
        {field: "name", headerName: "Name"},
        {field: "description", headerName: "Description"},
        {field: "inventory", headerName: "Inventory"},
        {
            field: "actions",
            headerName: "Actions",
            flex: 1,
            renderCell: (params: GridCellParams<ProductWithInventory>) => (
                <>
                    <Button onClick={() => handleEditClick(params.row)}>Edit</Button>
                    <Button onClick={() => handleDeleteClick()}>Delete</Button>
                    <Button onClick={() => handleInventoryClick(params.row)}>Inventory</Button>
                </>
            ),
        }
    ]

    return <Container maxWidth="lg">
        <Box sx={{height: 800, width: '100%'}}>
            <DataGrid columns={columns} rows={products}/>
            <Button variant="contained" onClick={handleCreateClick}>
                Create Product
            </Button>
        </Box>

        <ProductDialog
            key={selectedProduct.id}
            isOpen={productDialog}
            initialProduct={selectedProduct}
            onClose={handleClose}
            onSave={createOrUpdateProduct}/>

        <InventoryDialog
            key={selectedProduct.id}
            isOpen={inventoryDialog}
            initialInventory={selectedProduct.inventory}
            onClose={handleInventoryClose}
            onSave={async inventory => {
                await InventoryService.putInventory(selectedProduct.id || "", inventory);
                await loadProducts();
                handleInventoryClose();
            }}/>

    </Container>
}