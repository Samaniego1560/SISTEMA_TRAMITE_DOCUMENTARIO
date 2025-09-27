<?php
namespace App\Services\CalendarSetting;

use App\Models\CalendarSetting;

class CreateCalendarSettingService
{
    public function create(array $data)
    {
        return CalendarSetting::create($data);
    }

}