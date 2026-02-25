<?php

namespace App\Services\CarouselSetting;

use App\Models\CarouselSetting;
use Illuminate\Http\UploadedFile;

class CreateCarouselSettingService
{
    public function create(array $data)
    {
        // Handle image upload
        if (isset($data['image']) && $data['image'] instanceof UploadedFile) {
            $imagePath = $this->uploadImage($data['image']);
            $data['image_path'] = $imagePath;
            unset($data['image']);
        }

        // Set default order if not provided
        if (!isset($data['order'])) {
            $maxOrder = CarouselSetting::max('order') ?? -1;
            $data['order'] = $maxOrder + 1;
        }

        return CarouselSetting::create($data);
    }

    private function uploadImage(UploadedFile $image)
    {
        // Create carousel directory if it doesn't exist
        $publicPath = public_path('carousel');
        if (!file_exists($publicPath)) {
            mkdir($publicPath, 0755, true);
        }

        // Generate unique filename
        $filename = time() . '_' . uniqid() . '.' . $image->getClientOriginalExtension();

        // Move file to public/carousel
        $image->move($publicPath, $filename);

        return $filename;
    }
}
