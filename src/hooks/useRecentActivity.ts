import { useQuery } from "@tanstack/react-query";
import { getRecentActivity } from "../api/analytics";

export const useRecentActivity = (limit: number = 20) => {
    return useQuery({
        queryKey: ["analytics", "recent-activity", limit],
        queryFn: () => getRecentActivity(limit),
    });
}