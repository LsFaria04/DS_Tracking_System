import React from "react";
import {
  Button,
  Dialog,
  DialogContent,
  DialogActions,
  CircularProgress,
  Box,
} from "@mui/material";

interface ModalProps {
  show: boolean;
  onClose: () => void;
  onUpdate: () => void;
  isUpdating: boolean;
  children: React.ReactNode;
}

export default function UpdateModal ({
  show,
  onClose,
  onUpdate,
  isUpdating,
  children,
}: ModalProps) {

  return (
    <Dialog open={show} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogContent>
        <Box sx={{ py: 2 }}>{children}</Box>
      </DialogContent>
      <DialogActions sx={{ gap: 1, p: 2 }}>
        <Button
          onClick={onUpdate}
          id="send-update"
          disabled={isUpdating}
          variant="contained"
          color="primary"
          sx={{
            minWidth: "120px",
          }}
        >
          {isUpdating ? (
            <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
              <CircularProgress size={16} color="inherit" />
              Updating...
            </Box>
          ) : (
            "Update"
          )}
        </Button>

        <Button
          onClick={onClose}
          id="cancel-update"
          disabled={isUpdating}
          variant="outlined"
          color="inherit"
          sx={{
            minWidth: "120px",
          }}
        >
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  );
}