import React from "react";
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Box } from "@mui/material";

interface AddCommandDialogProps {
  open: boolean;
  onClose: () => void;
  category: string;
  name: string;
  command: string;
  description: string;
  onCategoryChange: (value: string) => void;
  onNameChange: (value: string) => void;
  onCommandChange: (value: string) => void;
  onDescriptionChange: (value: string) => void;
  onSave: () => void;
}

const AddCommandDialog: React.FC<AddCommandDialogProps> = ({
  open,
  onClose,
  category,
  name,
  command,
  description,
  onCategoryChange,
  onNameChange,
  onCommandChange,
  onDescriptionChange,
  onSave,
}) => {
  return (
    <Dialog 
      open={open} 
      onClose={onClose}
      maxWidth="sm"
      fullWidth
    >
      <DialogTitle>添加自定义命令</DialogTitle>
      <DialogContent>
        <Box sx={{ display: "flex", flexDirection: "column", gap: 2, pt: 1 }}>
          <TextField
            label="分类"
            value={category}
            onChange={(e) => onCategoryChange(e.target.value)}
            fullWidth
            required
          />
          <TextField
            label="名称"
            value={name}
            onChange={(e) => onNameChange(e.target.value)}
            fullWidth
            required
          />
          <TextField
            label="描述"
            value={description}
            onChange={(e) => onDescriptionChange(e.target.value)}
            fullWidth
            multiline
            rows={2}
            minRows={2}
          />
          <TextField
            label="命令"
            value={command}
            onChange={(e) => onCommandChange(e.target.value)}
            fullWidth
            multiline
            rows={4}
            required
            minRows={4}
          />
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} variant="outlined">取消</Button>
        <Button onClick={onSave} variant="contained">
          保存
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default AddCommandDialog;
