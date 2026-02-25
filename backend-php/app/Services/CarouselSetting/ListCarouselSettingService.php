<?php

namespace App\Services\CarouselSetting;

use App\Models\CarouselSetting;

class ListCarouselSettingService
{
    public function list()
    {
        return CarouselSetting::ordered()->get();
    }

    public function listEnabled()
    {
        return CarouselSetting::enabled()->get();
    }
}
