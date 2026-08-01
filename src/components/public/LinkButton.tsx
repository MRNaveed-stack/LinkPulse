import type { PublicLink } from '../../types/public';
import apiClient from '../../api/client';

interface LinkButtonProps {
  link: PublicLink;
  username: string;
}
const LinkButton = ({ link, username }: LinkButtonProps) => {
  const apiBaseUrl = apiClient.defaults.baseURL || '/api';
  const redirectUrl = `${apiBaseUrl}/u/${username}/${link.slug}`;
  return (
    <a
      href={redirectUrl}
      target="_blank"
      rel="noopener"
      referrerPolicy="no-referrer-when-downgrade"
      className="block w-full text-center px-5 py-3.5 bg-white border border-gray-200 rounded-lg text-gray-900 font-medium hover:border-indigo-300 hover:shadow-sm hover:bg-indigo-50/30 transition-all duration-200 active:scale-[0.98] focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
    >
      {link.title}
    </a>
  );
};
export default LinkButton;