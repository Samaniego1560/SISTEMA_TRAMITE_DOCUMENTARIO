<?php
namespace App\Services\CalendarSetting;

use App\Models\CalendarSetting;

class ListCalendarSettingService
{
    public function list()
    {
        return CalendarSetting::get();
    }

}