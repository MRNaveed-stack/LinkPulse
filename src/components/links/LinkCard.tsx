import { ExternalLink, Edit2, Trash2 } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import Card from '../common/Card';
import Badge from '../common/Badge';
import StatusSwitch from './StatusSwitch';
import CopyButton from './CopyButton';
import type { Link } from '../../types/link';

interface LinkCardProps {
  link: Link;
  onEdit: (link: Link) => void;
  onDelete: (link: Link) => void;
}

const LinkCard = ({ link, onEdit, onDelete }: LinkCardProps) => {
  const user = useAuthStore((state) => state.user);
  const API_BASE_URL = 'http://localhost:8080';
  const shortUrl = user ? `${API_BASE_URL}/u/${user.username}/${link.slug}` : '';

  return (
    <Card className="p-4 flex flex-col gap-3">
      <div className="flex justify-between items-start">
        <div className="min-w-0 flex-1">
          <h3 className="text-lg font-semibold text-gray-900 truncate">{link.title}</h3>
          <p className="text-sm text-gray-500 truncate">{link.destination_url}</p>
        </div>
        <StatusSwitch active={link.is_active} linkId={link.id} />
      </div>

      <div className="flex items-center gap-2">
        <Badge variant={link.is_active ? 'success' : 'warning'}>
          {link.is_active ? 'Active' : 'Inactive'}
        </Badge>
        {link.click_count > 0 && (
          <span className="text-sm text-gray-500">{link.click_count} clicks</span>
        )}
      </div>

      <div className="flex items-center gap-1">
        <CopyButton text={shortUrl} />
        <a
          href={shortUrl}
          target="_blank"
          rel="noopener"
          referrerPolicy="no-referrer-when-downgrade"
          className="p-2 text-gray-400 hover:text-gray-600 rounded-full hover:bg-gray-100"
          title="Open link"
        >
          <ExternalLink className="h-5 w-5" />
        </a>
        <button
          onClick={() => onEdit(link)}
          className="p-2 text-gray-400 hover:text-indigo-600 rounded-full hover:bg-gray-100"
          title="Edit link"
        >
          <Edit2 className="h-5 w-5" />
        </button>
        <button
          onClick={() => onDelete(link)}
          className="p-2 text-gray-400 hover:text-red-600 rounded-full hover:bg-gray-100"
          title="Delete link"
        >
          <Trash2 className="h-5 w-5" />
        </button>
      </div>
    </Card>
  );
};

export default LinkCard;