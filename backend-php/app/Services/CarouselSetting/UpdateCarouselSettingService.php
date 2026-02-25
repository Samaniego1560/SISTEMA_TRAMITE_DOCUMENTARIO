<?php

namespace App\Services\CarouselSetting;

use App\Models\CarouselSetting;
use Illuminate\Http\UploadedFile;

class UpdateCarouselSettingService
{
    public function update(array $data)
    {
        $updatedSettings = [];

        foreach ($data as $item) {
            if (isset($item['id'])) {
                $setting = CarouselSetting::find($item['id']);
                if ($setting) {
                    // Handle image upload if new image provided
                    if (isset($item['image']) && $item['image'] instanceof UploadedFile) {
                        // Delete old image
                        $this->deleteImage($setting->image_path);

                        // Upload new image
                        $imagePath = $this->uploadImage($item['image']);
                        $item['image_path'] = $imagePath;
                        unset($item['image']);
                    }

                    $setting->update($item);
                    $updatedSettings[] = $setting;
                }
            }
        }

        return $updatedSettings;
    }

    public function updateSingle($id, array $data)
    {
        $setting = CarouselSetting::findOrFail($id);

        // Handle image upload if new image provided
        if (isset($data['image']) && $data['image'] instanceof UploadedFile) {
            // Delete old image
            $this->deleteImage($setting->image_path);

            // Upload new image
            $imagePath = $this->uploadImage($data['image']);
            $data['image_path'] = $imagePath;
            unset($data['image']);
        }

        $setting->update($data);
        return $setting;
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

    private function deleteImage($imagePath)
    {
        $fullPath = public_path('carousel/' . $imagePath);
        if (file_exists($fullPath)) {
            unlink($fullPath);
        }
    }
}
