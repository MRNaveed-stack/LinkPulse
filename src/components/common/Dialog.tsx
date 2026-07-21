import Modal from './Modal';
import Button from './Button';

interface DialogProps {
  open: boolean;
  onClose: () => void;
  title: string;
  description: string;
  confirmLabel?: string;
  onConfirm: () => void;
  loading?: boolean;
  variant?: 'primary' | 'danger';
}

const Dialog = ({
  open,
  onClose,
  title,
  description,
  confirmLabel = 'Confirm',
  onConfirm,
  loading = false,
  variant = 'primary',
}: DialogProps) => (
  <Modal open={open} onClose={onClose}>
    <div>
      <h3 className="text-lg font-medium text-gray-900">{title}</h3>
      <p className="mt-2 text-sm text-gray-500">{description}</p>
    </div>
    <div className="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse gap-3">
      <Button variant={variant} onClick={onConfirm} loading={loading}>
        {confirmLabel}
      </Button>
      <Button variant="secondary" onClick={onClose} disabled={loading}>
        Cancel
      </Button>
    </div>
  </Modal>
);

export default Dialog;