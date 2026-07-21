import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { createLinkSchema, updateLinkSchema } from '../../schemas/link';
import Modal from '../common/Modal';
import Button from '../common/Button';
import type { Link } from '../../types/link';

interface LinkFormProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: any) => void;
  isLoading: boolean;
  initialData?: Link | null; // if editing
}

const LinkForm = ({ open, onClose, onSubmit, isLoading, initialData }: LinkFormProps) => {
  const schema = initialData ? updateLinkSchema : createLinkSchema;

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm({
    resolver: zodResolver(schema),
    defaultValues: {
      title: initialData?.title || '',
      slug: initialData?.slug || '',
      destination_url: initialData?.destination_url || '',
    },
  });

  useEffect(() => {
    if (initialData) {
      reset({
        title: initialData.title,
        slug: initialData.slug,
        destination_url: initialData.destination_url,
      });
    } else {
      reset({ title: '', slug: '', destination_url: '' });
    }
  }, [initialData, reset]);

  const handleFormSubmit = (data: any) => {
    onSubmit(data);
  };

  return (
    <Modal open={open} onClose={onClose}>
      <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
        <h3 className="text-lg font-medium text-gray-900">
          {initialData ? 'Edit Link' : 'Create New Link'}
        </h3>

        <div>
          <label className="block text-sm font-medium text-gray-700">Title</label>
          <input
            {...register('title')}
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.title && <p className="mt-1 text-sm text-red-600">{errors.title.message as string}</p>}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">Slug</label>
          <input
            {...register('slug')}
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.slug && <p className="mt-1 text-sm text-red-600">{errors.slug.message as string}</p>}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">Destination URL</label>
          <input
            {...register('destination_url')}
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.destination_url && (
            <p className="mt-1 text-sm text-red-600">{errors.destination_url.message as string}</p>
          )}
        </div>

        <div className="flex justify-end gap-3 pt-4">
          <Button variant="secondary" type="button" onClick={onClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button type="submit" loading={isLoading}>
            {initialData ? 'Save Changes' : 'Create Link'}
          </Button>
        </div>
      </form>
    </Modal>
  );
};

export default LinkForm;