import { Inbox } from 'lucide-react';

const EmptyProfile = () => (
  <div className="flex flex-col items-center py-12 text-center">
    <Inbox className="h-12 w-12 text-gray-300 mb-3" />
    <p className="text-gray-500 text-sm">No links available yet.</p>
  </div>
);

export default EmptyProfile;
