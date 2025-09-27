<?php

namespace App\Services\ApplicationNotice;

use Carbon\Carbon;
use App\Models\ApplicationNotice;
use Illuminate\Support\Facades\DB;

class UpdateApplicationNoticeService
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
            ApplicationNotice::truncate();

            ApplicationNotice::insert($data);
        });

        return ApplicationNotice::get();
    }
}
