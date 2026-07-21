export interface PublicProfile {
    username : string;
    display_name : string;
    bio: string;
    avatar_url: string;
    links: PublicLink[];
}

export interface PublicLink {
    title: string;
    slug: string;
}

