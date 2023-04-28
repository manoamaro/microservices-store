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
import ProductService, {PRODUCTS_HOST} from "../services/ProductService";

interface ImageDialogProps {
    isOpen: boolean;
    initialImages?: string[];
    onClose: () => void;
    onSave: (images: string[]) => void;
}

const ImageDialog: React.FC<ImageDialogProps> = ({isOpen, initialImages, onClose, onSave}) => {
    const [images, setImages] = useState<string[]>(initialImages || []);

    return (
        <Dialog open={isOpen} onClose={onClose}>
            <DialogTitle>Images</DialogTitle>
            <DialogContent>
                <ImageList sx={{width: 500, height: 450}}>
                    {images.map((image, idx) => (
                        <ImageListItem key={idx}>
                            <img
                                src={`${PRODUCTS_HOST}/public/assets/${image}`}
                                alt={image}
                                loading="lazy"
                            />
                            <ImageListItemBar
                                title=""
                                position="below"
                                actionPosition="left"
                                actionIcon={
                                    <IconButton>
                                        <Delete/>
                                    </IconButton>
                                }
                            />
                        </ImageListItem>
                    ))}
                </ImageList>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button onClick={() => onSave(images)}>Save</Button>
            </DialogActions>
        </Dialog>
    );
};

export default ImageDialog;