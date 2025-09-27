<?php

namespace App\Services\CalendarSetting;

use Carbon\Carbon;
use App\Models\CalendarSetting;
use Illuminate\Support\Facades\DB;

class UpdateCalendarSettingService
{
    public function update(array $data)
    {
        $currentTimestamp = Carbon::now();

        foreach ($data as &$item) {
            $item['created_at'] = $currentTimestamp;
            $item['updated_at'] = $currentTimestamp;
        }

        DB::transaction(function () use ($data) {
            // Borra todos los registros existentes
            CalendarSetting::truncate();

            CalendarSetting::insert($data);
        });

        return CalendarSetting::get();
    }
}
