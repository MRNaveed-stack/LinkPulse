import { useState } from 'react';
import { useLinks } from '../../hooks/useLinks';
import { useCreateLink } from '../../hooks/useCreateLink';
import { useUpdateLink } from '../../hooks/useUpdateLink';
import { useDeleteLink } from '../../hooks/useDeleteLink';
import LinkCard from '../../components/links/LinkCard';
import LinkForm from '../../components/links/LinkForm';
import DeleteDialog from '../../components/links/DeleteDialog';
import Button from '../../components/common/Button';
import EmptyState from '../../components/common/EmptyState';
import Skeleton from '../../components/common/Skeleton';
import { Search, Plus } from 'lucide-react';
import type { Link } from '../../types/link';

export default function LinksPage() {
  const { data: links, isLoading, isError, error } = useLinks();
  const createMutation = useCreateLink();
  const updateMutation = useUpdateLink();
  const deleteMutation = useDeleteLink();

  const [searchQuery, setSearchQuery] = useState('');
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingLink, setEditingLink] = useState<Link | null>(null);
  const [deletingLink, setDeletingLink] = useState<Link | null>(null);

  // Filter links client-side
  const filteredLinks = links
    ? links.filter(
        (link) =>
          link.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
          link.slug.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : [];

  const handleCreateOrUpdate = (formData: any) => {
    if (editingLink) {
      updateMutation.mutate(
        { id: editingLink.id, data: formData },
        {
          onSuccess: () => {
            setIsFormOpen(false);
            setEditingLink(null);
          },
        }
      );
    } else {
      createMutation.mutate(formData, {
        onSuccess: () => {
          setIsFormOpen(false);
        },
      });
    }
  };

  const handleDelete = () => {
    if (!deletingLink) return;
    deleteMutation.mutate(deletingLink.id, {
      onSuccess: () => setDeletingLink(null),
    });
  };

  if (isError) {
    return (
      <div className="py-12 text-center">
        <h2 className="text-lg font-medium text-red-600">Failed to load links</h2>
        <p className="text-sm text-gray-500 mt-2">
          {(error as any)?.message || 'An unexpected error occurred'}
        </p>
      </div>
    );
  }

  return (
    <div>
      <div className="sm:flex sm:items-center sm:justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-900">My Links</h2>
        <Button
          onClick={() => {
            setEditingLink(null);
            setIsFormOpen(true);
          }}
          className="mt-3 sm:mt-0"
        >
          <Plus className="h-5 w-5 mr-1" /> Create Link
        </Button>
      </div>

      {/* Search */}
      <div className="relative max-w-md mb-6">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
        <input
          type="text"
          placeholder="Search by title or slug..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        />
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="bg-white shadow rounded-lg p-5">
              <Skeleton className="h-5 w-3/4 mb-3" />
              <Skeleton className="h-4 w-full mb-2" />
              <Skeleton className="h-4 w-1/2 mb-4" />
              <Skeleton className="h-10 w-full" />
            </div>
          ))}
        </div>
      ) : filteredLinks.length === 0 && !isLoading ? (
        <EmptyState
          title={searchQuery ? 'No links match your search' : 'No links yet'}
          description={
            searchQuery
              ? 'Try a different search term.'
              : 'Create your first short link to get started.'
          }
          action={
            !searchQuery ? (
              <Button
                onClick={() => {
                  setEditingLink(null);
                  setIsFormOpen(true);
                }}
              >
                <Plus className="h-5 w-5 mr-1" /> Create Link
              </Button>
            ) : undefined
          }
        />
      ) : (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {filteredLinks.map((link) => (
            <LinkCard
              key={link.id}
              link={link}
              onEdit={(link) => {
                setEditingLink(link);
                setIsFormOpen(true);
              }}
              onDelete={(link) => setDeletingLink(link)}
            />
          ))}
        </div>
      )}

      {/* Create/Edit Modal */}
      <LinkForm
        open={isFormOpen}
        onClose={() => {
          setIsFormOpen(false);
          setEditingLink(null);
        }}
        onSubmit={handleCreateOrUpdate}
        isLoading={createMutation.isPending || updateMutation.isPending}
        initialData={editingLink}
      />

      {/* Delete Confirmation Dialog */}
      <DeleteDialog
        open={!!deletingLink}
        onClose={() => setDeletingLink(null)}
        onConfirm={handleDelete}
        isLoading={deleteMutation.isPending}
        linkTitle={deletingLink?.title}
      />
    </div>
  );
}