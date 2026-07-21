export interface Link{
    id : string;
    user_id: string;
    title: string;
    slug: string;
    destination_url: string;
    is_active: boolean;
    click_count: number;
    created_at : string;
    updated_at : string;
}

export interface CreateLinkRequest {
    title: string;
    slug: string;
    destination_url: string;
}

export interface UpdateLinkRequest {
    title ?: string;
    slug?: string;
    destination_url?: string;
}

export interface ToggleStatusRequest {
    is_active : boolean;
}

export interface LinkListResponse {
    links: Link[];
}

export interface LinkStatusResponse{
    message: string;
    is_active: boolean;
}