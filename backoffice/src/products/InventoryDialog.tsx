import React, {useState} from "react";
import {Button, Dialog, DialogActions, DialogContent, DialogTitle} from "@mui/material";
import TextField from "@mui/material/TextField";

interface InventoryDialogProps {
    isOpen: boolean;
    initialInventory?: number;
    onClose: () => void;
    onSave: (inventory: number) => void;
}

const InventoryDialog: React.FC<InventoryDialogProps> = ({isOpen, initialInventory, onClose, onSave}) => {
    const [inventory, setInventory] = useState<number>(initialInventory || 0);

    return (
        <Dialog open={isOpen} onClose={onClose}>
            <DialogTitle>Set Inventory</DialogTitle>
            <DialogContent>
                <TextField
                    autoFocus
                    margin="dense"
                    id="inventory"
                    label="Inventory"
                    type="number"
                    fullWidth
                    value={inventory}
                    onChange={({target: {value}}) => setInventory(parseInt(value))}
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button onClick={() => onSave(inventory)}>Save</Button>
            </DialogActions>
        </Dialog>
    );
};

export default InventoryDialog;