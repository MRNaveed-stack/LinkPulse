import { useToggleStatus } from '../../hooks/useToggleStatus';

interface StatusSwitchProps {
  active: boolean;
  linkId: string;
}

const StatusSwitch = ({ active, linkId }: StatusSwitchProps) => {
  const toggleMutation = useToggleStatus();

  return (
    <button
      type="button"
      onClick={() => toggleMutation.mutate({ id: linkId, is_active: !active })}
      disabled={toggleMutation.isPending}
      className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 ${
        active ? 'bg-indigo-600' : 'bg-gray-200'
      }`}
    >
      <span className="sr-only">Toggle status</span>
      <span
        className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
          active ? 'translate-x-6' : 'translate-x-1'
        }`}
      />
    </button>
  );
};

export default StatusSwitch;