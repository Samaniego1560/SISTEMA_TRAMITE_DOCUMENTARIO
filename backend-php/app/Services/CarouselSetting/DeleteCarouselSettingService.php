<?php

namespace App\Services\CarouselSetting;

use App\Models\CarouselSetting;

class DeleteCarouselSettingService
{
    public function delete($id)
    {
        $setting = CarouselSetting::findOrFail($id);

        // Delete image file
        $this->deleteImage($setting->image_path);

        // Delete database record
        return $setting->delete();
    }

    private function deleteImage($imagePath)
    {
        $fullPath = public_path('carousel/' . $imagePath);
        if (file_exists($fullPath)) {
            unlink($fullPath);
        }
    }
}
