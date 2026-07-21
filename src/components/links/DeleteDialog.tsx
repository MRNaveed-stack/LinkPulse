import Dialog from '../common/Dialog';

interface DeleteDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: () => void;
  isLoading: boolean;
  linkTitle?: string;
}

const DeleteDialog = ({ open, onClose, onConfirm, isLoading, linkTitle }: DeleteDialogProps) => (
  <Dialog
    open={open}
    onClose={onClose}
    title="Delete Link"
    description={`Are you sure you want to delete "${linkTitle || 'this link'}"? This action cannot be undone.`}
    confirmLabel="Delete"
    variant="danger"
    onConfirm={onConfirm}
    loading={isLoading}
  />
);

export default DeleteDialog;