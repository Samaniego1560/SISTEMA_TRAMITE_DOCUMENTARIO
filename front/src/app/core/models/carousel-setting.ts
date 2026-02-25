export interface ICarouselSetting {
    id?: number;
    image_path: string;
    image_url?: string;
    title?: string;
    description?: string;
    button_text?: string;
    button_link?: string;
    is_enabled: boolean;
    order: number;
    created_at?: string;
    updated_at?: string;
}
