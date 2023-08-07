import React, {useState} from "react";
import {
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle, IconButton, ImageList,
    ImageListItem,
    ImageListItemBar
} from "@mui/material";
import TextField from "@mui/material/TextField";
import {Delete} from "@mui/icons-material";
import ProductService, {EMPTY_PRODUCT, Product, ProductImage, PRODUCTS_HOST} from "../services/ProductService";

interface ImageDialogProps {
    isOpen: boolean;
    product?: Product;
    onClose: () => void;
}

const ImageDialog: React.FC<ImageDialogProps> = ({isOpen, product, onClose}) => {
    const [images, setImages] = useState<ProductImage[]>(product?.images || []);

    const handleUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
        if (event.target.files && product?.id) {
            const updatedProduct = await ProductService.postProductImages(product.id, event.target.files);
            setImages(updatedProduct.images);
        }
    }

    const handleDelete = async (image: ProductImage) => {
        if (product?.id) {
            const updatedProduct = await ProductService.deleteProductImage(product.id, image.id);
            setImages(updatedProduct.images || []);
        }
    }

    return (
        <Dialog open={isOpen} onClose={onClose}>
            <DialogTitle>Images</DialogTitle>
            <DialogContent>
                <ImageList sx={{width: 500, height: 450}}>
                    {images.map((image, idx) => (
                        <ImageListItem key={idx}>
                            <img
                                src={`${image.url}`}
                                alt={image.description}
                                loading="lazy"
                            />
                            <ImageListItemBar
                                title=""
                                position="below"
                                actionPosition="left"
                                actionIcon={
                                    <IconButton onClick={() => handleDelete(image)}>
                                        <Delete/>
                                    </IconButton>
                                }
                            />
                        </ImageListItem>
                    ))}
                </ImageList>
            </DialogContent>
            <DialogActions>
                <Button variant="contained" component="label">
                    Upload File
                    <input type="file" hidden accept="image/*" multiple onChange={handleUpload}/>
                </Button>
                <Button onClick={onClose}>Cancel</Button>
            </DialogActions>
        </Dialog>
    );
};

export default ImageDialog;