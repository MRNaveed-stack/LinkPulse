import { useState } from 'react';
import Modal from '../common/Modal';
import Button from '../common/Button';
import { useDeleteAccount } from '../../hooks/useDeleteAccount';
import { AlertTriangle } from 'lucide-react';

interface DeleteAccountDialogProps {
  open: boolean;
  onClose: () => void;
}

export default function DeleteAccountDialog({ open, onClose }: DeleteAccountDialogProps) {
  const deleteAccountMutation = useDeleteAccount();
  const [confirmText, setConfirmText] = useState('');

  const handleDelete = () => {
    if (confirmText !== 'DELETE') return;
    deleteAccountMutation.mutate(undefined, {
      onSuccess: () => {
        onClose();
      },
    });
  };

  const isDeleteEnabled = confirmText === 'DELETE';

  return (
    <Modal open={open} onClose={onClose}>
      <div className="sm:flex sm:items-start">
        <div className="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10">
          <AlertTriangle className="h-6 w-6 text-red-600" />
        </div>
        <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
          <h3 className="text-lg leading-6 font-semibold text-gray-900">
            Delete Account
          </h3>
          <div className="mt-2">
            <p className="text-sm text-gray-500 leading-relaxed">
              Are you absolutely sure you want to delete your account? This action is <strong className="text-red-600 font-semibold">irreversible</strong> and will permanently delete your:
            </p>
            <ul className="mt-2 list-disc pl-5 text-xs text-gray-500 space-y-1">
              <li>Personal public landing page profile</li>
              <li>All short links created by you</li>
              <li>All tracked click analytics and referrer history</li>
            </ul>
          </div>
        </div>
      </div>

      <div className="mt-6 border-t border-gray-100 pt-4">
        <label htmlFor="confirm-delete-input" className="block text-sm font-medium text-gray-700">
          To confirm deletion, type <strong className="text-gray-900 font-bold select-all bg-gray-50 px-1.5 py-0.5 rounded border border-gray-200">DELETE</strong> below:
        </label>
        <input
          type="text"
          id="confirm-delete-input"
          value={confirmText}
          onChange={(e) => setConfirmText(e.target.value)}
          placeholder="Type DELETE to confirm"
          className="mt-2 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm"
        />
      </div>

      <div className="mt-6 flex flex-col-reverse sm:flex-row sm:justify-end gap-3 border-t border-gray-100 pt-4">
        <Button
          variant="secondary"
          onClick={() => {
            setConfirmText('');
            onClose();
          }}
          className="w-full sm:w-auto"
        >
          Cancel
        </Button>
        <Button
          variant="danger"
          loading={deleteAccountMutation.isPending}
          disabled={!isDeleteEnabled}
          onClick={handleDelete}
          className="w-full sm:w-auto"
        >
          Delete My Account
        </Button>
      </div>
    </Modal>
  );
}
