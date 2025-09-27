<?php
namespace App\Services\ApplicationNotice;

use App\Models\ApplicationNotice;

class ListApplicationNoticeService
{
    public function list()
    {
        return ApplicationNotice::get();
    }

}