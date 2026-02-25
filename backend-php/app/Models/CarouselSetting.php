<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class CarouselSetting extends Model
{
    protected $fillable = [
        'image_path',
        'title',
        'description',
        'button_text',
        'button_link',
        'is_enabled',
        'order'
    ];

    protected $casts = [
        'is_enabled' => 'boolean',
        'order' => 'integer'
    ];

    /**
     * Get the full URL for the image
     */
    public function getImageUrlAttribute()
    {
        return url('carousel/' . $this->image_path);
    }

    /**
     * Scope to get only enabled items ordered
     */
    public function scopeEnabled($query)
    {
        return $query->where('is_enabled', true)->orderBy('order');
    }

    /**
     * Scope to get all items ordered
     */
    public function scopeOrdered($query)
    {
        return $query->orderBy('order');
    }
}
