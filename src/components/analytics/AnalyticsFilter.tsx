interface AnalyticsFilterProps {
  selectedDays: number;
  onDaysChange: (days: number) => void;
}

const filters = [
  { label: 'Today', value: 1 },
  { label: 'Last 7 Days', value: 7 },
  { label: 'Last 30 Days', value: 30 },
  { label: 'Last 90 Days', value: 90 },
];

const AnalyticsFilter = ({ selectedDays, onDaysChange }: AnalyticsFilterProps) => (
  <div className="flex items-center space-x-1 bg-gray-100 rounded-lg p-1">
    {filters.map((filter) => (
      <button
        key={filter.value}
        onClick={() => onDaysChange(filter.value)}
        className={`px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
          selectedDays === filter.value
            ? 'bg-white text-gray-900 shadow'
            : 'text-gray-600 hover:text-gray-900'
        }`}
      >
        {filter.label}
      </button>
    ))}
  </div>
);

export default AnalyticsFilter;